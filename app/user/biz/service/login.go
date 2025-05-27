package service

import (
	"context"
	"tiktok_e-commerce/app/user/biz/dao"
	user_service "tiktok_e-commerce/rpc_gen/kitex_gen/user"

	"github.com/cloudwego/kitex/pkg/klog"
	"golang.org/x/crypto/bcrypt"
)

type LoginService struct {
	ctx context.Context
} // NewLoginService new LoginService
func NewLoginService(ctx context.Context) *LoginService {
	return &LoginService{ctx: ctx}
}

// Run create note info
func (s *LoginService) Run(req *user_service.LoginReq) (resp *user_service.LoginResp, err error) {
	// Finish your business logic.

	userDAO := dao.GetUserDAO()
	if userDAO == nil {
		klog.Error("userDAO 未初始化!")
		return &user_service.LoginResp{
			ResponseStatus: buildResponse("系统错误，请稍后重试！", false),
		}, nil
	}

	user, err := userDAO.FindByEmail(req.Email)
	if err != nil {
		klog.Error("登陆失败，用户不存在: ", err)
		return &user_service.LoginResp{
			ResponseStatus: buildResponse("用户不存在！", false),
		}, nil
	}

	if user == nil {
		klog.Error("登陆失败，用户不存在: ", err)
		return &user_service.LoginResp{
			ResponseStatus: buildResponse("用户不存在！", false),
		}, nil
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password+user.Salt)); err != nil {
		klog.Error("密码验证错误: ", err)
		return &user_service.LoginResp{
			ResponseStatus: buildResponse("邮箱或密码错误!", false),
		}, nil
	}

	klog.Info("用户登录成功！")
	return &user_service.LoginResp{
		ResponseStatus: &user_service.ResponseStatus{
			Message: "登陆成功",
			Status:  true,
		},
	}, nil
}
