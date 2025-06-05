package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	user "tiktok_e-commerce/rpc_gen/kitex_gen/user"
	"tiktok_e-commerce/user/biz/dal/mysql"

	"tiktok_e-commerce/user/biz/model"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type UpdateUserService struct {
	ctx context.Context
} // NewUpdateUserService new UpdateUserService
func NewUpdateUserService(ctx context.Context) *UpdateUserService {
	return &UpdateUserService{ctx: ctx}
}

// Run create note info
func (s *UpdateUserService) Run(req *user.UpdateUserReq) (resp *user.UpdateUserResp, err error) {
	if err = model.UpdateUser(mysql.DB, s.ctx, req.UserId, req.Username, req.Email, req.Sex, req.Description, req.Avatar); err != nil {
		klog.Errorf("更新用户信息失败: req: %v, err: %v", req, err)
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return &user.UpdateUserResp{
				StatusCode: 1008,
				StatusMsg:  constant.GetMsg(1008),
			}, nil
		}
		return nil, errors.WithStack(err)
	}
	return &user.UpdateUserResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}, nil
}
