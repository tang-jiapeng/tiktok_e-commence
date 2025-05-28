package dal

import (
	"tiktok_e-commerce/app/product/biz/dal/mysql"
	"tiktok_e-commerce/app/product/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
