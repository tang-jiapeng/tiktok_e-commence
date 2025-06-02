package middleware

import (
	"context"
	"tiktok_e-commerce/api/infra/rpc"
	"tiktok_e-commerce/common/constant"
	"tiktok_e-commerce/rpc_gen/kitex_gen/auth"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/common/utils"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
)

func AuthorizationMiddleware() app.HandlerFunc {
	return func(ctx context.Context, c *app.RequestContext) {
		path := string(c.Request.URI().Path())
		if !(path == "/ping" || path == "/user/login" || path == "/user/register" || path == "user/refresh_token") {
			authClient := rpc.AuthClient
			verifyResp, err := authClient.VerifyTokenByRPC(ctx, &auth.VerifyTokenReq{
				RefreshToken: c.Request.Header.Get("refresh_token"),
				AccessToken:  c.Request.Header.Get("access_token"),
			})
			if err != nil {
				c.JSON(consts.StatusOK, utils.H{
					"status_code": 500,
					"status_msg":  constant.GetMsg(500),
				})
			} else {
				if verifyResp.StatusCode != 0 {
					c.JSON(consts.StatusOK, verifyResp)
				}
				ctx = context.WithValue(ctx, "user_id", verifyResp.UserId)
			}
		}
		c.Next(ctx)
	}
}
