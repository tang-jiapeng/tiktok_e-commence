package payment

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"tiktok_e-commerce/api/biz/service"
	"tiktok_e-commerce/api/biz/utils"
	payment "tiktok_e-commerce/api/hertz_gen/api/payment"
)

// Charge .
// @router /payment/charge [POST]
func Charge(ctx context.Context, c *app.RequestContext) {
	var err error
	var req payment.ChargeRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &payment.ChargeResponse{}
	resp, err = service.NewChargeService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
