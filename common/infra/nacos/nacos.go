package nacos

import (
	"os"

	kitexRigistry "github.com/cloudwego/kitex/pkg/registry"
	"github.com/kitex-contrib/registry-nacos/registry"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
)

func RegisterService() kitexRigistry.Registry {
	namingClient := GetNamingClient()
	r := registry.NewNacosRegistry(namingClient)
	return r
}

func GetNamingClient() naming_client.INamingClient {
	clientConfig, serverConfigs := GetNacosConfig()
	namingClient, err := clients.CreateNamingClient(map[string]interface{}{
		"serverConfigs": serverConfigs,
		"clientConfig":  clientConfig,
	})
	if err != nil {
		panic(err)
	}
	return namingClient
}

func GetNacosConfig() (constant.ClientConfig, []constant.ServerConfig) {
	env := os.Getenv("env")
	var logLevel string
	if env == "dev" {
		logLevel = "debug"
	} else {
		logLevel = "info"
	}
	clientConfig := constant.ClientConfig{
		NamespaceId: "e45ccc29-3e7d-4275-917b-febc49052d58",
		TimeoutMs:   5000,
		LogLevel:    logLevel,
	}
	//nacos_ip := os.Getenv("NACOS_ADDR")
	nacos_ip := "127.0.0.1"
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr: nacos_ip,
			Port:   8848,
		},
	}
	return clientConfig, serverConfigs
}
