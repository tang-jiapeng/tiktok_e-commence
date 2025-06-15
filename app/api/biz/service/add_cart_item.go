package service

import (
	"context"
	cart "tiktok_e-commerce/api/hertz_gen/api/cart"
	"tiktok_e-commerce/api/infra/rpc"
	rpccart "tiktok_e-commerce/rpc_gen/kitex_gen/cart"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
)

type AddCartItemService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddCartItemService(Context context.Context, RequestContext *app.RequestContext) *AddCartItemService {
	return &AddCartItemService{RequestContext: RequestContext, Context: Context}
}

func (h *AddCartItemService) Run(req *cart.AddCartItemReq) (resp *cart.AddCartItemResp, err error) {
	addItemResp, err := rpc.CartClient.AddItem(h.Context, &rpccart.AddItemReq{
		Item: &rpccart.CartItem{
			ProductId: req.Item.ProductId,
			Quantity:  req.Item.Quantity,
		},
		UserId: h.Context.Value("user_id").(int32),
	})
	if err != nil {
		hlog.CtxErrorf(h.Context, "rpc调用添加购物车失败，req：%v, err: %v", req, err)
		return nil, errors.New("添加购物车失败")
	}
	resp = &cart.AddCartItemResp{
		StatusCode: addItemResp.StatusCode,
		StatusMsg:  addItemResp.StatusMsg,
	}
	return
}
