package mysql

import (
	"tiktok_e-commerce/payment/biz/model"
	"tiktok_e-commerce/payment/conf"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)


func Init() {
	dsn := conf.GetConf().MySQL.DSN
	DB, err = gorm.Open(mysql.Open(dsn),
		&gorm.Config{
			PrepareStmt:    true,
			TranslateError: true,
		},
	)
	if err != nil {
		panic(err)
	}

	err := DB.AutoMigrate(&model.Payment{})	
	if err != nil {
		panic(err)
	}
}

