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
	RealStock   int64   `gorm:"-" json:"quantity"`
}

func (p Product) TableName() string {
	return "tb_product"
}

func (p *Product) AfterFind(tx *gorm.DB) (err error) {
	p.RealStock = p.Stock - p.LockStock
	return nil
}

func SelectProduct(db *gorm.DB, ctx context.Context, id int64) (product Product, err error) {
	product = Product{}
	result := db.WithContext(ctx).Where("id = ?", id).First(&product)
	err = result.Error
	return product, err
}

func UpdateProduct(db *gorm.DB, ctx context.Context, product *Product) (err error) {
	result := db.WithContext(ctx).Updates(&product)
	err = result.Error
	return err
}

func DeleteProduct(db *gorm.DB, ctx context.Context, id int64) (err error) {
	result := db.WithContext(ctx).Delete(&Product{Base: Base{ID: id}})
	err = result.Error
	return err
}

func CreateProduct(db *gorm.DB, ctx context.Context, product *Product) (err error) {
	result := db.WithContext(ctx).Create(product)
	err = result.Error
	return err
}

func SelectProductList(db *gorm.DB, ctx context.Context, ids []int64) (product []Product, err error) {
	product = []Product{}
	result := db.WithContext(ctx).Where("id IN ?", ids).Find(&product)
	err = result.Error
	return product, err
}

func SelectProductAll(db *gorm.DB, ctx context.Context, index, total int64) (product []Product, err error) {
	product = []Product{}
	result := db.WithContext(ctx).Model(&Product{}).Where("id%? = ?", total, index).Find(&product)
	err = result.Error
	return product, err
}

func UpdateLockStock(db *gorm.DB, ctx context.Context, productQuantityMap map[int64]int64) (err error) {
	err = db.Transaction(func(tx *gorm.DB) error {
		for id, quantity := range productQuantityMap {
			result := tx.WithContext(ctx).
				Model(&Product{}).
				Where("id =?", id).
				Where("stock >= lock_stock + ?", quantity).
				Update("lock_stock", gorm.Expr("lock_stock + ?", quantity))
			e := result.Error
			if e != nil {
				return e
			}
		}
		return nil
	})
	return err
}
