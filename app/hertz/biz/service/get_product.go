package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	common "tiktok_e-commerce/app/hertz/hertz_gen/hertz/common"
	product "tiktok_e-commerce/app/hertz/hertz_gen/hertz/product"
)

type GetProductService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetProductService(Context context.Context, RequestContext *app.RequestContext) *GetProductService {
	return &GetProductService{RequestContext: RequestContext, Context: Context}
}

func (h *GetProductService) Run(req *product.GetProductReq) (resp *common.Empty, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
