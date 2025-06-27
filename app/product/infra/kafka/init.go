package kafka

import (
	"tiktok_e-commerce/product/infra/kafka/consumer"
	"tiktok_e-commerce/product/infra/kafka/producer"
)

func Init() {
	consumer.InitConsumer()
	producer.InitProducerClient()
}
