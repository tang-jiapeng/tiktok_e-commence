package consumer

import (
	"tiktok_e-commerce/product/infra/kafka/consumer/add"
	del "tiktok_e-commerce/product/infra/kafka/consumer/delete"
	"tiktok_e-commerce/product/infra/kafka/consumer/update"
)

func InitConsumer() {
	go add.AddConsumer()
	go del.DeleteConsumer()
	go update.UpdateConsumer()
}
