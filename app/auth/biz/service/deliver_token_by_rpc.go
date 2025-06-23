package service

import (
	"context"
	"tiktok_e-commerce/auth/infra/kafka/producer"
	"tiktok_e-commerce/auth/utils/jwt"
	"tiktok_e-commerce/common/constant"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"

	"github.com/cloudwego/kitex/pkg/klog"
)

type DeliverTokenByRPCService struct {
	ctx context.Context
} // NewDeliverTokenByRPCService new DeliverTokenByRPCService
func NewDeliverTokenByRPCService(ctx context.Context) *DeliverTokenByRPCService {
	return &DeliverTokenByRPCService{ctx: ctx}
}

// Run create note info
func (s *DeliverTokenByRPCService) Run(req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	refreshToken, err := jwt.GenerateRefreshToken(s.ctx, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "生成refresh token失败，req: %v, error: %v", req, err)
		return nil, err
	}
	accessToken, err := jwt.GenerateAccessToken(s.ctx, req.UserId, req.Role)
	if err != nil {
		klog.CtxErrorf(s.ctx, "生成access token失败，req: %v, error: %v", req, err)
		return nil, err
	}

	producer.SendUserCacheMessage(req.UserId)
	return &auth.DeliveryResp{
		StatusCode:   0,
		StatusMsg:    constant.GetMsg(0),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}, nil
}
