package model

import (
	"context"

	"gorm.io/gorm"
)

type Brand struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

func (p *Brand) TableName() string {
	return "tb_brand"
}

func SelectBrand(db *gorm.DB, ctx context.Context, id int64) (brand Brand, err error) {
	result := db.WithContext(ctx).Where("id = ?", id).First(&brand)
	err = result.Error
	return
}

func CreateBrand(db *gorm.DB, ctx context.Context, brand *Brand) (err error) {
	result := db.WithContext(ctx).Create(&brand)
	err = result.Error
	return
}

func UpdateBrand(db *gorm.DB, ctx context.Context, brand *Brand) (err error) {
	result := db.WithContext(ctx).Updates(&brand)
	err = result.Error
	return
}

func DeleteBrand(db *gorm.DB, ctx context.Context, id int64) (err error) {
	result := db.WithContext(ctx).Where("id=?", id).Delete(&Brand{})
	err = result.Error
	return
}
