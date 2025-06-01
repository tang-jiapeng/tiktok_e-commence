package redis

import (
	"context"
	"tiktok_e-commerce/auth/biz/dal/redis"
	"time"

	"github.com/cloudwego/kitex/pkg/klog"
)

func SetVal(ctx context.Context, key string, val interface{}, expiration time.Duration) (string, error) {
	result, err := redis.RedisClient.Set(ctx, key, val, expiration).Result()
	if err != nil {
		klog.Error("Redis Set操作失败: %v", err)
	}
	return result, err
}

func GetVal(ctx context.Context, token string) (string, error) {
	result, err := redis.RedisClient.Get(ctx, token).Result()
	if err != nil {
		klog.Error("Redis Get操作失败: %v", err)
	}
	return result, err
}
