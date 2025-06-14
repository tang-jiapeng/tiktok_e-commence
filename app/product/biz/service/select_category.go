package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/klog"
)

type SelectCategoryService struct {
	ctx context.Context
} // NewSelectCategoryService new SelectCategoryService
func NewSelectCategoryService(ctx context.Context) *SelectCategoryService {
	return &SelectCategoryService{ctx: ctx}
}

// Run create note info
func (s *SelectCategoryService) Run(req *product.CategorySelectReq) (resp *product.CategorySelectResp, err error) {
	category, err := model.SelectCategory(mysql.DB, s.ctx, req.CategoryId)
	if err != nil {
		klog.Error("select category failed, error:%v", err)
		resp = &product.CategorySelectResp{
			StatusCode: 2017,
			StatusMsg:  constant.GetMsg(2017),
		}
		return
	}

	resp = &product.CategorySelectResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		Category: &product.Category{
			Id:          category.ID,
			Name:        category.Name,
			Description: category.Description,
		},
	}
	return
}
