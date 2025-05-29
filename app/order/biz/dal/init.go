package dal

import (
	"tiktok_e-commerce/order/biz/dal/mysql"
	"tiktok_e-commerce/order/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
