package service

import (
	"context"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/common/utils"
	"tiktok_e-commerce/order/biz/dal/mysql"
	"tiktok_e-commerce/order/biz/model"
	"tiktok_e-commerce/order/infra/rpc"
	order "tiktok_e-commerce/rpc_gen/kitex_gen/order"
	"tiktok_e-commerce/rpc_gen/kitex_gen/product"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pkg/errors"
)

type GetOrderService struct {
	ctx context.Context
} // NewGetOrderService new GetOrderService
func NewGetOrderService(ctx context.Context) *GetOrderService {
	return &GetOrderService{ctx: ctx}
}

// Run create note info
func (s *GetOrderService) Run(req *order.GetOrderReq) (resp *order.GetOrderResp, err error) {
	ctx := s.ctx
	o, err := model.GetOrder(ctx, mysql.DB, req.UserId, req.OrderId)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	orderIdList := []string{req.OrderId}
	orderItems, err := model.GetOrderItemsByOrderIdList(ctx, mysql.DB, orderIdList)
	if err != nil {
		klog.CtxErrorf(ctx, "数据库查询订单商品信息失败, error: %v", err)
		return nil, errors.WithStack(err)
	}
	productIdList := make([]int64, len(orderItems))
	for i, item := range orderItems {
		productIdList[i] = int64(item.ProductID)
	}
	productListReq := &product.SelectProductListReq{
		Ids: productIdList,
	}
	getProductListResp, err := rpc.ProductClient.SelectProductList(ctx, productListReq)
	if err != nil {
		klog.CtxErrorf(ctx, "rpc查询商品信息失败, req: %v, error: %v", productListReq, err)
		return nil, errors.WithStack(err)
	}
	productMap := make(map[int]*product.Product)
	for _, p := range getProductListResp.Products {
		productMap[int(p.Id)] = p
	}

	var products []*order.Product
	for _, item := range orderItems {
		p := productMap[int(item.ProductID)]
		if p == nil {
			continue
		}
		products = append(products, &order.Product{
			Id:       int32(p.Id),
			Name:     p.Name,
			Price:    p.Price,
			Quantity: item.Quantity,
		})
	}
	orderInfo := &order.Order{
		OrderId: o.OrderID,
		Address: &order.Address{
			Name:          o.Name,
			PhoneNumber:   o.PhoneNumber,
			Province:      o.Province,
			City:          o.City,
			Region:        o.Region,
			DetailAddress: o.DetailAddress,
		},
		Products:  products,
		Cost:      o.TotalCost,
		CreatedAt: utils.GetFormattedDateTime(o.CreatedAt),
	}

	return &order.GetOrderResp{
		StatusCode: 0,
		StatusMsg:  constant.GetMsg(0),
		Order:      orderInfo,
	}, nil
}
