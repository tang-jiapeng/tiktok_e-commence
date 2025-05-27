package service

import (
	"context"

	user "tiktok_e-commerce/app/hertz/hertz_gen/hertz/user"
	"tiktok_e-commerce/app/hertz/rpc_client/user_rpc"
	user_service "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	userrpccl "tiktok_e-commerce/rpc_gen/rpc/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

type GetUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserService(Context context.Context, RequestContext *app.RequestContext) *GetUserService {
	return &GetUserService{RequestContext: RequestContext, Context: Context}
}

func (h *GetUserService) Run(req *user.GetUserReq) (resp *user.GetUserResp, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	userrpcclient.InitUserRpcClient()

	res, err := userrpccl.GetUser(h.Context, &user_service.GetUserReq{
		UserId: req.UserId,
	})

	if err != nil {
		h.RequestContext.String(consts.StatusBadRequest, err.Error())
		return
	}

	resp = &user.GetUserResp{
		UserId:          res.UserId,
		UserName:        res.UserName,
		Email:           res.Email,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
		UserPermissions: res.UserPermissions,
		ResponseStatus: &user.ResponseStatus{
			Status:  res.ResponseStatus.Status,
			Message: res.ResponseStatus.Message,
		},
	}

	return resp, nil
}
