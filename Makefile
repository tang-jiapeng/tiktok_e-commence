ROOT_MOD = tiktok_e-commerce

# 生成所有 RPC 客户端代码
.PHONY: gen_rpc_clients
gen_rpc_clients:
	@cd rpc_gen && \
	cwgo client --type RPC --module ${ROOT_MOD}/rpc_gen -I ../idl --service auth --idl ../idl/auth.proto && \
	cwgo client --type RPC --module ${ROOT_MOD}/rpc_gen -I ../idl --service cart --idl ../idl/cart.proto && \
	cwgo client --type RPC --module ${ROOT_MOD}/rpc_gen -I ../idl --service checkout --idl ../idl/checkout.proto && \
	cwgo client --type RPC --module ${ROOT_MOD}/rpc_gen -I ../idl --service order --idl ../idl/order.proto && \
	cwgo client --type RPC --module ${ROOT_MOD}/rpc_gen -I ../idl --service payment --idl ../idl/payment.proto && \
	cwgo client --type RPC --module ${ROOT_MOD}/rpc_gen -I ../idl --service product --idl ../idl/product.proto && \
	cwgo client --type RPC --module ${ROOT_MOD}/rpc_gen -I ../idl --service user --idl ../idl/user.proto && \
	# cwgo client --type RPC --module ${ROOT_MOD}/rpc_gen -I ../idl --service doubao_ai --idl ../idl/doubao_ai.proto && \
	go mod tidy

# 生成所有 RPC 服务端代码
.PHONY: gen_rpc_servers
gen_rpc_servers:
	@cd app/auth && cwgo server --type RPC --service auth --module ${ROOT_MOD}/auth --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/auth.proto && go mod tidy
	@cd app/cart && cwgo server --type RPC --service cart --module ${ROOT_MOD}/cart --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/cart.proto && go mod tidy
	@cd app/checkout && cwgo server --type RPC --service checkout --module ${ROOT_MOD}/checkout --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/checkout.proto && go mod tidy
	@cd app/order && cwgo server --type RPC --service order --module ${ROOT_MOD}/order --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/order.proto && go mod tidy
	@cd app/payment && cwgo server --type RPC --service payment --module ${ROOT_MOD}/payment --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/payment.proto && go mod tidy
	@cd app/product && cwgo server --type RPC --service product --module ${ROOT_MOD}/product --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/product.proto && go mod tidy
	@cd app/user && cwgo server --type RPC --service user --module ${ROOT_MOD}/user --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/user.proto && go mod tidy
	# @cd app/doubao_ai && cwgo server --type RPC --service doubao_ai --module ${ROOT_MOD}/doubao_ai --pass "-use ${ROOT_MOD}/rpc_gen/kitex_gen" -I ../../idl --idl ../../idl/doubao_ai.proto && go mod tidy

# 生成所有 HTTP API 代码
.PHONY: gen_api
gen_api:
	@cd app/api && \
	cwgo server --type HTTP --idl ../../idl/api/user_api.proto --server_name api --module ${ROOT_MOD}/api && \
	# cwgo server --type HTTP --idl ../../idl/api/cart_api.proto --server_name api --module ${ROOT_MOD}/api && \
	# cwgo server --type HTTP --idl ../../idl/api/order_api.proto --server_name api --module ${ROOT_MOD}/api && \
	# cwgo server --type HTTP --idl ../../idl/api/product_api.proto --server_name api --module ${ROOT_MOD}/api && \
	# cwgo server --type HTTP --idl ../../idl/api/checkout_api.proto --server_name api --module ${ROOT_MOD}/api && \
	cwgo server --type HTTP --idl ../../idl/api/payment_api.proto --server_name api --module ${ROOT_MOD}/api && \
	go mod tidy

# 生成所有代码（RPC 客户端、服务端和 HTTP API）
.PHONY: gen_all
gen_all: gen_rpc_clients gen_rpc_servers gen_api