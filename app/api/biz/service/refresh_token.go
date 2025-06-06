package service

import (
	"context"

	user "tiktok_e-commerce/api/hertz_gen/api/user"
	"tiktok_e-commerce/api/infra/rpc"
	"tiktok_e-commerce/rpc_gen/kitex_gen/auth"

	"github.com/cloudwego/hertz/pkg/app"
)

type RefreshTokenService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRefreshTokenService(Context context.Context, RequestContext *app.RequestContext) *RefreshTokenService {
	return &RefreshTokenService{RequestContext: RequestContext, Context: Context}
}

func (h *RefreshTokenService) Run(req *user.Empty) (resp *user.LoginResponse, err error) {
	authClient := rpc.AuthClient
	refreshTokenResp, err := authClient.RefreshTokenByRPC(h.Context, &auth.RefreshTokenReq{
		RefreshToken: string(h.RequestContext.GetHeader("refresh_token")),
	})
	if err != nil {
		return nil, err
	}
	return &user.LoginResponse{
		AccessToken:  refreshTokenResp.AccessToken,
		RefreshToken: refreshTokenResp.RefreshToken,
		StatusCode:   int32(refreshTokenResp.StatusCode),
		StatusMsg:    refreshTokenResp.StatusMsg,
	}, nil
}
