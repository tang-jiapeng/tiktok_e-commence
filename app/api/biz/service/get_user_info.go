package service

import (
	"context"

	user "tiktok_e-commerce/api/hertz_gen/api/user"
	"tiktok_e-commerce/api/infra/rpc"
	"tiktok_e-commerce/common/constant"
	rpcuser "tiktok_e-commerce/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/hertz/pkg/app"
)

type GetUserInfoService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserInfoService(Context context.Context, RequestContext *app.RequestContext) *GetUserInfoService {
	return &GetUserInfoService{RequestContext: RequestContext, Context: Context}
}

func (h *GetUserInfoService) Run(req *user.Empty) (resp *user.GetUserInfoResponse, err error) {
	userClient := rpc.UserClient
	ctx := h.Context
	userInfo, err := userClient.GetUser(ctx, &rpcuser.GetUserReq{
		UserId: ctx.Value("user_id").(int32),
	})
	if err != nil {
		return nil, err
	}
	return &user.GetUserInfoResponse{
		StatusCode:  0,
		StatusMsg:   constant.GetMsg(0),
		Username:    userInfo.User.Username,
		Email:       userInfo.User.Email,
		Sex:         userInfo.User.Sex,
		Description: userInfo.User.Description,
		Avatar:      userInfo.User.Avatar,
		CreatedAt:   userInfo.User.CreatedAt,
	}, nil

}
