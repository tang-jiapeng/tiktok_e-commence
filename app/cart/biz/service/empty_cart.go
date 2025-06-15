package service

import (
	"context"
	"tiktok_e-commerce/cart/biz/dal/mysql"
	"tiktok_e-commerce/cart/biz/model"
	"tiktok_e-commerce/common/constant"
	cart "tiktok_e-commerce/rpc_gen/kitex_gen/cart"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pkg/errors"
)

type EmptyCartService struct {
	ctx context.Context
} // NewEmptyCartService new EmptyCartService
func NewEmptyCartService(ctx context.Context) *EmptyCartService {
	return &EmptyCartService{ctx: ctx}
}

// Run create note info
func (s *EmptyCartService) Run(req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	err = model.EmptyCart(mysql.DB, s.ctx, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "数据库操作清空购物车失败, userId: %d, err: %v", req.UserId, err)
		return nil, errors.WithStack(err)
	}
	resp = &cart.EmptyCartResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
