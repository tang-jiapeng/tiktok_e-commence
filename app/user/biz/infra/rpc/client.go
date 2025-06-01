package rpc

import (
	"os"
	"sync"
	"tiktok_e-commerce/common/clientsuite"
	"tiktok_e-commerce/rpc_gen/kitex_gen/auth/authservice"
	"tiktok_e-commerce/user/conf"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	AuthClient   authservice.Client
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
		initAuthClient()
	})
}

func initAuthClient() {
	AuthClient, err = authservice.NewClient("auth-service", commonSuite)
	if err != nil {
		klog.Fatal("init auth client failed: ", err)
	}
}
