package utils

import (
	"context"

	"github.com/bwmarrin/snowflake"

	myRedis "tiktok_e-commerce/order/biz/dal/redis"
	redisKeys "tiktok_e-commerce/order/utils/redis"
)

var (
	node *snowflake.Node
	err  error
	ctx  = context.Background()
)

func InitSnowflake() {
	nodeId := allocateNodeID()
	node, err = snowflake.NewNode(nodeId)
	if err != nil {
		panic(err)
	}
}

func GetSnowFlakeID() string {
	return node.Generate().String()
}

func allocateNodeID() int64 {
	nodeId, err := myRedis.RedisClient.Incr(ctx, redisKeys.OrderServiceNodeIdKey).Result()
	if err != nil {
		panic(err)
	}
	return nodeId
}
