package kafka

import (
	"context"
	"strings"
	"sync"
	"tiktok_e-commerce/product/conf"

	"github.com/IBM/sarama"
	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	Producer sarama.AsyncProducer
	Consumer sarama.Consumer
	Topic    string
	once     sync.Once
	err      error
)

func InitClient() {
	// 配置Topic
	Topic = conf.GetConf().Kafka.BizKafka.ProductTopicId
	once.Do(func() {
		// 配置生产者
		err = InitProducer()
		if err != nil {
			panic(err)
		}
		// 配置消费者
		err = InitConsumer()
		if err != nil {
			panic(err)
		}
	})
}

func InitConsumer() (err error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = sarama.V1_1_0_0

	brokers := strings.Split(conf.GetConf().Kafka.BizKafka.BootstrapServers, ",")
	klog.Infof("Consumer brokers: %v", brokers)

	// 创建消费者
	Consumer, err = sarama.NewConsumer(brokers, config)
	if err != nil {
		klog.Errorf("Failed to start consumer: %v", err)
		return err
	}
	klog.Info("Successfully connected to Kafka consumer")

	// 创建消费者组
	group, err := sarama.NewConsumerGroup(brokers, "product_group", config)
	if err != nil {
		klog.Errorf("Failed to start consumer group: %v", err)
		return err
	}

	// 在单独 goroutine 中运行消费者组
	go func() {
		for {
			err := group.Consume(context.Background(), []string{Topic}, ProductKafkaConsumer{})
			if err != nil {
				klog.Errorf("Consumer group error: %v", err)
			}
			// 可添加重试逻辑或退出条件
		}
	}()

	return
}

func InitProducer() (err error) {
	config := sarama.NewConfig()
	// 配置生产者
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = 0
	config.Version = sarama.V1_1_0_0
	config.Producer.Compression = sarama.CompressionGZIP
	// 创建生产者
	brokers := strings.Split(conf.GetConf().Kafka.BizKafka.BootstrapServers, ",")
	Producer, err = sarama.NewAsyncProducer(brokers, config)
	if err != nil {
		klog.Error("Failed to start producer:", err)
	}
	klog.Info("Successfully connected to kafka", Producer)
	return
}
