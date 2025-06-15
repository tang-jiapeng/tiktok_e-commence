package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type ProductSelectService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewProductSelectService(Context context.Context, RequestContext *app.RequestContext) *ProductSelectService {
	return &ProductSelectService{RequestContext: RequestContext, Context: Context}
}

func (h *ProductSelectService) Run(req *product.ProductSelectRequest) (resp *product.ProductSelectResponse, err error) {
	selectProduct, err := rpc.ProductClient.SelectProduct(h.Context, &rpcproduct.SelectProductReq{Id: req.Id})
	if err != nil {
		resp = &product.ProductSelectResponse{
			StatusCode: selectProduct.StatusCode,
			StatusMsg:  selectProduct.StatusMsg,
		}
		return
	}
	resp = &product.ProductSelectResponse{
		StatusCode: selectProduct.StatusCode,
		StatusMsg:  selectProduct.StatusMsg,
		Product: &product.Product{
			Id:            selectProduct.Product.Id,
			Name:          selectProduct.Product.Name,
			Description:   selectProduct.Product.Description,
			Price:         selectProduct.Product.Price,
			Stock:         selectProduct.Product.Stock,
			Sale:          selectProduct.Product.Sale,
			PublishStatus: selectProduct.Product.PublishStatus,
			Picture:       selectProduct.Product.Picture,
			Categories:    selectProduct.Product.Categories,
			BrandId:       selectProduct.Product.BrandId,
			CategoryId:    selectProduct.Product.CategoryId,
		},
	}
	return
}
