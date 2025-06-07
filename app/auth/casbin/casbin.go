package casbin

import (
	"context"
	"fmt"
	myredis "tiktok_e-commerce/auth/biz/dal/redis"
	"tiktok_e-commerce/auth/conf"

	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var Enforcer *casbin.Enforcer

func InitCasbin() {
	dsn := conf.GetConf().MySQL.DSN
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		klog.Errorf("连接数据库失败: %v", err)
		panic(err)
	}

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		klog.Errorf("创建casbin适配器失败: %v", err)
		panic(err)
	}

	Enforcer, err = casbin.NewEnforcer(conf.GetConf().Casbin.ModelPath, adapter)
	if err != nil {
		klog.Errorf("创建casbin执行器失败: %v", err)
		panic(err)
	}

	printAllPolicies()

	subscribeToRedisChannel(myredis.RedisClient, context.Background())
}

func subscribeToRedisChannel(redisClient *redis.Client, ctx context.Context) {
	pubsub := redisClient.Subscribe(ctx, "casbin_policy_updates")
	go func() {
		defer pubsub.Close()
		for {
			msg, err := pubsub.ReceiveMessage(ctx)
			if err != nil {
				klog.Errorf("接收 Redis 订阅消息失败: %v", err)
				continue
			}
			klog.Infof("接收到 Redis 订阅消息: %s", msg.Payload)

			// 重新加载Casbin权限
			err = Enforcer.LoadPolicy()
			if err != nil {
				klog.Errorf("重新加载Casbin权限失败: %v", err)
			} else {
				klog.Infof("重新加载Casbin权限成功")

			}
		}
	}()
}

func printAllPolicies() {
	policies, _ := Enforcer.GetFilteredPolicy(0) //  0 表示不过滤任何字段
	fmt.Println("当前所有权限策略：")
	for _, policy := range policies {
		fmt.Printf("%v\n", policy)
	}
}

func AddPolicy(sub, obj, act string) error {
	_, err := Enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		klog.Errorf("添加权限策略失败: %v", err)
		return errors.WithStack(err)
	} else {
		klog.Infof("添加权限策略成功")
	}
	return nil
}
