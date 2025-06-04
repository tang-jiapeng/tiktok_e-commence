package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/rpc_gen/kitex_gen/auth"
	user "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	"tiktok_e-commerce/user/biz/dal/mysql"
	"tiktok_e-commerce/user/biz/infra/rpc"
	"tiktok_e-commerce/user/biz/model"

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
func (s *LoginService) Run(req *user.LoginReq) (resp *user.LoginResp, err error) {
	loginUser, err := model.GetUserByUsername(mysql.DB, s.ctx, req.Username)
	if err != nil {
		// 数据库错误
		klog.Error("用户登录失败,req: %v, err: %v", req, err)
		return nil, err
	}
	comparePwdErr := bcrypt.CompareHashAndPassword([]byte(loginUser.Password), []byte(req.Password))
	if loginUser == nil || comparePwdErr != nil {
		// 邮箱密码错误
		resp = &user.LoginResp{
			StatusCode: 1003,
			StatusMsg:  constant.GetMsg(1003),
		}
		return
	}
	client := rpc.AuthClient
	deliveryTokenResp, err := client.DeliverTokenByRPC(s.ctx, &auth.DeliverTokenReq{UserId: loginUser.ID})
	if err != nil {
		klog.Error("调用用户授权服务发放令牌失败，UserId: %v, err: %v", loginUser.ID, err)
		return nil, err
	}
	resp = &user.LoginResp{
		StatusCode:   0,
		StatusMsg:    constant.GetMsg(0),
		AccessToken:  deliveryTokenResp.AccessToken,
		RefreshToken: deliveryTokenResp.RefreshToken,
	}
	return
}
