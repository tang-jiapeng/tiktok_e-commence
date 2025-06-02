package user

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"tiktok_e-commerce/api/biz/service"
	"tiktok_e-commerce/api/biz/utils"
	user "tiktok_e-commerce/api/hertz_gen/api/user"
)

// Login .
// @router /user/login [GET]
func Login(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.LoginRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &user.LoginResponse{}
	resp, err = service.NewLoginService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// Register .
// @router /user/register [POST]
func Register(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.RegisterRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &user.RegisterResponse{}
	resp, err = service.NewRegisterService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// GetUserInfo .
// @router /user/info [GET]
func GetUserInfo(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.Empty
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &user.GetUserInfoResponse{}
	resp, err = service.NewGetUserInfoService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}

// RefreshToken .
// @router /user/refresh_token [POST]
func RefreshToken(ctx context.Context, c *app.RequestContext) {
	var err error
	var req user.Empty
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp, err := service.NewRefreshTokenService(ctx, c).Run(&req)

	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}
	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
