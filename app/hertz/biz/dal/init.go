package dal

import (
	"tiktok_e-commerce/app/hertz/biz/dal/mysql"
	"tiktok_e-commerce/app/hertz/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
