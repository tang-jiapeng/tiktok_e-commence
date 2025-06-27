package update

import (
	"context"
	"strconv"
	"tiktok_e-commerce/product/biz/dal/redis"
	"tiktok_e-commerce/product/biz/vo"
	"tiktok_e-commerce/product/infra/kafka/model"
	"time"

	"github.com/bytedance/sonic"
	"github.com/cloudwego/kitex/pkg/klog"
)

func UpdateProductToRedis(ctx context.Context, product *model.UpdateProductSendToKafka) (err error) {
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
		return err
	}
	_, err = redis.RedisClient.Set(ctx, key, marshal, 1*time.Hour).Result()
	if err != nil {
		klog.Error("redis set error ", err)
		return err
	}
	return nil
}
