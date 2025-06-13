package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	product "tiktok_e-commerce/api/hertz_gen/api/product"
)

type ProductInsertService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewProductInsertService(Context context.Context, RequestContext *app.RequestContext) *ProductInsertService {
	return &ProductInsertService{RequestContext: RequestContext, Context: Context}
}

func (h *ProductInsertService) Run(req *product.ProductInsertRequest) (resp *product.ProductInsertResponse, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
