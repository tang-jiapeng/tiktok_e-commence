package kafka

import "tiktok_e-commerce/auth/infra/kafka/producer"

func Init() {
	producer.InitUserCacheMessageProducer()
}
