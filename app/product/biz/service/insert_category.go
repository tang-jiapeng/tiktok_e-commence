package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/klog"
)

type InsertCategoryService struct {
	ctx context.Context
} // NewInsertCategoryService new InsertCategoryService
func NewInsertCategoryService(ctx context.Context) *InsertCategoryService {
	return &InsertCategoryService{ctx: ctx}
}

// Run create note info
func (s *InsertCategoryService) Run(req *product.CategoryInsertReq) (resp *product.CategoryInsertResp, err error) {
	category := model.Category{
		Name:        req.Name,
		Description: req.Description,
	}
	err = model.CreateCategory(mysql.DB, s.ctx, &category)
	if err != nil {
		klog.Error("insert category failed, error:%v", err)
		resp = &product.CategoryInsertResp{
			StatusCode: 2014,
			StatusMsg:  constant.GetMsg(2014),
		}
		return
	}
	resp = &product.CategoryInsertResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
