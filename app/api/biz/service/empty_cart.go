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

type EmptyCartService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewEmptyCartService(Context context.Context, RequestContext *app.RequestContext) *EmptyCartService {
	return &EmptyCartService{RequestContext: RequestContext, Context: Context}
}

func (h *EmptyCartService) Run(req *cart.EmptyCartReq) (resp *cart.EmptyCartResp, err error) {
	emptyCartResp, err := rpc.CartClient.EmptyCart(h.Context, &rpccart.EmptyCartReq{
		UserId: h.Context.Value("user_id").(int32),
	})
	if err != nil {
		hlog.CtxErrorf(h.Context, "rpc调用清空购物车失败, err: %v", err)
		return nil, errors.New("清空购物车失败")
	}
	resp = &cart.EmptyCartResp{
		StatusCode: emptyCartResp.StatusCode,
		StatusMsg:  emptyCartResp.StatusMsg,
	}
	return
}
