package dal

import (
	"tiktok_e-commerce/auth/biz/dal/mysql"
	"tiktok_e-commerce/auth/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
