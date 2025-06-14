package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"

	"github.com/cloudwego/hertz/pkg/app"
	product "tiktok_e-commerce/api/hertz_gen/api/product"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"
)

type CategoryInsertService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCategoryInsertService(Context context.Context, RequestContext *app.RequestContext) *CategoryInsertService {
	return &CategoryInsertService{RequestContext: RequestContext, Context: Context}
}

func (h *CategoryInsertService) Run(req *product.CategoryInsertRequest) (resp *product.CategoryInsertResponse, err error) {
	category, err := rpc.ProductClient.InsertCategory(h.Context, &rpcproduct.CategoryInsertReq{
		Name:        req.Name,
		Description: req.Description,
	})
	resp = &product.CategoryInsertResponse{
		StatusCode: category.StatusCode,
		StatusMsg:  category.StatusMsg,
	}
	return
}
