package model

import (
	"context"

	"gorm.io/gorm"
)

type Product struct {
	Base
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Picture     string  `json:"picture"`
	Price       float32 `json:"price"`
	Stock       int64   `json:"stock"`
	Sale        int64   `json:"sale"`
	PublicState uint32  `json:"public_state"`
	LockStock   int64   `json:"lock_stock"`
	CategoryId  int64   `json:"category_id"`
	BrandId     int64   `json:"brand_id"`
}

func (u Product) TableName() string {
	return "tb_product"
}

func SelectProduct(db *gorm.DB, ctx context.Context, id int64) (product Product, err error) {
	product = Product{}
	result := db.WithContext(ctx).Where("id = ?", id).First(&product)
	err = result.Error
	return
}

func UpdateProduct(db *gorm.DB, ctx context.Context, product *Product) (err error) {
	result := db.WithContext(ctx).Updates(&product)
	err = result.Error
	return
}

func DeleteProduct(db *gorm.DB, ctx context.Context, id int64) (err error) {
	result := db.WithContext(ctx).Delete(&Product{Base: Base{ID: id}})
	err = result.Error
	return
}

func CreateProduct(db *gorm.DB, ctx context.Context, product *Product) (err error) {
	result := db.WithContext(ctx).Create(product)
	err = result.Error
	return
}
