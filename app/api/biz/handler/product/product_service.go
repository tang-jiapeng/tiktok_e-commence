package product

import (
	"context"

	"github.com/cloudwego/hertz/pkg/app"
	"github.com/cloudwego/hertz/pkg/protocol/consts"
	"tiktok_e-commerce/api/biz/service"
	"tiktok_e-commerce/api/biz/utils"
	product "tiktok_e-commerce/api/hertz_gen/api/product"
)

// Search .
// @router /product/search [POST]
func Search(ctx context.Context, c *app.RequestContext) {
	var err error
	var req product.ProductRequest
	err = c.BindAndValidate(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	resp := &product.ProductResponse{}
	resp, err = service.NewSearchService(ctx, c).Run(&req)
	if err != nil {
		utils.SendErrResponse(ctx, c, consts.StatusOK, err)
		return
	}

	utils.SendSuccessResponse(ctx, c, consts.StatusOK, resp)
}
