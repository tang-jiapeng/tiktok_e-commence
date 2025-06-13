package middleware

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/rpc_gen/kitex_gen/auth"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

var (
	Whitelist = map[string]struct{}{
		"/ping":               {},
		"/user/register":      {},
		"/user/login":         {},
		"/user/refresh_token": {},
		"/product/search":		 {},
	}
)

func AuthorizationMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())
		if _, exist := Whitelist[path]; !exist {
			authClient := rpc.AuthClient
			verifyResp, err := authClient.VerifyTokenByRPC(ctx, &auth.VerifyTokenReq{
				RefreshToken: c.Request.Header.Get("refresh_token"),
				AccessToken:  c.Request.Header.Get("access_token"),
				Path:         path,
				Method:       string(c.Request.Method()),
			})
			if err != nil {
				hlog.CtxErrorf(ctx, "rpc权限校验失败，err: %v", err)
				c.JSON(consts.StatusOK, utils.H{
					"status_code": 500,
					"status_msg":  constant.GetMsg(500)})
				c.Abort()
			} else {
				if verifyResp.StatusCode != 0 {
					c.JSON(consts.StatusOK, verifyResp)
					c.Abort()
				}
				ctx = context.WithValue(ctx, "user_id", verifyResp.UserId)
			}
		}
		c.Next(ctx)
	}
}
