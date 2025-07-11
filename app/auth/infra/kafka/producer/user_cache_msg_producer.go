package producer

import (
	"strconv"
	"tiktok_e-commerce/auth/conf"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
)

var (
	producer sarama.AsyncProducer
	err      error
)

func InitUserCacheMessageProducer() {
	config := sarama.NewConfig()
	// 只等待leader确认，接受不完全保证消息不丢失
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewHashPartitioner

	producer, err = sarama.NewAsyncProducer(conf.GetConf().Kafka.BizKafka.BootstrapServers, config)
	if err != nil {
		panic(err.Error())
	}
	go func() {
		for msg := range producer.Successes() {
			klog.Infof("消息发送成功 topic:%s partition:%d offset:%d\n", msg.Topic, msg.Partition, msg.Offset)
		}
	}()
	go func() {
		for err = range producer.Errors() {
			klog.Errorf("消息发送失败: %v\n", err)
		}
	}()

	server.RegisterShutdownHook(func() {
		_ = producer.Close()
	})
}

func sendMessage(topic string, message []byte, key string) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
		Key:   sarama.StringEncoder(key),
	}
	producer.Input() <- msg
}

type UserCacheMessage struct {
	UserId int32 `json:"user_id"`
}

func SendUserCacheMessage(userId int32) {
	msg := UserCacheMessage{
		UserId: userId,
	}
	msgStr, _ := sonic.Marshal(msg)
	sendMessage("auth_service_deliver_token", msgStr, strconv.FormatInt(int64(userId), 10))
}
