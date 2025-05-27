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

type LogoutService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewLogoutService(Context context.Context, RequestContext *app.RequestContext) *LogoutService {
	return &LogoutService{RequestContext: RequestContext, Context: Context}
}

func (h *LogoutService) Run(req *user.LogoutReq) (resp *user.LogoutResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	userrpcclient.InitUserRpcClient()

	res, err := userrpccl.Logout(h.Context, &user_service.LogoutReq{
		UserId: req.UserId,
	})
	if err != nil {
		log.Printf("Failed to logout user: %v", err)
		h.RequestContext.JSON(consts.StatusBadRequest, &user.LogoutResp{
			ResponseStatus: &user.ResponseStatus{
				Status:  false,
				Message: err.Error(),
			},
		})
		return
	}
	resp = &user.LogoutResp{
		ResponseStatus: &user.ResponseStatus{
			Status:  res.ResponseStatus.Status,
			Message: res.ResponseStatus.Message,
		},
	}
	return resp, nil
}
