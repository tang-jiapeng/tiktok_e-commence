package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type ProductUpdateService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewProductUpdateService(Context context.Context, RequestContext *app.RequestContext) *ProductUpdateService {
	return &ProductUpdateService{RequestContext: RequestContext, Context: Context}
}

func (h *ProductUpdateService) Run(req *product.ProductUpdateRequest) (resp *product.ProductUpdateResponse, err error) {
	updateProduct, err := rpc.ProductClient.UpdateProduct(h.Context, &rpcproduct.UpdateProductReq{
		Id:            req.Id,
		Name:          req.Name,
		Description:   req.Description,
		Picture:       req.Picture,
		Price:         req.Price,
		Categories:    req.Categories,
		Stock:         req.Stock,
		Sale:          req.Sale,
		PublishStatus: req.PublishStatus,
	})
	resp = &product.ProductUpdateResponse{
		StatusCode: updateProduct.StatusCode,
		StatusMsg:  updateProduct.StatusMsg,
	}
	return
}
