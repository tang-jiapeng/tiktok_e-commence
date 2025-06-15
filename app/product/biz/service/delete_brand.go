package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"
)

type DeleteBrandService struct {
	ctx context.Context
} // NewDeleteBrandService new DeleteBrandService
func NewDeleteBrandService(ctx context.Context) *DeleteBrandService {
	return &DeleteBrandService{ctx: ctx}
}

// Run create note info
func (s *DeleteBrandService) Run(req *product.BrandDeleteReq) (resp *product.BrandDeleteResp, err error) {
	err = model.DeleteBrand(mysql.DB, s.ctx, req.BrandId)
	if err != nil {
		resp = &product.BrandDeleteResp{
			StatusCode: 2019,
			StatusMsg:  constant.GetMsg(2019),
		}
		return
	}
	resp = &product.BrandDeleteResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
