package mtl

import (
	"net"
	"net/http"
	"tiktok_e-commerce/common/infra/nacos"

	"github.com/cloudwego/kitex/pkg/registry"
	"github.com/cloudwego/kitex/server"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var Register *prometheus.Registry

func InitMetric(serviceName, metricsPort string) {
	Register = prometheus.NewRegistry()
	Register.MustRegister(collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))

	r := nacos.RegisterService()

	addr, _ := net.ResolveTCPAddr("tcp", metricsPort)

	registryInfo := &registry.Info{
		ServiceName: "prometheus-" + serviceName,
		Addr:        addr,
		Weight:      1,
		Tags:        map[string]string{"service": serviceName},
	}

	_ = r.Register(registryInfo)

	server.RegisterShutdownHook(func() {
		r.Deregister(registryInfo)
	})

	http.Handle("/metrics", promhttp.HandlerFor(Register, promhttp.HandlerOpts{}))
	go http.ListenAndServe(metricsPort, nil)
}
