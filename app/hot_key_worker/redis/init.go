package redis

import (
	"context"
	"hot_key/constant"
	"hot_key/model/key"
	"tiktok_e-commerce/product/conf"

	"github.com/bytedance/sonic"
	"github.com/redis/go-redis/v9"
)

var Rdb *redis.Client

func Init() {
	Rdb = redis.NewClient(&redis.Options{
		Addr:     conf.GetConf().Redis.Address,
		Username: conf.GetConf().Redis.Username,
		Password: conf.GetConf().Redis.Password,
		DB:       conf.GetConf().Redis.DB,
	})
}

func PublishClientChannel(hotkeyModel key.HotKeyModel) (err error) {
	marshal, _ := sonic.Marshal(hotkeyModel)
	err = Rdb.Publish(context.Background(), constant.ClientChannel, marshal).Err()
	if err != nil {
		return err
	}
	return nil
}
