package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type BrandInsertService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewBrandInsertService(Context context.Context, RequestContext *app.RequestContext) *BrandInsertService {
	return &BrandInsertService{RequestContext: RequestContext, Context: Context}
}

func (h *BrandInsertService) Run(req *product.BrandInsertRequest) (resp *product.BrandInsertResponse, err error) {
	brand, err := rpc.ProductClient.InsertBrand(h.Context, &rpcproduct.BrandInsertReq{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	})
	resp = &product.BrandInsertResponse{
		StatusCode: brand.StatusCode,
		StatusMsg:  brand.StatusMsg,
	}
	return
}
