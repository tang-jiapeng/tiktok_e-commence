package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type CategoryUpdateService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewCategoryUpdateService(Context context.Context, RequestContext *app.RequestContext) *CategoryUpdateService {
	return &CategoryUpdateService{RequestContext: RequestContext, Context: Context}
}

func (h *CategoryUpdateService) Run(req *product.CategoryUpdateRequest) (resp *product.CategoryUpdateResponse, err error) {
	category, err := rpc.ProductClient.UpdateCategory(h.Context, &rpcproduct.CategoryUpdateReq{
		CategoryId:  req.CategoryId,
		Name:        req.Name,
		Description: req.Description,
	})
	resp = &product.CategoryUpdateResponse{
		StatusCode: category.StatusCode,
		StatusMsg:  category.StatusMsg,
	}
	return
}
