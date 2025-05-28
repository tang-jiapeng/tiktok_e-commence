ROOT_MOD = tiktok_e-commerce

# 生成 RPC 客户端和服务端代码
.PHONY: gen_rpc
gen_rpc:
	@if [ -z "$(SERVICE)" ]; then echo "Error: SERVICE variable is not set. Usage: make gen_rpc SERVICE=<service_name>"; exit 1; fi
	@# 生成客户端代码
	@cd rpc_gen && cwgo client --type RPC --service $(SERVICE) --module ${ROOT_MOD}/rpc_gen -I ../idl/rpc --idl ../idl/rpc/$(SERVICE).proto && go mod tidy
	@# 生成服务端代码
	@cd app/$(SERVICE) && cwgo server --type RPC --service $(SERVICE) --module ${ROOT_MOD}/app/$(SERVICE) -I ../../idl/rpc --idl ../../idl/rpc/$(SERVICE).proto --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" && go mod tidy

# 生成 API 代码
.PHONY: gen_api
gen_api:
	@if [ -z "$(SERVICE)" ]; then echo "Error: SERVICE variable is not set. Usage: make gen_api SERVICE=<service_name>"; exit 1; fi
	@cd app/hertz && cwgo server --type HTTP --service hertz --module ${ROOT_MOD}/app/hertz --idl ../../idl/hertz/$(SERVICE)_api.proto && go mod tidy