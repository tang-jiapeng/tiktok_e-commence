package service

import (
	"context"
	"testing"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"
)

func TestRevokeTokenByRPC_Run(t *testing.T) {
	ctx := context.Background()
	s := NewRevokeTokenByRPCService(ctx)
	// init req and assert value

	req := &auth.RevokeTokenReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
