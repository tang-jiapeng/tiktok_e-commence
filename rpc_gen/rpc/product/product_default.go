package product

import (
	"context"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
	product "tiktok_e-commerce/rpc_gen/kitex_gen/product"
)

func ListProducts(ctx context.Context, req *product.ListProductsReq, callOptions ...callopt.Option) (resp *product.ListProductsResp, err error) {
	resp, err = defaultClient.ListProducts(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "ListProducts call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetProduct(ctx context.Context, req *product.GetProductReq, callOptions ...callopt.Option) (resp *product.GetProductResp, err error) {
	resp, err = defaultClient.GetProduct(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetProduct call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func SearchProducts(ctx context.Context, req *product.SearchProductsReq, callOptions ...callopt.Option) (resp *product.SearchProductsResp, err error) {
	resp, err = defaultClient.SearchProducts(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SearchProducts call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func InsertProduct(ctx context.Context, req *product.InsertProductReq, callOptions ...callopt.Option) (resp *product.InsertProductResp, err error) {
	resp, err = defaultClient.InsertProduct(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "InsertProduct call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func SelectProduct(ctx context.Context, req *product.SelectProductReq, callOptions ...callopt.Option) (resp *product.SelectProductResp, err error) {
	resp, err = defaultClient.SelectProduct(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SelectProduct call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func SelectProductList(ctx context.Context, req *product.SelectProductListReq, callOptions ...callopt.Option) (resp *product.SelectProductListResp, err error) {
	resp, err = defaultClient.SelectProductList(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SelectProductList call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteProduct(ctx context.Context, req *product.DeleteProductReq, callOptions ...callopt.Option) (resp *product.DeleteProductResp, err error) {
	resp, err = defaultClient.DeleteProduct(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteProduct call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateProduct(ctx context.Context, req *product.UpdateProductReq, callOptions ...callopt.Option) (resp *product.UpdateProductResp, err error) {
	resp, err = defaultClient.UpdateProduct(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateProduct call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func LockProductQuantity(ctx context.Context, req *product.ProductLockQuantityRequest, callOptions ...callopt.Option) (resp *product.ProductLockQuantityResponse, err error) {
	resp, err = defaultClient.LockProductQuantity(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "LockProductQuantity call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func SelectCategory(ctx context.Context, req *product.CategorySelectReq, callOptions ...callopt.Option) (resp *product.CategorySelectResp, err error) {
	resp, err = defaultClient.SelectCategory(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SelectCategory call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func InsertCategory(ctx context.Context, req *product.CategoryInsertReq, callOptions ...callopt.Option) (resp *product.CategoryInsertResp, err error) {
	resp, err = defaultClient.InsertCategory(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "InsertCategory call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteCategory(ctx context.Context, req *product.CategoryDeleteReq, callOptions ...callopt.Option) (resp *product.CategoryDeleteResp, err error) {
	resp, err = defaultClient.DeleteCategory(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteCategory call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateCategory(ctx context.Context, req *product.CategoryUpdateReq, callOptions ...callopt.Option) (resp *product.CategoryUpdateResp, err error) {
	resp, err = defaultClient.UpdateCategory(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateCategory call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func SelectBrand(ctx context.Context, req *product.BrandSelectReq, callOptions ...callopt.Option) (resp *product.BrandSelectResp, err error) {
	resp, err = defaultClient.SelectBrand(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "SelectBrand call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func InsertBrand(ctx context.Context, req *product.BrandInsertReq, callOptions ...callopt.Option) (resp *product.BrandInsertResp, err error) {
	resp, err = defaultClient.InsertBrand(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "InsertBrand call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteBrand(ctx context.Context, req *product.BrandDeleteReq, callOptions ...callopt.Option) (resp *product.BrandDeleteResp, err error) {
	resp, err = defaultClient.DeleteBrand(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteBrand call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateBrand(ctx context.Context, req *product.BrandUpdateReq, callOptions ...callopt.Option) (resp *product.BrandUpdateResp, err error) {
	resp, err = defaultClient.UpdateBrand(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateBrand call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
