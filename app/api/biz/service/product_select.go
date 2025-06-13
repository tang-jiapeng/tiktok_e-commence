package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	product "tiktok_e-commerce/api/hertz_gen/api/product"
)

type ProductSelectService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewProductSelectService(Context context.Context, RequestContext *app.RequestContext) *ProductSelectService {
	return &ProductSelectService{RequestContext: RequestContext, Context: Context}
}

func (h *ProductSelectService) Run(req *product.ProductSelectRequest) (resp *product.ProductSelectResponse, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
