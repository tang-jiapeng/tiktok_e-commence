package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	user "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	"tiktok_e-commerce/user/biz/dal/mysql"
	myRedis "tiktok_e-commerce/user/biz/dal/redis"
	"tiktok_e-commerce/user/biz/model"
	redisUtils "tiktok_e-commerce/user/utils/redis"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
)

type GetReceiveAddressService struct {
	ctx context.Context
} // NewGetReceiveAddressService new GetReceiveAddressService
func NewGetReceiveAddressService(ctx context.Context) *GetReceiveAddressService {
	return &GetReceiveAddressService{ctx: ctx}
}

// Run create note info
func (s *GetReceiveAddressService) Run(req *user.GetReceiveAddressReq) (resp *user.GetReceiveAddressResp, err error) {
	ctx := s.ctx
	userId := req.UserId

	// 判断key是否存在
	addressKey := redisUtils.GetUserAddressKey(userId)
	exists, err := myRedis.RedisClient.Exists(ctx, addressKey).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	var receiveAddresses []*user.ReceiveAddress
	if exists == 0 {
		addressList, err := s.SelectAndCacheUserAddresses(ctx, userId)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		for _, address := range addressList {
			receiveAddresses = append(receiveAddresses, &user.ReceiveAddress{
				Id:            address.ID,
				Name:          address.Name,
				PhoneNumber:   address.PhoneNumber,
				DefaultStatus: int32(address.DefaultStatus),
				Province:      address.Province,
				City:          address.City,
				Region:        address.Region,
				DetailAddress: address.DetailAddress,
			})
		}
	} else {
		cachedAddresses, err := myRedis.RedisClient.LRange(ctx, addressKey, 0, -1).Result()
		if err != nil && errors.Is(err, redis.Nil) {
			return nil, errors.WithStack(err)
		}
		for _, cachedAddress := range cachedAddresses {
			var receiveAddress user.ReceiveAddress
			_ = sonic.Unmarshal([]byte(cachedAddress), &receiveAddress)
			receiveAddresses = append(receiveAddresses, &receiveAddress)
		}
	}
	return &user.GetReceiveAddressResp{
		StatusCode:     0,
		StatusMsg:      constant.GetMsg(0),
		ReceiveAddress: receiveAddresses,
	}, nil
}

func (s *GetReceiveAddressService) SelectAndCacheUserAddresses(ctx context.Context, userId int32) ([]model.AddressInfo, error) {
	addresses, err := model.GetAddressList(mysql.DB, ctx, userId)
	if err != nil {
		klog.CtxErrorf(ctx, "数据库查询地址信息失败，err: %v", err)
		return nil, errors.WithStack(err)
	}
	if len(addresses) == 0 {
		return addresses, nil
	}

	// 缓存并设置过期时间和access token的过期时间一致
	luaScript := `
		local key = KEYS[1]
		if redis.call('EXISTS', key) == 0 then
			redis.call('RPUSH', key, unpack(ARGV))
			redis.call('EXPIRE', key, 7200)
		end
		return 1
	`
	key := redisUtils.GetUserAddressKey(userId)
	redisClient := myRedis.RedisClient
	addressStrs := make([]string, len(addresses))
	for i, address := range addresses {
		addressStr, _ := sonic.Marshal(address)
		addressStrs[i] = string(addressStr)
	}
	err = redisClient.Eval(ctx, luaScript, []string{key}, addressStrs).Err()
	if err != nil {
		klog.CtxErrorf(ctx, "redis缓存地址信息失败，err: %v", err)
		return nil, errors.WithStack(err)
	}
	return addresses, nil
}
