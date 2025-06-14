package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"
)

type UpdateCategoryService struct {
	ctx context.Context
} // NewUpdateCategoryService new UpdateCategoryService
func NewUpdateCategoryService(ctx context.Context) *UpdateCategoryService {
	return &UpdateCategoryService{ctx: ctx}
}

// Run create note info
func (s *UpdateCategoryService) Run(req *product.CategoryUpdateReq) (resp *product.CategoryUpdateResp, err error) {
	err = model.UpdateCategory(mysql.DB, s.ctx, &model.Category{
		Base: model.Base{
			ID: req.CategoryId,
		},
		Name:        req.Name,
		Description: req.Description,
	})
	if err != nil {
		resp = &product.CategoryUpdateResp{
			StatusCode: 2016,
			StatusMsg:  constant.GetMsg(2016),
		}
		return
	}
	resp = &product.CategoryUpdateResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
