package service

import (
	"context"
	"tiktok_e-commerce/auth/conf"
	"tiktok_e-commerce/common/constant"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"
)

type CheckIfUserBannedService struct {
	ctx context.Context
} // NewCheckIfUserBannedService new CheckIfUserBannedService
func NewCheckIfUserBannedService(ctx context.Context) *CheckIfUserBannedService {
	return &CheckIfUserBannedService{ctx: ctx}
}

// Run create note info
func (s *CheckIfUserBannedService) Run(req *auth.CheckIfUserBannedReq) (resp *auth.CheckIfUserBannedResp, err error) {
	_ , exist := conf.BannedUserList[req.UserId]
	resp = &auth.CheckIfUserBannedResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		IsBanned:   exist,
	}
	return resp, nil
}
