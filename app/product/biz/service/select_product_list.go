package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/klog"
)

type SelectProductListService struct {
	ctx context.Context
} // NewSelectProductListService new SelectProductListService
func NewSelectProductListService(ctx context.Context) *SelectProductListService {
	return &SelectProductListService{ctx: ctx}
}

// Run create note info
func (s *SelectProductListService) Run(req *product.SelectProductListReq) (resp *product.SelectProductListResp, err error) {
	products, err := model.SelectProductList(mysql.DB, s.ctx, req.Ids)
	if err != nil {
		klog.CtxErrorf(s.ctx, "查询商品列表失败, error:%v", err)
		resp = &product.SelectProductListResp{
			StatusCode: 2003,
			StatusMsg:  constant.GetMsg(2003),
		}
		return
	}

	var productList []*product.Product
	for i := range products {
		productList = append(productList, &product.Product{
			Id:           products[i].ProductId,
			Name:         products[i].ProductName,
			Description:  products[i].ProductDescription,
			Picture:      products[i].ProductPicture,
			Price:        products[i].ProductPrice,
			CategoryName: products[i].CategoryName,
			CategoryId:   products[i].CategoryID,
		})
	}
	resp = &product.SelectProductListResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		Products:   productList,
	}
	return
}
