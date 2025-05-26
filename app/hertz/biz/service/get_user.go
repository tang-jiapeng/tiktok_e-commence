package service

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	common "tiktok_e-commerce/app/hertz/hertz_gen/hertz/common"
	user "tiktok_e-commerce/app/hertz/hertz_gen/hertz/user"
)

type GetUserService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewGetUserService(Context context.Context, RequestContext *app.RequestContext) *GetUserService {
	return &GetUserService{RequestContext: RequestContext, Context: Context}
}

func (h *GetUserService) Run(req *user.GetUserReq) (resp *common.Empty, err error) {
	//defer func() {
	// hlog.CtxInfof(h.Context, "req = %+v", req)
	// hlog.CtxInfof(h.Context, "resp = %+v", resp)
	//}()
	// todo edit your code
	return
}
