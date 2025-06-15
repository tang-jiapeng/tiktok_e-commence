package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpccart "tiktok_e-commerce/rpc_gen/kitex_gen/cart"

	cart "tiktok_e-commerce/api/hertz_gen/api/cart"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
)

type GetCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetCartService(Context context.Context, RequestContext *app.RequestContext) *GetCartService {
	return &GetCartService{RequestContext: RequestContext, Context: Context}
}

func (h *GetCartService) Run(req *cart.GetCartReq) (resp *cart.GetCartResp, err error) {
	getCartResp, err := rpc.CartClient.GetCart(h.Context, &rpccart.GetCartReq{
		UserId: h.Context.Value("user_id").(int32),
	})
	if err != nil {
		hlog.CtxErrorf(h.Context, "rpc调用获取购物车失败，req：%v, err: %v", req, err)
		return nil, errors.New("获取购物车失败")
	}
	items := make([]*cart.Product, len(getCartResp.Products))
	for i, item := range getCartResp.Products {
		items[i] = &cart.Product{
			Id:          item.Id,
			Name:        item.Name,
			Description: item.Description,
			Picture:     item.Picture,
			Price:       item.Price,
			Quantity:    item.Quantity,
		}
	}
	resp = &cart.GetCartResp{
		StatusCode: getCartResp.StatusCode,
		StatusMsg:  getCartResp.StatusMsg,
		Products:   items,
	}
	return
}
