package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/model"
	kf "tiktok_e-commerce/product/infra/kafka"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/klog"
)

type UpdateProductService struct {
	ctx context.Context
} // NewUpdateProductService new UpdateProductService
func NewUpdateProductService(ctx context.Context) *UpdateProductService {
	return &UpdateProductService{ctx: ctx}
}

// Run create note info
func (s *UpdateProductService) Run(req *product.UpdateProductReq) (resp *product.UpdateProductResp, err error) {
	pro := model.Product{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Stock:       req.Stock,
		Sale:        req.Sale,
		PublicState: req.PublishStatus,
		LockStock:   req.Stock,
		CategoryId:  req.CategoryId,
		BrandId:     req.BrandId,
	}
	err = model.UpdateProduct(mysql.DB, s.ctx, &pro)
	if err != nil {
		resp = &product.UpdateProductResp{
			StatusCode: 2000,
			StatusMsg:  constant.GetMsg(2000),
		}
		return
	}
	// 发送到kafka
	go func() {
		err := kf.UpdateProduct(&pro)
		if err != nil {
			klog.Error("update product error:%v", err)
		}
	}()
	resp = &product.UpdateProductResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
