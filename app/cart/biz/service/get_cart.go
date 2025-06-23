package service

import (
	"context"
	"tiktok_e-commerce/cart/biz/dal/mysql"
	"tiktok_e-commerce/cart/biz/model"
	"tiktok_e-commerce/cart/infra/rpc"
	"tiktok_e-commerce/common/constant"
	cart "tiktok_e-commerce/rpc_gen/kitex_gen/cart"
	"tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pkg/errors"
)

type GetCartService struct {
	ctx context.Context
} // NewGetCartService new GetCartService
func NewGetCartService(ctx context.Context) *GetCartService {
	return &GetCartService{ctx: ctx}
}

// Run create note info
func (s *GetCartService) Run(req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	cartItems, err := model.GetCartItemByUserId(mysql.DB, s.ctx, req.UserId)
	if err != nil {
		klog.CtxErrorf(s.ctx, "数据库查询购物车信息失败，req: %v, err: %v", req, err)
		return nil, err
	}
	if len(cartItems) == 0 {
		return &cart.GetCartResp{
			StatusCode: 0,
			StatusMsg:  constant.GetMsg(0),
			Products:   make([]*cart.Product, 0),
		}, nil
	}
	productIds := make([]int64, len(cartItems))
	for i, item := range cartItems {
		productIds[i] = int64(item.ProductId)
	}
	getProductListResp, err := rpc.ProductClient.SelectProductList(s.ctx, &product.SelectProductListReq{
		Ids: productIds,
	})
	if err != nil {
		klog.CtxErrorf(s.ctx, "rpc查询商品信息失败，req: %v, err: %v", req, err)
		return nil, errors.WithStack(err)
	}

	productMap := make(map[int]*product.Product)
	for _, p := range getProductListResp.Products {
		productMap[int(p.Id)] = p
	}

	productItems := make([]*cart.Product, len(cartItems))
	for i, item := range cartItems {
		p := productMap[int(item.ProductId)]
		if p == nil {
			// 商品不存在，返回空数据
			productItems[i] = &cart.Product{
				Id:          item.ProductId,
				Name:        "",
				Description: "",
				Picture:     "",
				Price:       0,
				Quantity:    item.Quantity,
			}
			continue
		}
		productItems[i] = &cart.Product{
			Id:          item.ProductId,
			Name:        p.Name,
			Description: p.Description,
			Picture:     p.Picture,
			Price:       p.Price,
			Quantity:    item.Quantity,
		}
	}
	resp = &cart.GetCartResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		Products:   productItems,
	}
	return
}
