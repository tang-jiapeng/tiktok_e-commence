package service

import (
	"context"
	"tiktok_e-commerce/auth/utils/jwt"
	"tiktok_e-commerce/auth/utils/redis"
	"tiktok_e-commerce/common/constant"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"

	"github.com/pkg/errors"
)

type VerifyTokenByRPCService struct {
	ctx context.Context
} // NewVerifyTokenByRPCService new VerifyTokenByRPCService
func NewVerifyTokenByRPCService(ctx context.Context) *VerifyTokenByRPCService {
	return &VerifyTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *VerifyTokenByRPCService) Run(req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	
	// 校验access token
	claims, status := jwt.ParseJWT(req.AccessToken)
	if status != jwt.TokenValid {
		return &auth.VerifyResp{
			StatusCode: 1004,
			StatusMsg:  constant.GetMsg(1004),
		}, nil
	}
	userId := int32(claims["userId"].(float64))
	// role := claims["role"].(string)
	// 校验 redis 中的access token
	savedAccessToken, err := redis.GetVal(s.ctx, redis.GetAccessTokenKey(userId))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if savedAccessToken != req.AccessToken {
		return &auth.VerifyResp{
			StatusCode: 1005,
			StatusMsg:  constant.GetMsg(1005),
		}, nil
	}

	// casbin 鉴权


	return &auth.VerifyResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		UserId:     userId,
	}, nil
}
