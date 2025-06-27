package del

import (
	"tiktok_e-commerce/product/infra/kafka/model"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
)

type DeleteProductHandler struct{}

func (DeleteProductHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (DeleteProductHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (DeleteProductHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) (err error) {
	ctx := session.Context()
	for msg := range claim.Messages() {
		klog.Infof("消费者接受到消息，Received message: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s \n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		session.MarkMessage(msg, "start")
		value := msg.Value
		dataVo := model.DeleteProductSendToKafka{}
		//将消息反序列化成Product结构体
		_ = sonic.Unmarshal(value, &dataVo)

		// TODO: 接入链路追踪

		// 发去redis
		err = DeleteProductToRedis(ctx, &dataVo)
		if err != nil {
			session.MarkMessage(msg, err.Error())
			return err
		}
		//发去es
		err = DeleteProductToElasticSearch(ctx, &dataVo)
		if err != nil {
			session.MarkMessage(msg, err.Error())
			return err
		}
		session.MarkMessage(msg, "done")
		session.Commit()
	}
	return nil
}
