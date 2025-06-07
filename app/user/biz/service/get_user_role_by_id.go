package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	user "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	"tiktok_e-commerce/user/biz/dal/mysql"
	"tiktok_e-commerce/user/biz/model"

	"github.com/cloudwego/kitex/pkg/klog"
)

type GetUserRoleByIdService struct {
	ctx context.Context
} // NewGetUserRoleByIdService new GetUserRoleByIdService
func NewGetUserRoleByIdService(ctx context.Context) *GetUserRoleByIdService {
	return &GetUserRoleByIdService{ctx: ctx}
}

// Run create note info
func (s *GetUserRoleByIdService) Run(req *user.GetUserRoleByIdReq) (resp *user.GetUserRoleByIdResp, err error) {
	role, err := model.GetUserRoleById(mysql.DB, s.ctx, req.UserId)
	if err != nil {
		klog.Errorf("查询用户角色失败，userId: %d, err: %v", req.UserId, err)
		return nil, err
	}
	resp = &user.GetUserRoleByIdResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		Role:       role,
	}
	return
}
