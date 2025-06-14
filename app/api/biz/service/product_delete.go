package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type ProductDeleteService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewProductDeleteService(Context context.Context, RequestContext *app.RequestContext) *ProductDeleteService {
	return &ProductDeleteService{RequestContext: RequestContext, Context: Context}
}

func (h *ProductDeleteService) Run(req *product.ProductDeleteRequest) (resp *product.ProductDeleteResponse, err error) {
	deleteProduct, err := rpc.ProductClient.DeleteProduct(h.Context, &rpcproduct.DeleteProductReq{
		Id: req.Id,
	})
	resp = &product.ProductDeleteResponse{
		StatusCode: deleteProduct.StatusCode,
		StatusMsg:  deleteProduct.StatusMsg,
	}
	return
}
