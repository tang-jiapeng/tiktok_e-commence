package clientsuite

import (
	"tiktok_e-commerce/common/infra/nacos"

	"github.com/cloudwego/kitex/client"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/pkg/transmeta"
	"github.com/cloudwego/kitex/transport"
	"github.com/kitex-contrib/registry-nacos/resolver"
)

type CommonGrpcClientSuite struct {
	CurrentServiceName string
	RegistryAddr       string
}

func (s CommonGrpcClientSuite) Options() []client.Option {
	r := resolver.NewNacosResolver(nacos.GetNamingClient())
	opts := []client.Option{
		client.WithResolver(r),
		client.WithMetaHandler(transmeta.ClientHTTP2Handler),
		client.WithTransportProtocol(transport.GRPC),
	}

	opts = append(opts,
		client.WithClientBasicInfo(&rpcinfo.EndpointBasicInfo{
			ServiceName: s.CurrentServiceName,
		}),
	)

	return opts
}
