package dal

import (
	"tiktok_e-commerce/app/user/biz/dal/mysql"
	"tiktok_e-commerce/app/user/biz/dal/redis"
)

func Init() {
	redis.Init()
	mysql.Init()
}
