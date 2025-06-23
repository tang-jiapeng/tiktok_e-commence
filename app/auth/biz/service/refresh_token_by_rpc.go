package service

import (
	"context"
	"tiktok_e-commerce/auth/infra/kafka/producer"
	"tiktok_e-commerce/auth/utils/jwt"
	"tiktok_e-commerce/common/constant"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"
)

type RefreshTokenByRPCService struct {
	ctx context.Context
} // NewRefreshTokenByRPCService new RefreshTokenByRPCService
func NewRefreshTokenByRPCService(ctx context.Context) *RefreshTokenByRPCService {
	return &RefreshTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *RefreshTokenByRPCService) Run(req *auth.RefreshTokenReq) (resp *auth.RefreshTokenResp, err error) {

	userId, newAccessToken, newRefreshToken, success := jwt.RefreshAccessToken(s.ctx, req.RefreshToken)
	if success {
		resp = &auth.RefreshTokenResp{
			StatusCode:   0,
			StatusMsg:    constant.GetMsg(0),
			AccessToken:  newAccessToken,
			RefreshToken: newRefreshToken,
		}
		producer.SendUserCacheMessage(userId)
		return
	}
	resp = &auth.RefreshTokenResp{
		StatusCode: 1006,
		StatusMsg:  constant.GetMsg(1006),
	}
	return
}
