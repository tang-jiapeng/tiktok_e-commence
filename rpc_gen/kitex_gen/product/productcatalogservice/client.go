// Code generated by Kitex v0.9.1. DO NOT EDIT.

package productcatalogservice

import (
	"context"
	client "github.com/cloudwego/kitex/client"
	callopt "github.com/cloudwego/kitex/client/callopt"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"
)

// Client is designed to provide IDL-compatible methods with call-option parameter for kitex framework.
type Client interface {
	ListProducts(ctx context.Context, Req *product.ListProductsReq, callOptions ...callopt.Option) (r *product.ListProductsResp, err error)
	GetProduct(ctx context.Context, Req *product.GetProductReq, callOptions ...callopt.Option) (r *product.GetProductResp, err error)
	SearchProducts(ctx context.Context, Req *product.SearchProductsReq, callOptions ...callopt.Option) (r *product.SearchProductsResp, err error)
	InsertProduct(ctx context.Context, Req *product.InsertProductReq, callOptions ...callopt.Option) (r *product.InsertProductResp, err error)
	SelectProduct(ctx context.Context, Req *product.SelectProductReq, callOptions ...callopt.Option) (r *product.SelectProductResp, err error)
	SelectProductList(ctx context.Context, Req *product.SelectProductListReq, callOptions ...callopt.Option) (r *product.SelectProductListResp, err error)
	DeleteProduct(ctx context.Context, Req *product.DeleteProductReq, callOptions ...callopt.Option) (r *product.DeleteProductResp, err error)
	UpdateProduct(ctx context.Context, Req *product.UpdateProductReq, callOptions ...callopt.Option) (r *product.UpdateProductResp, err error)
	LockProductQuantity(ctx context.Context, Req *product.ProductLockQuantityRequest, callOptions ...callopt.Option) (r *product.ProductLockQuantityResponse, err error)
	SelectCategory(ctx context.Context, Req *product.CategorySelectReq, callOptions ...callopt.Option) (r *product.CategorySelectResp, err error)
	InsertCategory(ctx context.Context, Req *product.CategoryInsertReq, callOptions ...callopt.Option) (r *product.CategoryInsertResp, err error)
	DeleteCategory(ctx context.Context, Req *product.CategoryDeleteReq, callOptions ...callopt.Option) (r *product.CategoryDeleteResp, err error)
	UpdateCategory(ctx context.Context, Req *product.CategoryUpdateReq, callOptions ...callopt.Option) (r *product.CategoryUpdateResp, err error)
	SelectBrand(ctx context.Context, Req *product.BrandSelectReq, callOptions ...callopt.Option) (r *product.BrandSelectResp, err error)
	InsertBrand(ctx context.Context, Req *product.BrandInsertReq, callOptions ...callopt.Option) (r *product.BrandInsertResp, err error)
	DeleteBrand(ctx context.Context, Req *product.BrandDeleteReq, callOptions ...callopt.Option) (r *product.BrandDeleteResp, err error)
	UpdateBrand(ctx context.Context, Req *product.BrandUpdateReq, callOptions ...callopt.Option) (r *product.BrandUpdateResp, err error)
}

// NewClient creates a client for the service defined in IDL.
func NewClient(destService string, opts ...client.Option) (Client, error) {
	var options []client.Option
	options = append(options, client.WithDestService(destService))

	options = append(options, opts...)

	kc, err := client.NewClient(serviceInfo(), options...)
	if err != nil {
		return nil, err
	}
	return &kProductCatalogServiceClient{
		kClient: newServiceClient(kc),
	}, nil
}

// MustNewClient creates a client for the service defined in IDL. It panics if any error occurs.
func MustNewClient(destService string, opts ...client.Option) Client {
	kc, err := NewClient(destService, opts...)
	if err != nil {
		panic(err)
	}
	return kc
}

type kProductCatalogServiceClient struct {
	*kClient
}

