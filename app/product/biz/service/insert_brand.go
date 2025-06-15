package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/klog"
)

type InsertBrandService struct {
	ctx context.Context
} // NewInsertBrandService new InsertBrandService
func NewInsertBrandService(ctx context.Context) *InsertBrandService {
	return &InsertBrandService{ctx: ctx}
}

// Run create note info
func (s *InsertBrandService) Run(req *product.BrandInsertReq) (resp *product.BrandInsertResp, err error) {
	brand := model.Brand{
		Name:        req.Name,
		Description: req.Description,
		Icon:        req.Icon,
	}
	err = model.CreateBrand(mysql.DB, s.ctx, &brand)
	if err != nil {
		klog.Error("insert brand failed, error:%v", err)
		resp = &product.BrandInsertResp{
			StatusCode: 2018,
			StatusMsg:  constant.GetMsg(2018),
		}
		return
	}
	resp = &product.BrandInsertResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
