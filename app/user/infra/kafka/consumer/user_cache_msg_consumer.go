package consumer

import (
	"context"
	"fmt"
	"tiktok_e-commerce/user/biz/dal/mysql"
	"tiktok_e-commerce/user/biz/dal/redis"
	"tiktok_e-commerce/user/biz/model"
	"tiktok_e-commerce/user/conf"
	redisUtils "tiktok_e-commerce/user/utils/redis"
	"time"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/cloudwego/kitex/server"
)

type msgConsumerGroup struct {
}

func (m msgConsumerGroup) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (m msgConsumerGroup) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (m msgConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		klog.Infof("收到消息，topic:%q partition:%d offset:%d  value:%s\n", msg.Topic, msg.Partition, msg.Offset, string(msg.Value))
		userCacheMsg := UserCacheMessage{}
		err := sonic.Unmarshal(msg.Value, &userCacheMsg)
		if err != nil {
			klog.Errorf("反序列化消息失败，err：%v", err)
			return err
		}
		err = selectAndCacheUserInfo(session.Context(), userCacheMsg.UserId)
		if err != nil {
			klog.Errorf("缓存用户信息失败，err：%v", err)
			return err
		}
		err = selectAndCacheUserAddresses(session.Context(), userCacheMsg.UserId)
		if err != nil {
			klog.Errorf("缓存用户地址失败，err：%v", err)
			return err
		}
		session.MarkMessage(msg, "")
		session.Commit()
	}
	return nil
}

func selectAndCacheUserAddresses(ctx context.Context, userId int32) error {
	addresses, err := model.GetAdressList(mysql.DB, ctx, userId)
	if err != nil {
		return err
	}
	if len(addresses) == 0 {
		return nil
	}
	luaScript := `
		if redis.call('EXISTS', KEYS[1]) == 0 then
			return redis.call('RPUSH', KEYS[1], unpack(ARGV))
		else
			return 0
		end
	`
	key := redisUtils.GetUserAddressesKey(userId)
	addressStrs := make([]string, len(addresses))
	for i, address := range addresses {
		addressStr, err := sonic.Marshal(address)
		if err != nil {
			return err
		}
		addressStrs[i] = string(addressStr)
	}
	err = redis.RedisClient.Eval(ctx, luaScript, []string{key}, addressStrs).Err()
	if err != nil {
		return err
	}
	// 设置过期时间和access token的过期时间一致
	err = redis.RedisClient.Expire(ctx, key, time.Hour*2).Err()
	if err != nil {
		return err
	}
	return nil
}

func selectAndCacheUserInfo(ctx context.Context, userId int32) error {
	user, err := model.GetUserBasicInfoById(mysql.DB, ctx, userId)
	if err != nil {
		return err
	}
	key := redisUtils.GetUserKey(userId)
	err = redis.RedisClient.HSet(ctx, key, user.ToMap()).Err()
	if err != nil {
		return err
	}
	// 设置过期时间和access token的过期时间一致
	err = redis.RedisClient.Expire(ctx, key, time.Hour*2).Err()
	if err != nil {
		return err
	}
	return nil
}

type UserCacheMessage struct {
	UserId int32 `json:"user_id"`
}

func InitUserCacheMessageConsumer() {
	consumerConfig := sarama.NewConfig()
	consumerConfig.Version = sarama.V3_5_0_0
	consumerConfig.Consumer.Offsets.AutoCommit.Enable = false
	consumerConfig.Consumer.Offsets.Initial = sarama.OffsetNewest
	consumerConfig.Consumer.Offsets.Retry.Max = 3

	consumerGroup, err := sarama.NewConsumerGroup(
		conf.GetConf().Kafka.BizKafka.BootstrapServers, "cache-user-info-dev", consumerConfig)
	if err != nil {
		panic(err)
	}
	go func() {
		for {
			err = consumerGroup.Consume(
				context.Background(),
				[]string{"auth_service_deliver_token"},
				msgConsumerGroup{})
			if err != nil {
				fmt.Println(err.Error())
				break
			}
		}
	}()

	server.RegisterShutdownHook(func() {
		_ = consumerGroup.Close()
	})
}
