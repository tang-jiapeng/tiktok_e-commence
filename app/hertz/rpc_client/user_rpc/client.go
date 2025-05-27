package userrpcclient

import (
	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/klog"
	etcdRegistry "github.com/kitex-contrib/registry-etcd"
	config "tiktok_e-commerce/app/hertz/conf"
	"tiktok_e-commerce/rpc_gen/kitex_gen/user/userservice"
	userrpccl "tiktok_e-commerce/rpc_gen/rpc/user"
)

func GetUserRpcClient() (cli userservice.Client, err error) {
	conf := config.GetConf()
	r, err := etcdRegistry.NewEtcdResolver(conf.Registry.RegistryAddress)
	if err != nil {
		return nil, err
	}
	cli, err = userservice.NewClient(
		"user",
		client.WithResolver(r),
	)
	if err != nil {
		return nil, err
	}
	return cli, err
}

func InitUserRpcClient() {
	conf := config.GetConf()
	r, err := etcdRegistry.NewEtcdResolver(conf.Registry.RegistryAddress)
	if err != nil {
		klog.Error(err)
	}
	userrpccl.InitClient("user", client.WithResolver(r))
}
