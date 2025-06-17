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

type InsertProductService struct {
	ctx context.Context
} // NewInsertProductService new InsertProductService
func NewInsertProductService(ctx context.Context) *InsertProductService {
	return &InsertProductService{ctx: ctx}
}

// Run create note info
func (s *InsertProductService) Run(req *product.InsertProductReq) (resp *product.InsertProductResp, err error) {
	pro := model.Product{
		Name:        req.Name,
		Description: req.Description,
		Picture:     req.Picture,
		Price:       req.Price,
		Stock:       req.Stock,
		Sale:        0,
		PublicState: 1,
		LockStock:   req.Stock,
		CategoryId:  req.CategoryId,
		BrandId:     req.BrandId,
	}

	err = model.CreateProduct(mysql.DB, s.ctx, &pro)
	if err != nil {
		klog.Error("insert product error:%v", err)
		resp = &product.InsertProductResp{
			StatusCode: 2002,
			StatusMsg:  constant.GetMsg(2002),
		}
		return
	}

	//发送到kafka
	defer func() {
		err := kf.AddProduct(&pro)
		if err != nil {
			klog.Error("insert product error:%v", err)
		}
	}()

	resp = &product.InsertProductResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}
	return
}
