package service

import (
	"context"
	"testing"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"
)

func TestUpdateCategory_Run(t *testing.T) {
	ctx := context.Background()
	s := NewUpdateCategoryService(ctx)
	// init req and assert value

	req := &product.CategoryUpdateReq{}
	resp, err := s.Run(req)
	t.Logf("err: %v", err)
	t.Logf("resp: %v", resp)

	// todo: edit your unit test

}
