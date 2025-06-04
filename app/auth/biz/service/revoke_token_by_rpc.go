package service

import (
	"context"
	"tiktok_e-commerce/auth/biz/dal/redis"
	redisUtils "tiktok_e-commerce/auth/utils/redis"
	"tiktok_e-commerce/common/constant"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"

	"github.com/pkg/errors"
)

type RevokeTokenByRPCService struct {
	ctx context.Context
} // NewRevokeTokenByRPCService new RevokeTokenByRPCService
func NewRevokeTokenByRPCService(ctx context.Context) *RevokeTokenByRPCService {
	return &RevokeTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *RevokeTokenByRPCService) Run(req *auth.RevokeTokenReq) (resp *auth.RevokeResp, err error) {
	// Finish your business logic.
	err = redis.RedisClient.Del(s.ctx, redisUtils.GetAccessTokenKey(req.UserId), redisUtils.GetRefreshTokenKey(req.UserId)).Err()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return &auth.RevokeResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}, nil
}
