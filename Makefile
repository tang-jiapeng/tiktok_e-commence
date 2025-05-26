ROOT_MOD = tiktok_e-commerce

# RPC 模块列表，例如：RPC_MOD=product order
service = user

# 生成 RPC 客户端和服务端代码
.PHONY: gen-rpc
gen-rpc:
		@# generate client
		@cd rpc_gen && cwgo client --type RPC --service ${service} --module ${ROOT_MOD}/rpc_gen -I ../idl/rpc --idl ../idl/rpc/${service}.proto && go mod tidy
		@# generate server
		@cd app/${service} && cwgo server --type RPC --service ${service} --module ${ROOT_MOD}/app/${service} -I ../../idl/rpc --idl ../../idl/rpc/${service}.proto --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" && go mod tidy


# 生成 API 代码
.PHONY: gen-api
gen-api:
	@cd app/hertz && cwgo server --type HTTP --service hertz --module ${ROOT_MOD}/app/hertz --idl ../../idl/hertz/${service}.proto && go mod tidy
