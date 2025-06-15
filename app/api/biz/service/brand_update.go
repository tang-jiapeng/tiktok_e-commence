package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type BrandUpdateService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewBrandUpdateService(Context context.Context, RequestContext *app.RequestContext) *BrandUpdateService {
	return &BrandUpdateService{RequestContext: RequestContext, Context: Context}
}

func (h *BrandUpdateService) Run(req *product.BrandUpdateRequest) (resp *product.BrandUpdateResponse, err error) {
	brand, err := rpc.ProductClient.UpdateBrand(h.Context, &rpcproduct.BrandUpdateReq{
		BrandId:     req.BrandId,
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	})
	resp = &product.BrandUpdateResponse{
		StatusCode: brand.StatusCode,
		StatusMsg:  brand.StatusMsg,
	}
	return
}
