package service

import (
	"context"
	"testing"
	user "tiktok_e-commerce/rpc_gen/kitex_gen/user"
)

func TestGetReceiveAddress_Run(t *testing.T) {
	ctx := context.Background()
	s := NewGetReceiveAddressService(ctx)
	// init req and assert value

	req := &user.GetReceiveAddressReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
