package service

import (
	"context"
	"testing"
	auth "tiktok_e-commerce/rpc_gen/kitex_gen/auth"
)

func TestCheckIfUserBanned_Run(t *testing.T) {
	ctx := context.Background()
	s := NewCheckIfUserBannedService(ctx)
	// init req and assert value

	req := &auth.CheckIfUserBannedReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
