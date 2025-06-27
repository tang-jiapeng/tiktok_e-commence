package main

import (
	"context"
	"net"
	"tiktok_e-commerce/product/biz/task"
	"tiktok_e-commerce/product/infra/elastic"
	"tiktok_e-commerce/product/infra/kafka"
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
	"github.com/xxl-job/xxl-job-executor-go"
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
	elastic.InitClient()
	kafka.Init()
	opts := kitexInit()

	svr := productcatalogservice.NewServer(new(ProductCatalogServiceImpl), opts...)

	//将任务注册到xxl-job中
	go xxlJobInit()

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

func xxlJobInit() {
	xxljobAddr := conf.GetConf().XxlJob.XxlJobAddress
	exec := xxl.NewExecutor(
		xxl.ServerAddr(xxljobAddr+"/xxl-job-admin"),
		xxl.AccessToken(conf.GetConf().XxlJob.AccessToken), //请求令牌(默认为空)
		xxl.ExecutorIp(conf.GetConf().XxlJob.ExecutorIp),
		xxl.ExecutorPort("7777"),
		xxl.RegistryKey("tiktok-e-commerce-product-service"), //执行器名称
		xxl.SetLogger(&logger{}),
	)
	exec.Init()
	exec.Use(customMiddleware)
	//设置日志查看handler
	exec.LogHandler(customLogHandle)
	//注册任务handler
	exec.RegTask("task.RefreshElastic", task.RefreshElastic)

	klog.Fatal(exec.Run())
}

// 自定义日志处理器
func customLogHandle(req *xxl.LogReq) *xxl.LogRes {
	return &xxl.LogRes{
		Code: xxl.SuccessCode,
		Msg:  "",
		Content: xxl.LogResContent{
			FromLineNum: req.FromLineNum,
			ToLineNum:   2,
			LogContent:  "自定义日志handler",
			IsEnd:       true,
		},
	}
}

// xxl.Logger接口实现
type logger struct{}

func (l *logger) Info(format string, a ...interface{}) {
	klog.CtxInfof(context.Background(), format, a...)
}

func (l *logger) Error(format string, a ...interface{}) {
	klog.CtxErrorf(context.Background(), format, a...)
}
func (l *logger) Debug(format string, a ...interface{}) {
	klog.CtxDebugf(context.Background(), format, a...)
}
func (l *logger) Warn(format string, a ...interface{}) {
	klog.CtxWarnf(context.Background(), format, a...)
}

func customMiddleware(tf xxl.TaskFunc) xxl.TaskFunc {
	return func(cxt context.Context, param *xxl.RunReq) string {
		klog.CtxInfof(context.Background(), "xxl-job 定时任务启动")
		res := tf(cxt, param)
		klog.CtxInfof(context.Background(), "xxl-job 定时任务结束")
		return res
	}
}
