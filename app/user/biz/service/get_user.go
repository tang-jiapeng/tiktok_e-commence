package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	user "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	"tiktok_e-commerce/user/biz/dal/mysql"
	"tiktok_e-commerce/user/biz/model"
)

type GetUserService struct {
	ctx context.Context
} // NewGetUserService new GetUserService
func NewGetUserService(ctx context.Context) *GetUserService {
	return &GetUserService{ctx: ctx}
}

// Run create note info
func (s *GetUserService) Run(req *user.GetUserReq) (resp *user.GetUserResp, err error) {
	userInfo, err := model.GetUserById(mysql.DB, s.ctx, int32(req.UserId))
	if err != nil {
		return nil, err
	}
	resp = &user.GetUserResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		User: &user.User{
			Username: userInfo.Username,
			Email: userInfo.Email,
			Sex: model.SexToString(userInfo.Sex),
			Description: userInfo.Description,
			Avatar: userInfo.Avatar,
			CreatedAt: userInfo.CreatedAt.Format("2006-01-02 15:04:05"),
		},
	}
	return
}
