package consumer

import (
	"context"
	"fmt"
	"strings"
	"tiktok_e-commerce/order/biz/dal/mysql"
	"tiktok_e-commerce/order/biz/model"
	"tiktok_e-commerce/order/conf"
	"tiktok_e-commerce/order/infra/kafka/constant"
	model2 "tiktok_e-commerce/order/infra/kafka/model"
	"tiktok_e-commerce/order/infra/kafka/producer"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
)

type msgConsumerGroup struct{}

func (msgConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (msgConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (h msgConsumerGroup) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ctx := sess.Context()
	for msg := range claim.Messages() {
		topic := msg.Topic
		klog.Infof("收到消息，topic:%q partition:%d offset:%d  value:%s\n", topic, msg.Partition, msg.Offset, string(msg.Value))
		if !strings.HasPrefix(string(msg.Key), constant.DelayCancelOrderKeyPrefix) {
			sess.MarkMessage(msg, "")
			sess.Commit()
			continue
		}

		var delayCancelOrderMessage model2.DelayOrderMessage
		err := sonic.Unmarshal(msg.Value, &delayCancelOrderMessage)
		if err != nil {
			klog.Errorf("解析消息失败，topic:%q partition:%d offset:%d  value:%s\n", topic, msg.Partition, msg.Offset, string(msg.Value))
			sess.MarkMessage(msg, "")
			sess.Commit()
			return err
		}

		orderId := delayCancelOrderMessage.OrderID
		unpaid, err := checkOrderUnpaid(ctx, orderId)
		if err != nil {
			return err
		}
		if unpaid {
			// 继续发消息或取消订单
			switch delayCancelOrderMessage.Level {
			case constant.DelayTopic1mLevel:
				{
					producer.SendDelayOrder(orderId, constant.DelayTopic4mLevel)
				}
			case constant.DelayTopic4mLevel:
				{
					producer.SendDelayOrder(orderId, constant.DelayTopic5mLevel)
				}
			case constant.DelayTopic5mLevel:
				{
					// 订单超时，取消订单，有兜底方案，不管结果；取消订单是幂等的，不考虑重复消费
					cancelOrder(ctx, orderId)
				}
			}
		}
		sess.MarkMessage(msg, "")
		sess.Commit()
	}
	return nil
}

func cancelOrder(ctx context.Context, orderId string) {
	// 先取消支付，再取消订单
	err := cacnelCharge(ctx, orderId)
	if err != nil {
		klog.Errorf("取消支付失败，orderId:%s，err:%v", orderId, err)
		return
	}
	affectedRows, err := model.CancelOrder(ctx, mysql.DB, orderId)
	if err != nil || affectedRows == 0 {
		klog.Errorf("取消订单失败，orderId:%s，err:%v", orderId, err)
		return
	}
	klog.Info("订单取消成功，orderId:%s", orderId)
}

func cacnelCharge(ctx context.Context, orderId string) error {
	// TODO: 取消支付
	return nil
}

func checkOrderUnpaid(ctx context.Context, orderId string) (bool, error) {
	unpaid, err := model.CheckOrderUnpaid(ctx, mysql.DB, orderId)
	if err != nil {
		klog.Errorf("查询订单状态失败，orderId:%s，err:%v", orderId, err)
		return false, err
	}
	return unpaid, nil
}

func InitDelayOrderConsumer() {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Version = sarama.V3_5_0_0
	consumerConfig.Consumer.Offsets.AutoCommit.Enable = false
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumerConfig.Consumer.Offsets.Retry.Max = 3

	groupId := constant.DelayCancelOrderGroupId
	cGroup, err := sarama.NewConsumerGroup(conf.GetConf().Kafka.BizKafka.BootstrapServers, groupId, consumerConfig)
	if err != nil {
		panic(err)
	}

	go func() {
		for {
			err = cGroup.Consume(
				context.Background(),
				[]string{constant.DelayCancelOrderTopic},
				msgConsumerGroup{},
			)
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}()

	server.RegisterShutdownHook(func() {
		_ = cGroup.Close()
	})
}
