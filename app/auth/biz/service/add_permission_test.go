package service

import (
	"context"
	"testing"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"
)

func TestAddPermission_Run(t *testing.T) {
	ctx := context.Background()
	s := NewAddPermissionService(ctx)
	// init req and assert value

	req := &auth.AddPermissionReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
