package model

import (
	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type CartItem struct {
	Base
	UserId    int32 `gorm:"not null;type:int;index:idx_user_id"`
	ProductId int32 `gorm:"not null;type:int;"`
	Quantity  int32 `gorm:"not null;type:int;"`
}

func (c CartItem) TableName() string {
	return "tb_cart_items"
}

func AddCartItem(db *gorm.DB, ctx context.Context, item *CartItem) error {
	return db.WithContext(ctx).Create(item).Error
}

func GetCartItemByUserId(db *gorm.DB, ctx context.Context, userId int32) ([]*CartItem, error) {
	var items []*CartItem
	err := db.WithContext(ctx).Where(&CartItem{UserId: userId}).Find(&items).Error
	return items, err
}

func EmptyCart(db *gorm.DB, ctx context.Context, userId int32) error {
	err := db.WithContext(ctx).Where(&CartItem{UserId: userId}).Delete(&CartItem{}).Error
	return err
}
