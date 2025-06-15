package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	"tiktok_e-commerce/common/constant"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type ProductSelectListService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewProductSelectListService(Context context.Context, RequestContext *app.RequestContext) *ProductSelectListService {
	return &ProductSelectListService{RequestContext: RequestContext, Context: Context}
}

func (h *ProductSelectListService) Run(req *product.ProductSelectListRequest) (resp *product.ProductSelectListResponse, err error) {
	selectProduct, err := rpc.ProductClient.SelectProductList(h.Context, &rpcproduct.SelectProductListReq{
		Ids: req.Id,
	})
	if err != nil {
		return nil, err
	}
	var products []*product.Product
	for i := range selectProduct.Products {
		products = append(products, &product.Product{
			Id:            selectProduct.Products[i].Id,
			Name:          selectProduct.Products[i].Name,
			Description:   selectProduct.Products[i].Description,
			Price:         selectProduct.Products[i].Price,
			Stock:         selectProduct.Products[i].Stock,
			Sale:          selectProduct.Products[i].Sale,
			PublishStatus: selectProduct.Products[i].PublishStatus,
			Picture:       selectProduct.Products[i].Picture,
			Categories:    selectProduct.Products[i].Categories,
			BrandId:       selectProduct.Products[i].BrandId,
			CategoryId:    selectProduct.Products[i].CategoryId,
		})
	}
	resp = &product.ProductSelectListResponse{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		Products:   products,
	}
	return
}
