package rpc

import (
	"os"
	"sync"
	"tiktok_e-commerce/common/clientsuite"
	"tiktok_e-commerce/rpc_gen/kitex_gen/auth/authservice"
	"tiktok_e-commerce/rpc_gen/kitex_gen/cart/cartservice"
	"tiktok_e-commerce/rpc_gen/kitex_gen/payment/paymentservice"
	"tiktok_e-commerce/rpc_gen/kitex_gen/product/productcatalogservice"
	"tiktok_e-commerce/rpc_gen/kitex_gen/user/userservice"
	"time"

	"github.com/cloudwego/hertz/pkg/common/hlog"
	"github.com/cloudwego/kitex/client"
)

var (
	AuthClient    authservice.Client
	UserClient    userservice.Client
	PaymentClient paymentservice.Client
	ProductClient productcatalogservice.Client
	CartClient    cartservice.Client
	once          sync.Once
	err           error
	registryAddr  string
	commonSuite   client.Option
)

func InitClient() {
	once.Do(func() {
		registryAddr = os.Getenv("NACOS_ADDR")
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: "api",
		})
		initUserClient()
		initAuthClient()
		initProductClient()
		initPaymentClient()
	})
}

func initUserClient() {
	UserClient, err = userservice.NewClient("user-service", commonSuite, client.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Fatal("init user client failed: ", err)
	}
}

func initAuthClient() {
	AuthClient, err = authservice.NewClient("auth-service", commonSuite, client.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Fatal("init auth client failed: ", err)
	}
}

func initProductClient() {
	ProductClient, err = productcatalogservice.NewClient("product-service", commonSuite, client.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Fatal("init product client failed: ", err)
	}
}

func initPaymentClient() {
	PaymentClient, err = paymentservice.NewClient("payment-service", commonSuite, client.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Fatal("init payment client failed: ", err)
	}
}

func initCartClient() {
	CartClient, err = cartservice.NewClient("cart-service", commonSuite, client.WithRPCTimeout(3*time.Second))
	if err != nil {
		hlog.Fatal("init cart client failed: ", err)
	}
}
