package service

import (
	"context"
	"tiktok_e-commerce/app/user/biz/dao"
	"tiktok_e-commerce/app/user/biz/model"
	user_service "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
)

type RegisterService struct {
	ctx context.Context
} // NewRegisterService new RegisterService
func NewRegisterService(ctx context.Context) *RegisterService {
	return &RegisterService{ctx: ctx}
}

// Run create note info
func (s *RegisterService) Run(req *user_service.RegisterReq) (resp *user_service.RegisterResp, err error) {
	// Finish your business logic.

	// 检查用户名和邮箱是否已存在
	userDAO := dao.GetUserDAO()
	if existingUser, err := userDAO.FindByEmail(req.Email); err != nil {
		klog.Error("检查邮箱是否存在时出错:", err)
		return &user_service.RegisterResp{
			ResponseStatus: buildResponse("注册失败，请稍后重试！", false),
		}, nil
	} else if existingUser != nil {
		return &user_service.RegisterResp{
			ResponseStatus: buildResponse("该邮箱已注册", false),
		}, nil
	}

	if existingUser, err := userDAO.FindByUsername(req.UserName); err != nil {
		klog.Error("检查用户名是否存在时出错:", err)
		return &user_service.RegisterResp{
			ResponseStatus: buildResponse("注册失败，请稍后重试！", false),
		}, nil
	} else if existingUser != nil {
		return &user_service.RegisterResp{
			ResponseStatus: buildResponse("该用户名已注册", false),
		}, nil
	}

	// 验证密码和确认密码是否一致
	if req.Password != req.ConfirmPassword {
		return &user_service.RegisterResp{
			ResponseStatus: buildResponse("两次密码不一致", false),
		}, nil
	}

	// 生成随机盐值
	salt, err := generateSalt(16)
	if err != nil {
		klog.Error("生成随机盐值时出错:", err)
		return &user_service.RegisterResp{
			ResponseStatus: buildResponse("注册失败，请稍后重试！", false),
		}, nil
	}
	// 加密密码
	hashedPassword, err := hashPasswordWithSalt(req.Password, salt)
	if err != nil {
		klog.Error("加密密码时出错:", err)
		return &user_service.RegisterResp{
			ResponseStatus: buildResponse("注册失败，请稍后重试！", false),
		}, nil
	}

	data := &model.User{
		Email:           req.Email,
		Username:        req.UserName,
		Password:        hashedPassword,
		Salt:            salt,
		UserPermissions: req.UserPermissions,
		Created_at:      time.Now(),
		Updated_at:      time.Now(),
	}

	if err := userDAO.Insert(data); err != nil {
		klog.Error("用户注册失败:", err)
		return &user_service.RegisterResp{
			ResponseStatus: buildResponse("注册失败，请稍后重试！", false),
		}, nil
	}

	return &user_service.RegisterResp{
		UserId:         data.UserId,
		ResponseStatus: buildResponse("注册成功", true),
	}, nil
}
