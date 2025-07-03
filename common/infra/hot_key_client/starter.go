package hotKeyClient

import (
	"tiktok_e-commerce/common/infra/hot_key_client/constants"
	"tiktok_e-commerce/common/infra/hot_key_client/listener"
	"tiktok_e-commerce/common/infra/hot_key_client/publisher"
	clientRedis "tiktok_e-commerce/common/infra/hot_key_client/redis"

	"github.com/redis/go-redis/v9"
)

func Start(redisClient *redis.Client, serviceName string) {
	clientRedis.Init(redisClient)
	constants.Init(serviceName)
	go publisher.PublishStarter()
	go listener.ListenStarter()
}
