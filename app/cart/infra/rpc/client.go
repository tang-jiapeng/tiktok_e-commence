package rpc

import (
	"os"
	"sync"
	"tiktok_e-commerce/cart/conf"
	"tiktok_e-commerce/common/clientsuite"
	"tiktok_e-commerce/rpc_gen/kitex_gen/product/productcatalogservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	ProductClient productcatalogservice.Client
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
			CurrentServiceName: conf.GetConf().Kitex.Service,
		})
		initProductClient()
	})
}

func initProductClient() {
	ProductClient, err = productcatalogservice.NewClient("product-service", commonSuite)
	if err != nil {
		klog.Fatal("init product client failed: ", err)
	}
}
