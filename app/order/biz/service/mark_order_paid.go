package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/order/biz/model"
	order "tiktok_e-commerce/rpc_gen/kitex_gen/order"
	"tiktok_e-commerce/user/biz/dal/mysql"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pkg/errors"
)

type MarkOrderPaidService struct {
	ctx context.Context
} // NewMarkOrderPaidService new MarkOrderPaidService
func NewMarkOrderPaidService(ctx context.Context) *MarkOrderPaidService {
	return &MarkOrderPaidService{ctx: ctx}
}

// Run create note info
func (s *MarkOrderPaidService) Run(req *order.MarkOrderPaidReq) (resp *order.MarkOrderPaidResp, err error) {
	ctx := s.ctx
	orderId := req.OrderId
	rowAffected, err := model.MarkOrderPaid(ctx, mysql.DB, orderId)
	if err != nil {
		klog.CtxErrorf(ctx, "数据库修改订单状态失败,orderId:%s,err:%v", orderId, err)
		return nil, errors.New("数据库修改订单状态失败")
	}
	if rowAffected == 0 {
		return &order.MarkOrderPaidResp{
			StatusCode: 3000,
			StatusMsg:  constant.GetMsg(3000),
		}, nil
	}
	return &order.MarkOrderPaidResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}, nil
}
