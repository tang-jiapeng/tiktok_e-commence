package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type ProductInsertService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewProductInsertService(Context context.Context, RequestContext *app.RequestContext) *ProductInsertService {
	return &ProductInsertService{RequestContext: RequestContext, Context: Context}
}

func (h *ProductInsertService) Run(req *product.ProductInsertRequest) (resp *product.ProductInsertResponse, err error) {
	productReq := rpcproduct.InsertProductReq{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Stock:       req.Stock,
	}
	insertProduct, err := rpc.ProductClient.InsertProduct(h.Context, &productReq)
	resp = &product.ProductInsertResponse{
		StatusCode: insertProduct.StatusCode,
		StatusMsg:  insertProduct.StatusMsg,
	}
	return
}
