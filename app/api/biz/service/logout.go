package service

import (
	"context"
	"errors"
	"fmt"

	user "tiktok_e-commerce/api/hertz_gen/api/user"
	"tiktok_e-commerce/rpc_gen/kitex_gen/auth"
	"tiktok_e-commerce/api/infra/rpc"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type LogoutService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLogoutService(Context context.Context, RequestContext *app.RequestContext) *LogoutService {
	return &LogoutService{RequestContext: RequestContext, Context: Context}
}

func (h *LogoutService) Run(req *user.Empty) (resp *user.LogoutResponse, err error) {
	logoutResp , err := rpc.AuthClient.RevokeTokenByRPC(h.Context , &auth.RevokeTokenReq{
		UserId: h.Context.Value("user_id").(int32),
	}) 
	fmt.Println(logoutResp)
	if err != nil {
		hlog.CtxErrorf(h.Context, "rpc调用登出接口失败: %v", err)
		return nil, errors.New("登出失败，请稍后再试")
	}
	resp = &user.LogoutResponse{
		StatusCode: logoutResp.StatusCode,
		StatusMsg:	logoutResp.StatusMsg,
	}
	return
}
