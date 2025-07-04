package mysql

import (
	"tiktok_e-commerce/cart/biz/model"
	"tiktok_e-commerce/cart/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

func Init() {
	DB, err = gorm.Open(mysql.Open(conf.GetConf().MySQL.DSN),
		&gorm.Config{
			PrepareStmt:    true,
			TranslateError: true,
		},
	)
	err = DB.AutoMigrate(&model.CartItem{})
	if err != nil {
		panic(err)
	}
}
