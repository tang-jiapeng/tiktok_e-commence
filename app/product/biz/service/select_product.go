package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/klog"
)

type SelectProductService struct {
	ctx context.Context
} // NewSelectProductService new SelectProductService
func NewSelectProductService(ctx context.Context) *SelectProductService {
	return &SelectProductService{ctx: ctx}
}

// Run create note info
func (s *SelectProductService) Run(req *product.SelectProductReq) (resp *product.SelectProductResp, err error) {
	pro, err := model.SelectProduct(mysql.DB, s.ctx, req.Id)
	if err != nil {
		klog.CtxErrorf(s.ctx, "查询商品失败, error:%v", err)
		resp = &product.SelectProductResp{
			StatusCode: 2003,
			StatusMsg:  constant.GetMsg(2003),
		}
		return
	}

	return &product.SelectProductResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		Product: &product.Product{
			Id:            pro.ID,
			Name:          pro.Name,
			Description:   pro.Description,
			Picture:       pro.Picture,
			Price:         pro.Price,
			Stock:         pro.Stock,
			Sale:          pro.Sale,
			PublishStatus: pro.PublicState,
			CategoryId:    pro.CategoryId,
			BrandId:       pro.BrandId,
			CategoryName:  pro.Category.Name,
		},
	}, nil
}
