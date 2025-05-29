package model

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserId          int64          `gorm:"column:user_id;primaryKey;autoIncrement" redis:"user_id"` // 用户ID
	Username        string         `gorm:"column:user_name;not null" redis:"username"`              // 用户名
	Email           string         `gorm:"column:email;not null" redis:"email"`                     // 用户邮箱
	Password        string         `gorm:"column:password;not null" redis:"password"`               // 用户密码
	Salt            string         `gorm:"column:salt;not null" redis:"salt"`                       //密码盐值
	UserPermissions int32          `gorm:"user_permissions;not null" redis:"user_permissions"`      //用户权限
	Created_at      time.Time      `gorm:"column:created_at;not null" redis:"created_at"`           // 创建时间
	Updated_at      time.Time      `gorm:"column:updated_at;not null" redis:"updated_at"`           // 更新时间
	DeletedAT       gorm.DeletedAt `gorm:"index;column:deleted_at"`
}

type Tabler interface {
	TableName() string
}

func (u *User) TableName() string {
	return "user"
}

// GetCacheKey 返回用户缓存键，格式为 "user:UserId:%d"
func (u *User) GetCacheKey() string {
	return fmt.Sprintf("%s:UserId:%d", u.TableName(), u.UserId)
}