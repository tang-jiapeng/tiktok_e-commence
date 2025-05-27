package service

import (
	"context"
	"log"

	userrpcclient "tiktok_e-commerce/app/hertz/rpc_client/user_rpc"
	user "tiktok_e-commerce/app/hertz/hertz_gen/hertz/user"
	user_service "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	userrpccl "tiktok_e-commerce/rpc_gen/rpc/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type DeleteUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteUserService(Context context.Context, RequestContext *app.RequestContext) *DeleteUserService {
	return &DeleteUserService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteUserService) Run(req *user.DeleteUserReq) (resp *user.DeleteUserResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code

	userrpcclient.InitUserRpcClient()

	res, err := userrpccl.DeleteUser(h.Context, &user_service.DeleteUserReq{
		UserId: req.UserId,
	})

	if err != nil {
		log.Printf("Failed to delete user: %v", err)
		h.RequestContext.JSON(consts.StatusBadRequest, &user.DeleteUserResp{
			ResponseStatus: &user.ResponseStatus{
				Status:  false,
				Message: err.Error(),
			},
		})
		return
	}

	resp = &user.DeleteUserResp{
		ResponseStatus: &user.ResponseStatus{
			Status:  res.ResponseStatus.Status,
			Message: res.ResponseStatus.Message,
		},
	}

	return resp, nil
}
