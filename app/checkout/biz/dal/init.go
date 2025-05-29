package dal

import (
	"tiktok_e-commerce/checkout/biz/dal/mysql"
	"tiktok_e-commerce/checkout/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
