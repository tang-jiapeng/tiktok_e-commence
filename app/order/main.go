package main

import (
	"net"
	"tiktok_e-commerce/common/infra/nacos"
	"tiktok_e-commerce/order/biz/task"
	"tiktok_e-commerce/order/infra/kafka"
	"tiktok_e-commerce/order/infra/rpc"
	"tiktok_e-commerce/order/utils"
	"time"

	"tiktok_e-commerce/common/mtl"
	"tiktok_e-commerce/order/conf"
	"tiktok_e-commerce/product/biz/dal"
	"tiktok_e-commerce/rpc_gen/kitex_gen/order/orderservice"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"github.com/xxl-job/xxl-job-executor-go"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {

	mtl.InitMetric(conf.GetConf().Kitex.Service, conf.GetConf().Kitex.MetricsPort)

	dal.Init()
	rpc.InitClient()
	kafka.Init()
	utils.InitSnowflake()

	xxljobInit()

	opts := kitexInit()

	svr := orderservice.NewServer(new(OrderServiceImpl), opts...)

	err := svr.Run()
	if err != nil {
		klog.Error(err.Error())
	}
}

func kitexInit() (opts []server.Option) {
	// address
	addr, err := net.ResolveTCPAddr("tcp", conf.GetConf().Kitex.Address)
	if err != nil {
		panic(err)
	}
	opts = append(opts, server.WithServiceAddr(addr))

	// service info
	opts = append(opts, server.WithServerBasicInfo(&rpcinfo.EndpointBasicInfo{
		ServiceName: conf.GetConf().Kitex.Service,
	}))

	r := nacos.RegisterService()
	opts = append(opts, server.WithRegistry(r))

	// klog
	logger := kitexlogrus.NewLogger()
	klog.SetLogger(logger)
	klog.SetLevel(conf.LogLevel())
	asyncWriter := &zapcore.BufferedWriteSyncer{
		WS: zapcore.AddSync(&lumberjack.Logger{
			Filename:   conf.GetConf().Kitex.LogFileName,
			MaxSize:    conf.GetConf().Kitex.LogMaxSize,
			MaxBackups: conf.GetConf().Kitex.LogMaxBackups,
			MaxAge:     conf.GetConf().Kitex.LogMaxAge,
		}),
		FlushInterval: time.Minute,
	}
	klog.SetOutput(asyncWriter)
	server.RegisterShutdownHook(func() {
		asyncWriter.Sync()
	})
	return
}

func xxljobInit() {
	exec := xxl.NewExecutor(
		xxl.ServerAddr(conf.GetConf().XxlJob.XxlJobAddress+"/xxl-job-admin"),
		xxl.AccessToken(conf.GetConf().XxlJob.AccessToken),
		xxl.ExecutorIp(conf.GetConf().XxlJob.ExecutorIp),
		xxl.ExecutorPort("7777"),
		xxl.RegistryKey("tiktok-e-commerce-order-service"),
	)
	exec.Init()
	server.RegisterShutdownHook(func() {
		exec.Stop()
	})

	exec.RegTask("CleanNodeIDTask", task.CleanNodeIDTask)
	exec.RegTask("CancelOrderTask", task.CancelOrderTask)

	err := exec.Run()
	if err != nil {
		return
	}
}
