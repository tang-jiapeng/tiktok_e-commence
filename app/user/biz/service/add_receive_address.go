package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	user "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	"tiktok_e-commerce/user/biz/dal/mysql"
	"tiktok_e-commerce/user/biz/dal/redis"
	"tiktok_e-commerce/user/biz/model"
	redisUtils "tiktok_e-commerce/user/utils/redis"

	"github.com/pkg/errors"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/gorm"
)

type AddReceiveAddressService struct {
	ctx context.Context
} // NewAddReceiveAddressService new AddReceiveAddressService
func NewAddReceiveAddressService(ctx context.Context) *AddReceiveAddressService {
	return &AddReceiveAddressService{ctx: ctx}
}

// Run create note info
func (s *AddReceiveAddressService) Run(req *user.AddReceiveAddressReq) (resp *user.AddReceiveAddressResp, err error) {
	ctx := s.ctx
	addr := req.ReceiveAddress
	err = mysql.DB.Transaction(func(tx *gorm.DB) error {
		if addr.DefaultStatus == model.AddressDefaultStatusYes {
			existingAddr, err := model.ExistDefaultAddress(tx, ctx, req.UserId)
			if err != nil {
				if !errors.Is(err, gorm.ErrRecordNotFound) {
					klog.CtxErrorf(ctx, "查询默认地址是否存在失败，req：%v，err：%v", req, err)
					return errors.WithStack(err)
				}
			} else {
				existingAddr.DefaultStatus = model.AddressDefaultStatusNo
				err = model.UpdateAddress(mysql.DB, ctx, existingAddr)
				if err != nil {
					klog.CtxErrorf(ctx, "更新默认地址失败，req：%v，err：%v", req, err)
					return errors.WithStack(err)
				}
			}
		}
		address := model.Address{
			UserId:        req.UserId,
			Name:          addr.Name,
			PhoneNumber:   addr.PhoneNumber,
			DefaultStatus: int8(addr.DefaultStatus),
			Province:      addr.Province,
			City:          addr.City,
			Region:        addr.Region,
			DetailAddress: addr.DetailAddress,
		}
		err = model.CreateAddress(mysql.DB, ctx, &address)
		if err != nil {
			klog.CtxErrorf(ctx, "添加收货地址失败，req：%v，err：%v", req, err)
			return errors.WithStack(err)
		}
		// 缓存删除失败则回滚事务，防止数据不一致
		err = redis.RedisClient.Del(ctx, redisUtils.GetUserAddressKey(req.UserId)).Err()
		if err != nil {
			klog.CtxErrorf(ctx, "删除用户地址缓存失败，req：%v，err：%v", req, err)
			return errors.WithStack(err)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}
	resp = &user.AddReceiveAddressResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return resp, nil
}
