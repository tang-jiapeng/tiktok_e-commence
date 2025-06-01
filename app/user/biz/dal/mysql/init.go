package mysql

import (
	"tiktok_e-commerce/user/biz/model"
	"tiktok_e-commerce/user/conf"

	"github.com/cloudwego/kitex/pkg/klog"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	err error
)

// func Init() {
// 	dsn := conf.GetConf().MySQL.DSN
// 	DB, err = gorm.Open(mysql.Open(dsn),
// 		&gorm.Config{
// 			PrepareStmt:    true,
// 			TranslateError: true,
// 		},
// 	)
// 	if err != nil {
// 		panic(err)
// 	}

// 	err := DB.AutoMigrate(&model.User{})	
// 	if err != nil {
// 		panic(err)
// 	}
// }


func Init() {
	dsn := conf.GetConf().MySQL.DSN
	klog.Info("Initializing MySQL with DSN: %s", dsn)
	if dsn == "" {
			klog.Fatal("MySQL DSN is empty")
	}
	DB, err = gorm.Open(mysql.Open(dsn),
			&gorm.Config{
					PrepareStmt:    true,
					TranslateError: true,
			},
	)
	if err != nil {
			klog.Fatal("failed to connect to MySQL: %v", err)
	}
	sqlDB, err := DB.DB()
	if err != nil {
			klog.Fatal("failed to get sql.DB: %v", err)
	}
	if err := sqlDB.Ping(); err != nil {
			klog.Fatal("failed to ping MySQL: %v", err)
	}
	klog.Info("MySQL connected successfully")
	err = DB.AutoMigrate(&model.User{})
	if err != nil {
			klog.Fatal("failed to migrate User table: %v", err)
	}
	klog.Info("User table migrated successfully")
}