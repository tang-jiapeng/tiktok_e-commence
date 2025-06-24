package task

import (
	"context"
	"strconv"
	"strings"
	"tiktok_e-commerce/product/biz/dal/redis"
	"tiktok_e-commerce/product/biz/model"
	"tiktok_e-commerce/product/biz/vo"
	"tiktok_e-commerce/product/conf"
	"tiktok_e-commerce/product/infra/elastic"
	"tiktok_e-commerce/product/infra/kafka"
	"time"

	"github.com/IBM/sarama"
	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
	"github.com/elastic/go-elasticsearch/v7/esapi"
)

func AddProduct(product *model.Product) (err error) {
	sonicData, err := sonic.Marshal(vo.ProductKafkaDataVO{
		Type: vo.Type{
			Name: vo.Add,
		},
		Product: vo.ProductSendToKafka{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
			Price:       product.Price,
			Picture:     product.Picture,
		},
	})
	_, _, err = kafka.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: kafka.Topic,
		Value: sarama.ByteEncoder(sonicData),
	})
	if err != nil {
		return err
	}
	return
}

func DeleteProduct(product *model.Product) (err error) {
	sonicData, err := sonic.Marshal(vo.ProductKafkaDataVO{
		Type: vo.Type{
			Name: vo.Delete,
		},
		Product: vo.ProductSendToKafka{
			ID: product.ID,
		},
	})
	kafka.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: kafka.Topic,
		Value: sarama.ByteEncoder(sonicData),
	})
	return
}

func UpdateProduct(product *model.Product) (err error) {
	sonicData, err := sonic.Marshal(vo.ProductKafkaDataVO{
		Type: vo.Type{
			Name: vo.Update,
		},
		Product: vo.ProductSendToKafka{
			ID:          product.ID,
			Name:        product.Name,
			Description: product.Description,
		},
	})
	kafka.Producer.SendMessage(&sarama.ProducerMessage{
		Topic: kafka.Topic,
		Value: sarama.ByteEncoder(sonicData),
	})
	return
}

func DeleteProductToElasticSearch(id int64) {
	body := vo.ProductSearchQueryBody{
		Query: &vo.ProductSearchQuery{
			Term: &vo.ProductSearchTermQuery{
				"id": id,
			},
		},
	}
	sonicData, err := sonic.Marshal(body)
	if err != nil {
		return
	}
	request, _ := esapi.DeleteByQueryRequest{
		Index: []string{"product"},
		Body:  strings.NewReader(string(sonicData)),
	}.Do(context.Background(), elastic.ElasticClient)
	print(request.StatusCode)
	return
}

func DeletProductToRedis(id int64) {
	key := "product:" + strconv.FormatInt(id, 10)
	_, err := redis.RedisClient.Del(context.Background(), key).Result()
	if err != nil {
		klog.Error("redis delete error", err)
		return
	}
	return
}

func UpdateProductToElasticSearch(product *vo.ProductSendToKafka) {
	body := vo.ProductSearchQueryBody{
		Query: &vo.ProductSearchQuery{
			Term: &vo.ProductSearchTermQuery{
				"id": product.ID,
			},
		},
		Doc: &vo.ProductSearchDoc{
			Name:        product.Name,
			Description: product.Description,
		},
	}
	sonicData, err := sonic.Marshal(body)
	if err != nil {
		return
	}
	request, _ := esapi.UpdateByQueryRequest{
		Index: []string{"product"},
		Body:  strings.NewReader(string(sonicData)),
	}.Do(context.Background(), elastic.ElasticClient)
	print(request.StatusCode)
	return
}

func UpdateProductToRedis(product *vo.ProductSendToKafka) {
	pro := vo.ProductRedisDataVO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
		Price:       product.Price,
		Picture:     product.Picture,
	}
	key := "product:" + strconv.FormatInt(product.ID, 10)
	// 调用redis的set方法将数据导入到redis缓存中
	marshal, err := sonic.MarshalString(pro)
	if err != nil {
		klog.Error("序列化失败", err)
		return
	}
	_, err = redis.RedisClient.Set(context.Background(), key, marshal, 1*time.Hour).Result()
	if err != nil {
		klog.Error("redis set error", err)
		return
	}
	return
}

func AddProductToElasticSearch(product *vo.ProductSendToKafka) {
	pro := vo.ProductSearchDataVo{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
	}
	sonicData, _ := sonic.Marshal(pro)
	response, _ := esapi.IndexRequest{
		Index: "product",
		Body:  strings.NewReader(string(sonicData)),
	}.Do(context.Background(), elastic.ElasticClient)
	print(response.StatusCode)
	return
}

func AddProductToRedis(product *vo.ProductSendToKafka) {
	pro := vo.ProductRedisDataVO{
		ID:          product.ID,
		Name:        product.Name,
		Description: product.Description,
	}
	key := "product:" + strconv.FormatInt(product.ID, 10)
	marshal, err := sonic.MarshalString(pro)
	if err != nil {
		klog.Error("序列化失败", err)
		return
	}
	_, err = redis.RedisClient.Set(context.Background(), key, marshal, 1*time.Hour).Result()
	if err != nil {
		klog.Error("redis set error", err)
		return
	}
	return
}

type ProductKafkaHandler struct {
}

func (ProductKafkaHandler) Setup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (ProductKafkaHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (h ProductKafkaHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	count := 0
	batchSize := 10
	for msg := range claim.Messages() {
		klog.CtxInfof(context.Background(), "消费者接受到消息，Received message: Topic=%s, Partition=%d, Offset=%d, Key=%s, Value=%s \n",
			msg.Topic, msg.Partition, msg.Offset, string(msg.Key), string(msg.Value))
		session.MarkMessage(msg, "start")
		value := msg.Value
		dataVo := vo.ProductKafkaDataVO{}
		_ = sonic.Unmarshal(value, &dataVo)
		//将消息反序列化成Product结构体
		switch dataVo.Type.Name {
		case vo.Add:
			AddProductToElasticSearch(&dataVo.Product)
			AddProductToRedis(&dataVo.Product)
		case vo.Update:
			UpdateProductToElasticSearch(&dataVo.Product)
			UpdateProductToRedis(&dataVo.Product)
		case vo.Delete:
			DeleteProductToElasticSearch(dataVo.Product.ID)
			DeletProductToRedis(dataVo.Product.ID)
		}
		count++
		session.MarkMessage(msg, "done")
		if count >= batchSize {
			count = 0
			session.Commit()
		}
	}
	session.Commit()
	return nil
}

func Consumer() {
	config := sarama.NewConfig()
	startegies := make([]sarama.BalanceStrategy, 1)
	startegies[0] = sarama.NewBalanceStrategyRange()
	config.Consumer.Group.Rebalance.GroupStrategies = startegies
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	// 创建消费者
	brokers := strings.Split(conf.GetConf().Kafka.BizKafka.BootstrapServers, ",")
	consumer, err := sarama.NewConsumerGroup(brokers, "product_group", config)
	handler := ProductKafkaHandler{}
	for {
		err = consumer.Consume(
			context.Background(),
			[]string{conf.GetConf().Kafka.BizKafka.ProductTopicId},
			handler,
		)
		if err != nil {
			panic(err)
		}
	}
}
