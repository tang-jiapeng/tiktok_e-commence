package dal

import (
	"tiktok_e-commerce/product/biz/dal/mysql"
	"tiktok_e-commerce/product/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
