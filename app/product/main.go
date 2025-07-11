package main

import (
	"net"
	hotKeyClient "tiktok_e-commerce/common/infra/hot_key_client"
	"tiktok_e-commerce/common/infra/hot_key_client/constants"
	"tiktok_e-commerce/product/biz/dal/redis"
	"tiktok_e-commerce/product/infra/elastic"
	"tiktok_e-commerce/product/infra/kafka"
	"tiktok_e-commerce/product/infra/rpc"
	"tiktok_e-commerce/product/infra/xxl"
	"time"

	"tiktok_e-commerce/common/infra/nacos"
	"tiktok_e-commerce/common/mtl"
	"tiktok_e-commerce/product/biz/dal"
	"tiktok_e-commerce/product/conf"
	"tiktok_e-commerce/rpc_gen/kitex_gen/product/productcatalogservice"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/pkg/rpcinfo"
	"github.com/cloudwego/kitex/server"
	"github.com/joho/godotenv"
	kitexlogrus "github.com/kitex-contrib/obs-opentelemetry/logging/logrus"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		klog.Warn("Failed to load .env file: %v", err)
	}
	dal.Init()
	mtl.InitMetric(conf.GetConf().Kitex.Service, conf.GetConf().Kitex.MetricsPort)

	//启动hotKeyClient
	go hotKeyClient.Start(redis.RedisClient, constants.ProductService)

	elastic.InitClient()
	kafka.Init()
	rpc.InitClient()
	opts := kitexInit()

	svr := productcatalogservice.NewServer(new(ProductCatalogServiceImpl), opts...)

	//将任务注册到xxl-job中
	xxl.Init()
	err = svr.Run()
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
