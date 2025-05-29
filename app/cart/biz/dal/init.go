package dal

import (
	"tiktok_e-commerce/cart/biz/dal/mysql"
	"tiktok_e-commerce/cart/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
