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

type UpdateUserInfoService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewUpdateUserInfoService(Context context.Context, RequestContext *app.RequestContext) *UpdateUserInfoService {
	return &UpdateUserInfoService{RequestContext: RequestContext, Context: Context}
}

func (h *UpdateUserInfoService) Run(req *user.UpdateUserInfoRequest) (resp *user.UpdateUserInfoResponse, err error) {
	ctx := h.Context
	updateUserResp, err := rpc.UserClient.UpdateUser(ctx, &rpcuser.UpdateUserReq{
		UserId:      ctx.Value("user_id").(int32),
		Username:    req.Username,
		Email:       req.Email,
		Sex:         req.Sex,
		Description: req.Description,
		Avatar:      req.Avatar,
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "更新用户信息失败: %v", err)
		return nil, errors.New("更新用户信息失败，请稍后再试")
	}
	return &user.UpdateUserInfoResponse{
		StatusCode: updateUserResp.StatusCode,
		StatusMsg:  updateUserResp.StatusMsg,
	}, nil
}
