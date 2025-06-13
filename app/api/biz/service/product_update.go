package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	product "tiktok_e-commerce/api/hertz_gen/api/product"
)

type ProductUpdateService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewProductUpdateService(Context context.Context, RequestContext *app.RequestContext) *ProductUpdateService {
	return &ProductUpdateService{RequestContext: RequestContext, Context: Context}
}

func (h *ProductUpdateService) Run(req *product.ProductUpdateRequest) (resp *product.ProductUpdateResponse, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
