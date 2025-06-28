package service

import (
	"context"

	"github.com/pkg/errors"

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

	client := rpc.ProductClient
	res, err := client.SearchProducts(h.Context, &rpcproduct.SearchProductsReq{
		Query: req.ProductName,
	})

	if err != nil {
		hlog.CtxErrorf(h.Context, "商品搜索失败: %v", errors.WithStack(err))
		return nil, err
	}
	var productList []*product.Product
	for i := range res.Results {
		source := res.Results[i]
		productList = append(productList, &product.Product{
			Name:         source.Name,
			Description:  source.Description,
			CategoryName: source.CategoryName,
		})
	}
	resp = &product.ProductResponse{
		StatusCode: res.StatusCode,
		StatusMsg:  res.StatusMsg,
		Products:   productList,
	}
	return resp, nil
}
