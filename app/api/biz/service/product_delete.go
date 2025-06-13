package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	product "tiktok_e-commerce/api/hertz_gen/api/product"
)

type ProductDeleteService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewProductDeleteService(Context context.Context, RequestContext *app.RequestContext) *ProductDeleteService {
	return &ProductDeleteService{RequestContext: RequestContext, Context: Context}
}

func (h *ProductDeleteService) Run(req *product.ProductDeleteRequest) (resp *product.ProductDeleteResponse, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
