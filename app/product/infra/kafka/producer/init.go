package producer

import (
	"strings"
	"tiktok_e-commerce/product/conf"

	"github.com/IBM/sarama"
	"github.com/cloudwego/kitex/pkg/klog"
)

var (
	Producer sarama.SyncProducer
	err      error
)

func InitProducerClient() {
	// 配置生产者
	err = InitProducer()
	if err != nil {
		return
	}
}

func InitProducer() (err error) {
	config := sarama.NewConfig()
	// 配置生产者
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = 0
	config.Version = sarama.V3_5_0_0
	config.Producer.Compression = sarama.CompressionGZIP
	// 创建生产者
	brokers := strings.Split(conf.GetConf().Kafka.BizKafka.BootstrapServers, ",")
	Producer, err = sarama.NewSyncProducer(brokers, config)
	if err != nil {
		klog.Error("Failed to start producer:", err)
	}
	klog.Info("Successfully connected to kafka", Producer)
	return
}
