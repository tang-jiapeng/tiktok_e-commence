package service

import (
	"context"
	"tiktok_e-commerce/auth/utils"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"
)

type DeliverTokenByRPCService struct {
	ctx context.Context
} // NewDeliverTokenByRPCService new DeliverTokenByRPCService
func NewDeliverTokenByRPCService(ctx context.Context) *DeliverTokenByRPCService {
	return &DeliverTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *DeliverTokenByRPCService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	refreshToken, err := utils.GenerateRefreshToken(req.UserId)
	if err != nil {
		return nil , err 
	}
	accessToken , err := utils.GenerateAccessToken(req.UserId)
	if err != nil {
		return nil , err
	}
	return &auth.DeliveryResp{
		StatusCode:	0,
		StatusMsg:	"success",
		AccessToken:	accessToken,
		RefreshToken:	refreshToken,
	}, nil
}
