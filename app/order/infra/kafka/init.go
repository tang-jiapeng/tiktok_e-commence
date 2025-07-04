package kafka

import (
	"tiktok_e-commerce/order/infra/kafka/consumer"
	"tiktok_e-commerce/order/infra/kafka/producer"
)

func Init() {
	consumer.InitDelayOrderConsumer()
	producer.InitDelayOrderProducer()
}
