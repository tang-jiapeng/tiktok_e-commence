package task

import (
	"context"
	myRedis "tiktok_e-commerce/order/biz/dal/redis"
	redisKeys "tiktok_e-commerce/order/utils/redis"

	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/pkg/errors"
	"github.com/redis/go-redis/v9"
	"github.com/xxl-job/xxl-job-executor-go"
)

// CleanNodeIDTask 删除记录节点id的redis key，当value大于30时删除，防止节点id重复
func CleanNodeIDTask(ctx context.Context, param *xxl.RunReq) (msg string) {
	script := `
		local key = KEYS[1]
		local value = tonumber(redis.call('GET', key))
		if value and value > 30 then
    		redis.call('DEL', key)
		end
    `
	err := myRedis.RedisClient.Eval(ctx, script, []string{redisKeys.OrderServiceNodeIdKey}).Err()
	if err != nil || !errors.Is(err, redis.Nil) {
		klog.Errorf("清除节点id key失败，err: %v", err)
		return "清除节点id key失败：" + err.Error()
	}
	return "success"
}
