package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type CategoryDeleteService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCategoryDeleteService(Context context.Context, RequestContext *app.RequestContext) *CategoryDeleteService {
	return &CategoryDeleteService{RequestContext: RequestContext, Context: Context}
}

func (h *CategoryDeleteService) Run(req *product.CategoryDeleteRequest) (resp *product.CategoryDeleteResponse, err error) {
	category, err := rpc.ProductClient.DeleteCategory(h.Context, &rpcproduct.CategoryDeleteReq{
		CategoryId: req.CategoryId,
	})
	resp = &product.CategoryDeleteResponse{
		StatusCode: category.StatusCode,
		StatusMsg:  category.StatusMsg,
	}
	return
}
