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

type DeleteUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewDeleteUserService(Context context.Context, RequestContext *app.RequestContext) *DeleteUserService {
	return &DeleteUserService{RequestContext: RequestContext, Context: Context}
}

func (h *DeleteUserService) Run(req *user.Empty) (resp *user.DeleteUserResponse, err error) {
	ctx := h.Context
	deleteUserResp, err := rpc.UserClient.DeleteUser(ctx, &rpcuser.DeleteUserReq{
		UserId: ctx.Value("user_id").(int32),
	})
	if err != nil {
		hlog.Error("删除用户失败: %v", err)
		return nil, errors.New("删除用户失败，请稍后再试")
	}
	return &user.DeleteUserResponse{
		StatusCode: deleteUserResp.StatusCode,
		StatusMsg:  deleteUserResp.StatusMsg,
	}, nil
}
