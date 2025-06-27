package producer

import (
	"context"
	"tiktok_e-commerce/product/infra/kafka/constant"
	"tiktok_e-commerce/product/infra/kafka/model"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
)

func UpdateToKafka(ctx context.Context, product model.UpdateProductSendToKafka) error {
	sonicData, err := sonic.Marshal(model.AddProductSendToKafka{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Picture:     product.Picture,
	})
	if err != nil {
		return err
	}
	msg := &sarama.ProducerMessage{
		Topic: constant.UpdateTopic,
		Value: sarama.ByteEncoder(sonicData),
	}

	_, _, err = Producer.SendMessage(msg)
	if err != nil {
		return err
	}
	return nil
}
