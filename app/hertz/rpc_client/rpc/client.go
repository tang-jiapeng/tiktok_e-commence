package rpc

import (
	"sync"
	"tiktok_e-commerce/app/hertz/conf"

	"github.com/cloudwego/kitex/client"
	"tiktok_e-commerce/common/clientsuite"
)


var (
	once          sync.Once
	commonSuite   client.Option
)

func InitClient() {
	once.Do(func() {
		commonSuite = client.WithSuite(
			clientsuite.CommonGrpcClientSuite{
				RegistryAddr:       conf.GetConf().Hertz.RegistryAddr,
				CurrentServiceName: conf.GetConf().Hertz.Service,
			},
		)
	})
}