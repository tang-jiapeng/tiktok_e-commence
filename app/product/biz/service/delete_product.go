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
	//发送到kafka
	defer func() {
		err := kf.DeleteProduct(&model.Product{
			Base: model.Base{
				ID: req.Id,
			},
		})
		if err != nil {
			klog.Error("delete product error:%v", err)
		}
	}()
	resp = &product.DeleteProductResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
