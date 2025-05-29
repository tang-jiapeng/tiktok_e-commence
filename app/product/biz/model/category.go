package model

import (
	"fmt"
	"time"
)

type Category struct {
	ID        uint32    `gorm:"primarykey"`
	Name      string    `gorm:"uniqueIndex;type:varchar(32);not null"`
	Products  []Product `gorm:"many2many:product_category"`
	CreatedAt time.Time `gorm:"column:created_at;not null"`
	UpdatedAt time.Time `gorm:"column:updated_at;not null"`
}

func (Category) TableName() string {
	return "category"
}

// GetCacheKey 返回产品缓存键，格式为 "category:name:%s"
func (c *Category) GetCacheKey() string {
	return fmt.Sprintf("%s:name:%s", c.TableName(), c.Name)
}
