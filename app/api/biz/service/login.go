package service

import (
	"context"
	"errors"
	user "tiktok_e-commerce/api/hertz_gen/api/user"
	"tiktok_e-commerce/api/infra/rpc"
	rpcuser "tiktok_e-commerce/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type LoginService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLoginService(Context context.Context, RequestContext *app.RequestContext) *LoginService {
	return &LoginService{RequestContext: RequestContext, Context: Context}
}

func (h *LoginService) Run(req *user.LoginRequest) (resp *user.LoginResponse, err error) {
	client := rpc.UserClient
	res, err := client.Login(h.Context, &rpcuser.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		hlog.CtxErrorf(h.Context, "rpc调用登录失败, err: %v", err)
		return nil, errors.New("登录失败，请稍后再试")
	}
	resp = &user.LoginResponse{
		StatusCode:   res.StatusCode,
		StatusMsg:    res.StatusMsg,
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
	}
	return
}
