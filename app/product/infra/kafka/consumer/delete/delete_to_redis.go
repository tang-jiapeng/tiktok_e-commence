package del

import (
	"context"
	"strconv"
	"tiktok_e-commerce/product/biz/dal/redis"
	"tiktok_e-commerce/product/infra/kafka/model"

	"github.com/cloudwego/kitex/pkg/klog"
)

func DeleteProductToRedis(ctx context.Context, product *model.DeleteProductSendToKafka) (err error) {
	key := "product:" + strconv.FormatInt(product.ID, 10)
	// 调用redis的set方法将数据导入到redis缓存中
	_, err = redis.RedisClient.Del(ctx, key).Result()
	if err != nil {
		klog.Error("redis del error ", err)
		return err
	}
	return nil
}
