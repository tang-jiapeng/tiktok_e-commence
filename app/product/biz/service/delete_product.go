package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"
)

type DeleteProductService struct {
	ctx context.Context
} // NewDeleteProductService new DeleteProductService
func NewDeleteProductService(ctx context.Context) *DeleteProductService {
	return &DeleteProductService{ctx: ctx}
}

// Run create note info
func (s *DeleteProductService) Run(req *product.DeleteProductReq) (resp *product.DeleteProductResp, err error) {
	err = model.DeleteProduct(mysql.DB, s.ctx, req.Id)
	if err != nil {
		resp = &product.DeleteProductResp{
			StatusCode: 2001,
			StatusMsg:  constant.GetMsg(2001),
		}
		return
	}
	resp = &product.DeleteProductResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
