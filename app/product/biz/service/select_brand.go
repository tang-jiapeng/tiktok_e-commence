package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/klog"
)

type SelectBrandService struct {
	ctx context.Context
} // NewSelectBrandService new SelectBrandService
func NewSelectBrandService(ctx context.Context) *SelectBrandService {
	return &SelectBrandService{ctx: ctx}
}

// Run create note info
func (s *SelectBrandService) Run(req *product.BrandSelectReq) (resp *product.BrandSelectResp, err error) {
	brand, err := model.SelectBrand(mysql.DB, s.ctx, req.BrandId)
	if err != nil {
		klog.Error("select brand failed, error:%v", err)
		resp = &product.BrandSelectResp{
			StatusCode: 2021,
			StatusMsg:  constant.GetMsg(2021),
		}
		return
	}
	resp = &product.BrandSelectResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		Brand: &product.Brand{
			Id:          brand.ID,
			Name:        brand.Name,
			Description: brand.Description,
		},
	}
	return
}
