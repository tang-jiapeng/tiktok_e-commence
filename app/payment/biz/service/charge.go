package service

import (
	"context"
	"strconv"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/payment/biz/dal/alipay"
	payment "tiktok_e-commerce/rpc_gen/kitex_gen/payment"

	"github.com/cloudwego/kitex/pkg/klog"
)

type ChargeService struct {
	ctx context.Context
} // NewChargeService new ChargeService
func NewChargeService(ctx context.Context) *ChargeService {
	return &ChargeService{ctx: ctx}
}

// Run create note info
func (s *ChargeService) Run(req *payment.ChargeReq) (resp *payment.ChargeResp, err error) {
	orderId, err := strconv.ParseInt(req.OrderId, 0, 64)
	if err != nil {
		klog.Error("parse order id error: %v", err)
	}
	amount := req.Amount
	paymentUrl, err := alipay.Pay(s.ctx, orderId, amount)
	if err != nil {
		klog.Error("alipay pay error: %v", err)
		return &payment.ChargeResp{
			StatusCode: 5000,
			StatusMsg:  constant.GetMsg(5000),
			PaymentUrl: "",
		}, err
	}
	return &payment.ChargeResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		PaymentUrl: paymentUrl,
	}, nil
}
