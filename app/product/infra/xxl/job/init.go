package job

import (
	"tiktok_e-commerce/product/conf"
	"tiktok_e-commerce/product/infra/xxl/task"

	"github.com/cloudwego/kitex/server"
	"github.com/xxl-job/xxl-job-executor-go"
)

func XxlJobInit() {
	xxljobAddr := conf.GetConf().XxlJob.XxlJobAddress
	exec := xxl.NewExecutor(
		xxl.ServerAddr(xxljobAddr+"/xxl-admin"),
		xxl.AccessToken(conf.GetConf().XxlJob.AccessToken), //请求令牌(默认为空)
		xxl.ExecutorIp(conf.GetConf().XxlJob.ExecutorIp),
		xxl.ExecutorPort("7777"),
		xxl.RegistryKey("tiktok-e-commerce-product-service"), //执行器名称
	)
	exec.Init()

	server.RegisterShutdownHook(func() {
		exec.Stop()
	})

	//注册任务handler
	exec.RegTask("RefreshElastic", task.RefreshElastic)

	err := exec.Run()
	if err != nil {
		return
	}
}
