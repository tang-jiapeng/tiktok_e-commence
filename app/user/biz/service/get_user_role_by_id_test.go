package service

import (
	"context"
	"testing"
	user "tiktok_e-commerce/rpc_gen/kitex_gen/user"
)

func TestGetUserRoleById_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetUserRoleByIdService(ctx)
	// init req and assert value

	req := &user.GetUserRoleByIdReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
