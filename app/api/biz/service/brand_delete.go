package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type BrandDeleteService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewBrandDeleteService(Context context.Context, RequestContext *app.RequestContext) *BrandDeleteService {
	return &BrandDeleteService{RequestContext: RequestContext, Context: Context}
}

func (h *BrandDeleteService) Run(req *product.BrandDeleteRequest) (resp *product.BrandDeleteResponse, err error) {
	brand, err := rpc.ProductClient.DeleteBrand(h.Context, &rpcproduct.BrandDeleteReq{
		BrandId: req.BrandId,
	})
	resp = &product.BrandDeleteResponse{
		StatusCode: brand.StatusCode,
		StatusMsg:  brand.StatusMsg,
	}
	return
}
