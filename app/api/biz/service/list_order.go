package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"

	order "tiktok_e-commerce/api/hertz_gen/api/order"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"

	rpcOrder "tiktok_e-commerce/rpc_gen/kitex_gen/order"
)

type ListOrderService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewListOrderService(Context context.Context, RequestContext *app.RequestContext) *ListOrderService {
	return &ListOrderService{RequestContext: RequestContext, Context: Context}
}

func (h *ListOrderService) Run(req *order.ListOrderRequest) (resp *order.ListOrderResponse, err error) {
	ctx := h.Context
	listOrderResp, err := rpc.OrderClient.ListOrder(ctx, &rpcOrder.ListOrderReq{
		UserId: ctx.Value("user_id").(int32),
	})
	if err != nil {
		hlog.CtxErrorf(ctx, "rpc调用获取订单列表失败, err: %v", err)
		return nil, err
	}
	var orders []*order.Order
	for _, o := range listOrderResp.Orders {
		var products []*order.Product
		for _, p := range o.Products {
			products = append(products, &order.Product{
				Id:          p.Id,
				Name:        p.Name,
				Description: p.Description,
				Picture:     p.Picture,
				Price:       p.Price,
				Quantity:    p.Quantity,
			})
		}
		orders = append(orders, &order.Order{
			OrderId: o.OrderId,
			Address: &order.Address{
				Name:          o.Address.Name,
				PhoneName:     o.Address.PhoneNumber,
				Province:      o.Address.Province,
				City:          o.Address.City,
				Region:        o.Address.Region,
				DetailAddress: o.Address.DetailAddress,
			},
			Products:  products,
			Cost:      o.Cost,
			Status:    o.Status,
			CreatedAt: o.CreatedAt,
		})
	}
	return &order.ListOrderResponse{
		StatusCode: listOrderResp.StatusCode,
		StatusMsg:  listOrderResp.StatusMsg,
		Orders:     orders,
	}, nil
}
