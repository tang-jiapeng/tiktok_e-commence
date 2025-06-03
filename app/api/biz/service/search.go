package service

import (
	"context"
	"errors"

	product "tiktok_e-commerce/api/hertz_gen/api/product"
	"tiktok_e-commerce/api/infra/rpc"

	rpcproduct "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
)

type SearchService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewSearchService(Context context.Context, RequestContext *app.RequestContext) *SearchService {
	return &SearchService{RequestContext: RequestContext, Context: Context}
}

func (h *SearchService) Run(req *product.ProductRequest) (resp *product.ProductResponse, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	client := rpc.ProductClient
	res, err := client.SearchProducts(h.Context, &rpcproduct.SearchProductsReq{
		Query: req.ProductName,
	})
	if err != nil {
		hlog.Error("product search error: %v", err)
		return nil, errors.New("搜索失败，请稍后再试")
	}
	resp = &product.ProductResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
	}
	return
}
