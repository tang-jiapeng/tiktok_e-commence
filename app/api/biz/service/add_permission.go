package service

import (
	"context"

	user "tiktok_e-commerce/api/hertz_gen/api/user"
	"tiktok_e-commerce/api/infra/rpc"
	"tiktok_e-commerce/rpc_gen/kitex_gen/auth"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/pkg/errors"
)

type AddPermissionService struct {
	RequestContext *app.RequestContext
	Context        context.Context
}

func NewAddPermissionService(Context context.Context, RequestContext *app.RequestContext) *AddPermissionService {
	return &AddPermissionService{RequestContext: RequestContext, Context: Context}
}

func (h *AddPermissionService) Run(req *user.AddPermissionRequest) (resp *user.Empty, err error) {
	_, err = rpc.AuthClient.AddPermission(h.Context, &auth.AddPermissionReq{
		Role:   req.Role,
		Path:   req.Path,
		Method: req.Method,
	})
	if err != nil {
		hlog.Errorf("添加权限失败，req: %v, err: %v", req, err)
		return nil, errors.WithStack(err)
	}
	return &user.Empty{}, nil
}
