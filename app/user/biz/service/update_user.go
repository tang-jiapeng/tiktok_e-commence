package service

import (
	"context"
	"tiktok_e-commerce/app/user/biz/dao"
	user_service "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/crypto/bcrypt"
)

type UpdateUserService struct {
	ctx context.Context
} // NewUpdateUserService new UpdateUserService
func NewUpdateUserService(ctx context.Context) *UpdateUserService {
	return &UpdateUserService{ctx: ctx}
}

// Run create note info
func (s *UpdateUserService) Run(req *user_service.UpdateUserReq) (resp *user_service.UpdateUserResp, err error) {
	// Finish your business logic.

	userDAO := dao.GetUserDAO()
	user, err := userDAO.FindOne(s.ctx, req.UserId)
	if err != nil {
		klog.Info("用户不存在:", err)
		return &user_service.UpdateUserResp{
			ResponseStatus: buildResponse("查询用户失败，请稍后重试！", false),
		}, nil
	}
	if user == nil {
		return &user_service.UpdateUserResp{
			ResponseStatus: buildResponse("用户不存在", false),
		}, nil
	}

	// 验证密码
	saltedCurrentPassword := req.CurrentPassword + user.Salt
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(saltedCurrentPassword)); err != nil {
		klog.Info("密码错误:", err)
		return &user_service.UpdateUserResp{
			ResponseStatus: buildResponse("密码错误，请重新输入！", false),
		}, nil
	}

	// 更新邮箱
	if req.NewEmail != "" {
		user.Email = req.NewEmail
	}

	// 更新密码
	if req.NewPassword != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword+user.Salt), bcrypt.DefaultCost)
		if err != nil {
			klog.Info("密码加密失败:", err)
			return &user_service.UpdateUserResp{
				ResponseStatus: buildResponse("更新失败，请稍后重试！", false),
			}, nil
		}
		user.Password = string(hashedPassword)
	}

	// 更新用户名
	if req.NewPassword != "" {
		user.Username = req.NewUserName
	}

	// 更新数据库
	user.Updated_at = time.Now()
	if err := userDAO.Update(user); err != nil {
		klog.Info("更新用户失败:", err)
		return &user_service.UpdateUserResp{
			ResponseStatus: buildResponse("更新用户失败，请稍后重试！", false),
		}, nil
	}

	// 同步更新redis缓存
	cacheKey := user.GetCacheKey()
	if err := userDAO.Cache().Set(s.ctx, cacheKey, user, time.Hour).Err(); err != nil {
		klog.Error("更新redis缓存失败:", err)
	}
	klog.Info("更新成功")
	return &user_service.UpdateUserResp{
		ResponseStatus: buildResponse("更新成功", true),
	}, nil
}
