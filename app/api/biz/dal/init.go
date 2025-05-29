package dal

import (
	"tiktok_e-commerce/api/biz/dal/mysql"
	"tiktok_e-commerce/api/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
