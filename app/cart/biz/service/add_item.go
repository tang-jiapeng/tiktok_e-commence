package service

import (
	"context"
	"tiktok_e-commerce/cart/biz/dal/mysql"
	"tiktok_e-commerce/cart/biz/model"
	"tiktok_e-commerce/common/constant"
	cart "tiktok_e-commerce/rpc_gen/kitex_gen/cart"

	"github.com/cloudwego/kitex/pkg/klog"
)

type AddItemService struct {
	ctx context.Context
} // NewAddItemService new AddItemService
func NewAddItemService(ctx context.Context) *AddItemService {
	return &AddItemService{ctx: ctx}
}

// Run create note info
func (s *AddItemService) Run(req *cart.AddItemReq) (resp *cart.AddItemResp, err error) {
	err = model.AddCartItem(mysql.DB, s.ctx, &model.CartItem{
		UserId:    req.UserId,
		ProductId: req.Item.ProductId,
		Quantity:  req.Item.Quantity,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "add item failed %v", err)
		return nil, err
	}
	return &cart.AddItemResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
	}, nil
}
