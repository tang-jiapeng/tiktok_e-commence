package service

import (
	"context"
	"log"

	user "tiktok_e-commerce/app/hertz/hertz_gen/hertz/user"
	"tiktok_e-commerce/app/hertz/rpc_client/user_rpc"
	user_service "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	userrpccl "tiktok_e-commerce/rpc_gen/rpc/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type LoginService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLoginService(Context context.Context, RequestContext *app.RequestContext) *LoginService {
	return &LoginService{RequestContext: RequestContext, Context: Context}
}

func (h *LoginService) Run(req *user.LoginReq) (resp *user.LoginResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code

	userrpcclient.InitUserRpcClient()
	res, err := userrpccl.Login(h.Context, &user_service.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		log.Printf("Failed to login user: %v", err)
		h.RequestContext.JSON(consts.StatusBadRequest, &user.LoginResp{
			ResponseStatus: &user.ResponseStatus{
				Status:  false,
				Message: err.Error(),
			},
		})
		return
	}
	resp = &user.LoginResp{
		ResponseStatus: &user.ResponseStatus{
			Status:  res.ResponseStatus.Status,
			Message: res.ResponseStatus.Message,
		},
	}
	return resp, nil
}
