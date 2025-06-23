package kafka

import "tiktok_e-commerce/user/infra/kafka/consumer"

func Init() {
	consumer.InitUserCacheMessageConsumer()
}
