// Code generated by Kitex v0.9.1. DO NOT EDIT.
package userservice

import (
	server "github.com/cloudwego/kitex/server"
	user "tiktok_e-commerce/rpc_gen/kitex_gen/user"
)

// NewServer creates a server.Server with the given handler and options.
func NewServer(handler user.UserService, opts ...server.Option) server.Server {
	var options []server.Option

	options = append(options, opts...)

	svr := server.NewServer(options...)
	if err := svr.RegisterService(serviceInfo(), handler); err != nil {
		panic(err)
	}
	return svr
}

func RegisterService(svr server.Server, handler user.UserService, opts ...server.RegisterOption) error {
	return svr.RegisterService(serviceInfo(), handler, opts...)
}
