package dal

import (
	"tiktok_e-commerce/payment/biz/dal/alipay"
	"tiktok_e-commerce/payment/biz/dal/mysql"
	"tiktok_e-commerce/payment/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
	alipay.Init()
}
