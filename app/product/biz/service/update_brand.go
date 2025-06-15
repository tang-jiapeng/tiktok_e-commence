package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"
)

type UpdateBrandService struct {
	ctx context.Context
} // NewUpdateBrandService new UpdateBrandService
func NewUpdateBrandService(ctx context.Context) *UpdateBrandService {
	return &UpdateBrandService{ctx: ctx}
}

// Run create note info
func (s *UpdateBrandService) Run(req *product.BrandUpdateReq) (resp *product.BrandUpdateResp, err error) {
	err = model.UpdateBrand(mysql.DB, s.ctx, &model.Brand{
		Base: model.Base{
			ID: req.BrandId,
		},
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	})
	if err != nil {
		resp = &product.BrandUpdateResp{
			StatusCode: 2020,
			StatusMsg:  constant.GetMsg(2020),
		}
		return
	}
	resp = &product.BrandUpdateResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
