package service

import (
	"context"
	"fmt"
	"tiktok_e-commerce/app/user/biz/dao"
	user_service "tiktok_e-commerce/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/klog"
)

type DeleteUserService struct {
	ctx context.Context
} // NewDeleteUserService new DeleteUserService
func NewDeleteUserService(ctx context.Context) *DeleteUserService {
	return &DeleteUserService{ctx: ctx}
}

// Run create note info
func (s *DeleteUserService) Run(req *user_service.DeleteUserReq) (resp *user_service.DeleteUserResp, err error) {
	// Finish your business logic.

	userDAO := dao.GetUserDAO()
	err = userDAO.Delete(req.UserId)
	if err != nil {
		klog.Error("删除用户失败:", err)
		return &user_service.DeleteUserResp{
			ResponseStatus: buildResponse("删除用户失败", false),
		}, nil
	}

	// 同步删除redis缓存
	cacheKey := fmt.Sprintf("user:%d", req.UserId)
	if err := userDAO.Cache().Del(s.ctx, cacheKey).Err(); err != nil {
		klog.Error("删除用户redis缓存失败:", err)
	}
	klog.Info("删除用户成功")
	return &user_service.DeleteUserResp{
		ResponseStatus: buildResponse("删除用户成功", true),
	}, nil
}
