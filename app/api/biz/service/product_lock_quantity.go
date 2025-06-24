package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	product "tiktok_e-commerce/api/hertz_gen/api/product"

	"github.com/cloudwego/hertz/pkg/app"
)

type ProductLockQuantityService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewProductLockQuantityService(Context context.Context, RequestContext *app.RequestContext) *ProductLockQuantityService {
	return &ProductLockQuantityService{RequestContext: RequestContext, Context: Context}
}

func (h *ProductLockQuantityService) Run(req *product.ProductLockQuantityRequest) (resp *product.ProductLockQuantityResponse, err error) {
	var pro = make([]*rpcproduct.ProductLockQuantity, 0)
	originPro := req.Products
	for i := range originPro {
		pro = append(pro, &rpcproduct.ProductLockQuantity{
			Id:       originPro[i].Id,
			Quantity: originPro[i].Quantity,
		})
	}
	quantity, err := rpc.ProductClient.LockProductQuantity(h.Context,
		&rpcproduct.ProductLockQuantityRequest{Products: pro})
	resp = &product.ProductLockQuantityResponse{
		StatusCode: quantity.StatusCode,
		StatusMsg:  quantity.StatusMsg,
	}
	return
}
