package service

import (
	"context"
	"log"

	user "tiktok_e-commerce/app/hertz/hertz_gen/hertz/user"
	userrpcclient "tiktok_e-commerce/app/hertz/rpc_client/user_rpc"
	user_service "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	userrpccl "tiktok_e-commerce/rpc_gen/rpc/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type UpdateUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateUserService(Context context.Context, RequestContext *app.RequestContext) *UpdateUserService {
	return &UpdateUserService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateUserService) Run(req *user.UpdateUserReq) (resp *user.UpdateUserResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code

	userrpcclient.InitUserRpcClient()
	res, err := userrpccl.UpdateUser(h.Context, &user_service.UpdateUserReq{
		UserId:          req.UserId,
		NewEmail:        req.NewEmail,
		NewUserName:     req.NewUserName,
		CurrentPassword: req.CurrentPassword,
		NewPassword:     req.NewPassword,
	})

	if err != nil {
		log.Printf("Failed to update user: %v", err)
		h.RequestContext.JSON(consts.StatusBadRequest, &user.UpdateUserResp{
			ResponseStatus: &user.ResponseStatus{
				Status:  false,
				Message: err.Error(),
			},
		})
		return
	}

	resp = &user.UpdateUserResp{
		ResponseStatus: &user.ResponseStatus{
			Status:  res.ResponseStatus.Status,
			Message: res.ResponseStatus.Message,
		},
	}

	return resp, nil
}
