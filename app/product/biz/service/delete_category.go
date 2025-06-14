package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/klog"
)

type DeleteCategoryService struct {
	ctx context.Context
} // NewDeleteCategoryService new DeleteCategoryService
func NewDeleteCategoryService(ctx context.Context) *DeleteCategoryService {
	return &DeleteCategoryService{ctx: ctx}
}

// Run create note info
func (s *DeleteCategoryService) Run(req *product.CategoryDeleteReq) (resp *product.CategoryDeleteResp, err error) {
	err = model.DeleteCategory(mysql.DB, s.ctx, req.CategoryId)
	if err != nil {
		klog.Error("delete category failed, error:%v", err)
		resp = &product.CategoryDeleteResp{
			StatusCode: 2015,
			StatusMsg:  constant.GetMsg(2015),
		}
		return
	}
	resp = &product.CategoryDeleteResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