func (p *kProductCatalogServiceClient) ListProducts(ctx context.Context, Req *product.ListProductsReq, callOptions ...callopt.Option) (r *product.ListProductsResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.ListProducts(ctx, Req)
}

func (p *kProductCatalogServiceClient) GetProduct(ctx context.Context, Req *product.GetProductReq, callOptions ...callopt.Option) (r *product.GetProductResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.GetProduct(ctx, Req)
}

func (p *kProductCatalogServiceClient) SearchProducts(ctx context.Context, Req *product.SearchProductsReq, callOptions ...callopt.Option) (r *product.SearchProductsResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SearchProducts(ctx, Req)
}

func (p *kProductCatalogServiceClient) InsertProduct(ctx context.Context, Req *product.InsertProductReq, callOptions ...callopt.Option) (r *product.InsertProductResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.InsertProduct(ctx, Req)
}

func (p *kProductCatalogServiceClient) SelectProduct(ctx context.Context, Req *product.SelectProductReq, callOptions ...callopt.Option) (r *product.SelectProductResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SelectProduct(ctx, Req)
}

func (p *kProductCatalogServiceClient) SelectProductList(ctx context.Context, Req *product.SelectProductListReq, callOptions ...callopt.Option) (r *product.SelectProductListResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SelectProductList(ctx, Req)
}

func (p *kProductCatalogServiceClient) DeleteProduct(ctx context.Context, Req *product.DeleteProductReq, callOptions ...callopt.Option) (r *product.DeleteProductResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DeleteProduct(ctx, Req)
}

func (p *kProductCatalogServiceClient) UpdateProduct(ctx context.Context, Req *product.UpdateProductReq, callOptions ...callopt.Option) (r *product.UpdateProductResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.UpdateProduct(ctx, Req)
}

func (p *kProductCatalogServiceClient) LockProductQuantity(ctx context.Context, Req *product.ProductLockQuantityRequest, callOptions ...callopt.Option) (r *product.ProductLockQuantityResponse, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.LockProductQuantity(ctx, Req)
}

func (p *kProductCatalogServiceClient) SelectCategory(ctx context.Context, Req *product.CategorySelectReq, callOptions ...callopt.Option) (r *product.CategorySelectResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SelectCategory(ctx, Req)
}

func (p *kProductCatalogServiceClient) InsertCategory(ctx context.Context, Req *product.CategoryInsertReq, callOptions ...callopt.Option) (r *product.CategoryInsertResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.InsertCategory(ctx, Req)
}

func (p *kProductCatalogServiceClient) DeleteCategory(ctx context.Context, Req *product.CategoryDeleteReq, callOptions ...callopt.Option) (r *product.CategoryDeleteResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DeleteCategory(ctx, Req)
}

func (p *kProductCatalogServiceClient) UpdateCategory(ctx context.Context, Req *product.CategoryUpdateReq, callOptions ...callopt.Option) (r *product.CategoryUpdateResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.UpdateCategory(ctx, Req)
}

func (p *kProductCatalogServiceClient) SelectBrand(ctx context.Context, Req *product.BrandSelectReq, callOptions ...callopt.Option) (r *product.BrandSelectResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.SelectBrand(ctx, Req)
}

func (p *kProductCatalogServiceClient) InsertBrand(ctx context.Context, Req *product.BrandInsertReq, callOptions ...callopt.Option) (r *product.BrandInsertResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.InsertBrand(ctx, Req)
}

func (p *kProductCatalogServiceClient) DeleteBrand(ctx context.Context, Req *product.BrandDeleteReq, callOptions ...callopt.Option) (r *product.BrandDeleteResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.DeleteBrand(ctx, Req)
}

func (p *kProductCatalogServiceClient) UpdateBrand(ctx context.Context, Req *product.BrandUpdateReq, callOptions ...callopt.Option) (r *product.BrandUpdateResp, err error) {
	ctx = client.NewCtxWithCallOptions(ctx, callOptions)
	return p.kClient.UpdateBrand(ctx, Req)
}
