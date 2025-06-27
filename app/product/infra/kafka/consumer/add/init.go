package add

import (
	"context"
	"strings"
	"tiktok_e-commerce/product/conf"
	"tiktok_e-commerce/product/infra/kafka/constant"

	"github.com/IBM/sarama"
	"github.com/cloudwego/kitex/pkg/klog"
)

func AddConsumer() {
	config := sarama.NewConfig()
	strategies := make([]sarama.BalanceStrategy, 1)
	strategies[0] = sarama.NewBalanceStrategyRange()
	config.Consumer.Group.Rebalance.GroupStrategies = strategies
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// 创建消费者
	brokers := strings.Split(conf.GetConf().Kafka.BizKafka.BootstrapServers, ",")
	groupId := "product_group_add"
	consumer, err := sarama.NewConsumerGroup(brokers, groupId, config)
	handler := AddProductHandler{}
	for {
		err = consumer.Consume(
			context.Background(),
			[]string{constant.AddTopic},
			handler,
		)
		if err != nil {
			klog.Error("Error from consumer: ", err)
		}
	}
}
