package auth

import (
	"context"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"

	"tiktok_e-commerce/rpc_gen/kitex_gen/auth/authservice"
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/client/callopt"
)

type RPCClient interface {
	KitexClient() authservice.Client
	Service() string
	DeliverTokenByRPC(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryResp, err error)
	VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyResp, err error)
	RefreshTokenByRPC(ctx context.Context, Req *auth.RefreshTokenReq, callOptions ...callopt.Option) (r *auth.RefreshTokenResp, err error)
	RevokeTokenByRPC(ctx context.Context, Req *auth.RevokeTokenReq, callOptions ...callopt.Option) (r *auth.RevokeResp, err error)
	AddPermission(ctx context.Context, Req *auth.AddPermissionReq, callOptions ...callopt.Option) (r *auth.Empty, err error)
	CheckIfUserBanned(ctx context.Context, Req *auth.CheckIfUserBannedReq, callOptions ...callopt.Option) (r *auth.CheckIfUserBannedResp, err error)
}

func NewRPCClient(dstService string, opts ...client.Option) (RPCClient, error) {
	kitexClient, err := authservice.NewClient(dstService, opts...)
	if err != nil {
		return nil, err
	}
	cli := &clientImpl{
		service:     dstService,
		kitexClient: kitexClient,
	}

	return cli, nil
}

type clientImpl struct {
	service     string
	kitexClient authservice.Client
}

func (c *clientImpl) Service() string {
	return c.service
}

func (c *clientImpl) KitexClient() authservice.Client {
	return c.kitexClient
}

func (c *clientImpl) DeliverTokenByRPC(ctx context.Context, Req *auth.DeliverTokenReq, callOptions ...callopt.Option) (r *auth.DeliveryResp, err error) {
	return c.kitexClient.DeliverTokenByRPC(ctx, Req, callOptions...)
}

func (c *clientImpl) VerifyTokenByRPC(ctx context.Context, Req *auth.VerifyTokenReq, callOptions ...callopt.Option) (r *auth.VerifyResp, err error) {
	return c.kitexClient.VerifyTokenByRPC(ctx, Req, callOptions...)
}

func (c *clientImpl) RefreshTokenByRPC(ctx context.Context, Req *auth.RefreshTokenReq, callOptions ...callopt.Option) (r *auth.RefreshTokenResp, err error) {
	return c.kitexClient.RefreshTokenByRPC(ctx, Req, callOptions...)
}

func (c *clientImpl) RevokeTokenByRPC(ctx context.Context, Req *auth.RevokeTokenReq, callOptions ...callopt.Option) (r *auth.RevokeResp, err error) {
	return c.kitexClient.RevokeTokenByRPC(ctx, Req, callOptions...)
}

func (c *clientImpl) AddPermission(ctx context.Context, Req *auth.AddPermissionReq, callOptions ...callopt.Option) (r *auth.Empty, err error) {
	return c.kitexClient.AddPermission(ctx, Req, callOptions...)
}

func (c *clientImpl) CheckIfUserBanned(ctx context.Context, Req *auth.CheckIfUserBannedReq, callOptions ...callopt.Option) (r *auth.CheckIfUserBannedResp, err error) {
	return c.kitexClient.CheckIfUserBanned(ctx, Req, callOptions...)
}
