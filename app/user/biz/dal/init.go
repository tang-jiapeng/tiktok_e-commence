package dal

import (
	"tiktok_e-commerce/user/biz/dal/mysql"
	"tiktok_e-commerce/user/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
