package service

import (
	"context"
	"errors"

	payment "tiktok_e-commerce/api/hertz_gen/api/payment"
	"tiktok_e-commerce/api/infra/rpc"
	rpcpayment "tiktok_e-commerce/rpc_gen/kitex_gen/payment"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/kitex/pkg/klog"
)

type ChargeService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewChargeService(Context context.Context, RequestContext *app.RequestContext) *ChargeService {
	return &ChargeService{RequestContext: RequestContext, Context: Context}
}

func (h *ChargeService) Run(req *payment.ChargeRequest) (resp *payment.ChargeResponse, err error) {
	client := rpc.PaymentClient
	res, err := client.Charge(h.Context, &rpcpayment.ChargeReq{
		Amount:  req.Amount,
		OrderId: req.OrderId,
		UserId:  req.UserId,
	})
	if err != nil {
		klog.Error("payment charge error: %v", err)
		return nil, errors.New("支付失败，请稍后再试")
	}
	resp = &payment.ChargeResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
	}
	return
}
