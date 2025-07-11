package rpc

import (
	"os"
	"sync"
	"tiktok_e-commerce/common/clientsuite"
	"tiktok_e-commerce/product/conf"
	"tiktok_e-commerce/rpc_gen/kitex_gen/user/userservice"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	UserClient   userservice.Client
	once         sync.Once
	err          error
	registryAddr string
	commonSuite  client.Option
)

func InitClient() {
	once.Do(func() {
		registryAddr = os.Getenv("NACOS_ADDR")
		commonSuite = client.WithSuite(clientsuite.CommonGrpcClientSuite{
			RegistryAddr:       registryAddr,
			CurrentServiceName: conf.GetConf().Kitex.Service,
		})
		initUserClient()
	})
}

func initUserClient() {
	UserClient, err = userservice.NewClient("user-service", commonSuite)
	if err != nil {
		klog.Fatal("init user client failed: ", err)
	}
}
