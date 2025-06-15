package service

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	rpcuser "tiktok_e-commerce/rpc_gen/kitex_gen/user"

	user "tiktok_e-commerce/api/hertz_gen/api/user"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
)

type AddReceiveAddressService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddReceiveAddressService(Context context.Context, RequestContext *app.RequestContext) *AddReceiveAddressService {
	return &AddReceiveAddressService{RequestContext: RequestContext, Context: Context}
}

func (h *AddReceiveAddressService) Run(req *user.AddReceiveAddressRequest) (resp *user.AddReceiveAddressResponse, err error) {
	addReceiveAddressResp, err := rpc.UserClient.AddReceiveAddress(h.Context, &rpcuser.AddReceiveAddressReq{
		UserId:        h.Context.Value("user_id").(int32),
		Name:          req.Name,
		PhoneNumber:   req.PhoneNumber,
		DefaltStatus:  req.DefaultStatus,
		Province:      req.Province,
		City:          req.City,
		Region:        req.Region,
		DetailAddress: req.DetailAddress,
	})
	if err != nil {
		hlog.CtxErrorf(h.Context, "添加收货地址错误，req：%v，err：%v", req, err)
		return nil, errors.New("添加收货地址错误，请稍后再试")
	}
	resp = &user.AddReceiveAddressResponse{
		StatusCode: addReceiveAddressResp.StatusCode,
		StatusMsg:  addReceiveAddressResp.StatusMsg,
	}
	return
}
