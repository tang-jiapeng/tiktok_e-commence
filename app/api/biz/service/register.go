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

type RegisterService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewRegisterService(Context context.Context, RequestContext *app.RequestContext) *RegisterService {
	return &RegisterService{RequestContext: RequestContext, Context: Context}
}

func (h *RegisterService) Run(req *user.RegisterRequest) (resp *user.RegisterResponse, err error) {
	client := rpc.UserClient
	res, err := client.Register(h.Context, &rpcuser.RegisterReq{
		Username:        req.Username,
		Email:           req.Email,
		Password:        req.Password,
		ConfirmPassword: req.ConfirmPassword,
		Sex:             req.Sex,
		Description:     req.Description,
		Avatar:          req.Avatar,
	})

	if err != nil {
		hlog.Error("register failed, err: ", err)
		return nil, errors.New("注册失败，请稍后再试")
	}
	resp = &user.RegisterResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
	}
	return
}
