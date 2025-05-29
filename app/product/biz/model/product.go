package model

import (
	"fmt"
	"time"
)

type Product struct {
	ID          uint32 `gorm:"primarykey"`
	StoreId     uint32 `gorm:"index:store_id;not null"`
	Name        string `gorm:"index:query,class:FULLTEXT;index:name;type:varchar(64);not null"`
	Description string `gorm:"index:query,class:FULLTEXT;type:TEXT"`
	Picture     string
	Price       uint32     `gorm:"not null;default:0"` // use integer to avoid accuracy loss; represent float with a decimal precision of 2
	Stock       uint32     `gorm:"not null;default:0"`
	Categories  []Category `gorm:"many2many:product_category;"`
	CreatedAt   time.Time  `gorm:"column:created_at;not null"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;not null"`
}

func (Product) TableName() string {
	return "product"
}

// GetCacheKey 返回产品缓存键，格式为 "product:id:%d"
func (p *Product) GetCacheKey() string {
	return fmt.Sprintf("%s:id:%d", p.TableName(), p.ID)
}
