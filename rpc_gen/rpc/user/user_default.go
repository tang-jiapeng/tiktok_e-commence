package user

import (
	"context"
	"github.com/cloudwego/kitex/client/callopt"
	"github.com/cloudwego/kitex/pkg/klog"
	user "tiktok_e-commerce/rpc_gen/kitex_gen/user"
)

func Register(ctx context.Context, req *user.RegisterReq, callOptions ...callopt.Option) (resp *user.RegisterResp, err error) {
	resp, err = defaultClient.Register(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Register call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func Login(ctx context.Context, req *user.LoginReq, callOptions ...callopt.Option) (resp *user.LoginResp, err error) {
	resp, err = defaultClient.Login(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "Login call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetUser(ctx context.Context, req *user.GetUserReq, callOptions ...callopt.Option) (resp *user.GetUserResp, err error) {
	resp, err = defaultClient.GetUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func UpdateUser(ctx context.Context, req *user.UpdateUserReq, callOptions ...callopt.Option) (resp *user.UpdateUserResp, err error) {
	resp, err = defaultClient.UpdateUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "UpdateUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func DeleteUser(ctx context.Context, req *user.DeleteUserReq, callOptions ...callopt.Option) (resp *user.DeleteUserResp, err error) {
	resp, err = defaultClient.DeleteUser(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "DeleteUser call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetUserRoleById(ctx context.Context, req *user.GetUserRoleByIdReq, callOptions ...callopt.Option) (resp *user.GetUserRoleByIdResp, err error) {
	resp, err = defaultClient.GetUserRoleById(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetUserRoleById call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func AddReceiveAddress(ctx context.Context, req *user.AddReceiveAddressReq, callOptions ...callopt.Option) (resp *user.AddReceiveAddressResp, err error) {
	resp, err = defaultClient.AddReceiveAddress(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "AddReceiveAddress call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}

func GetReceiveAddress(ctx context.Context, req *user.GetReceiveAddressReq, callOptions ...callopt.Option) (resp *user.GetReceiveAddressResp, err error) {
	resp, err = defaultClient.GetReceiveAddress(ctx, req, callOptions...)
	if err != nil {
		klog.CtxErrorf(ctx, "GetReceiveAddress call failed,err =%+v", err)
		return nil, err
	}
	return resp, nil
}
