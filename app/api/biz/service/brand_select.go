package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type BrandSelectService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewBrandSelectService(Context context.Context, RequestContext *app.RequestContext) *BrandSelectService {
	return &BrandSelectService{RequestContext: RequestContext, Context: Context}
}

func (h *BrandSelectService) Run(req *product.BrandSelectRequest) (resp *product.BrandSelectResponse, err error) {
	brand, err := rpc.ProductClient.SelectBrand(h.Context, &rpcproduct.BrandSelectReq{
		BrandId: req.BrandId,
	})
	if err != nil {
		resp = &product.BrandSelectResponse{
			StatusCode: brand.StatusCode,
			StatusMsg:  brand.StatusMsg,
		}
		return
	}
	resp = &product.BrandSelectResponse{
		StatusCode: brand.StatusCode,
		StatusMsg:  brand.StatusMsg,
		Brand: &product.Brand{
			Id:          brand.Brand.Id,
			Name:        brand.Brand.Name,
			Description: brand.Brand.Description,
			Icon:        brand.Brand.Icon,
		},
	}
	return
}
