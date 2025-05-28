package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	common "tiktok_e-commerce/app/hertz/hertz_gen/hertz/common"
	product "tiktok_e-commerce/app/hertz/hertz_gen/hertz/product"
)

type GetProductListService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetProductListService(Context context.Context, RequestContext *app.RequestContext) *GetProductListService {
	return &GetProductListService{RequestContext: RequestContext, Context: Context}
}

func (h *GetProductListService) Run(req *product.GetProductListReq) (resp *common.Empty, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
