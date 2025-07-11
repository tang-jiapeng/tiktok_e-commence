package main

import (
	"context"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"
	"tiktok_e-commerce/auth/biz/service"
)

// AuthServiceImpl implements the last service interface defined in the IDL.
type AuthServiceImpl struct{}

// DeliverTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) DeliverTokenByRPC(ctx context.Context, req *auth.DeliverTokenReq) (resp *auth.DeliveryResp, err error) {
	resp, err = service.NewDeliverTokenByRPCService(ctx).Run(req)

	return resp, err
}

// VerifyTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) VerifyTokenByRPC(ctx context.Context, req *auth.VerifyTokenReq) (resp *auth.VerifyResp, err error) {
	resp, err = service.NewVerifyTokenByRPCService(ctx).Run(req)

	return resp, err
}

// RefreshTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RefreshTokenByRPC(ctx context.Context, req *auth.RefreshTokenReq) (resp *auth.RefreshTokenResp, err error) {
	resp, err = service.NewRefreshTokenByRPCService(ctx).Run(req)

	return resp, err
}

// RevokeTokenByRPC implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) RevokeTokenByRPC(ctx context.Context, req *auth.RevokeTokenReq) (resp *auth.RevokeResp, err error) {
	resp, err = service.NewRevokeTokenByRPCService(ctx).Run(req)

	return resp, err
}

// AddPermission implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) AddPermission(ctx context.Context, req *auth.AddPermissionReq) (resp *auth.Empty, err error) {
	resp, err = service.NewAddPermissionService(ctx).Run(req)

	return resp, err
}

// CheckIfUserBanned implements the AuthServiceImpl interface.
func (s *AuthServiceImpl) CheckIfUserBanned(ctx context.Context, req *auth.CheckIfUserBannedReq) (resp *auth.CheckIfUserBannedResp, err error) {
	resp, err = service.NewCheckIfUserBannedService(ctx).Run(req)

	return resp, err
}
