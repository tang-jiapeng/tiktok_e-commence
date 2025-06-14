package model

import (
	"context"

	"gorm.io/gorm"
)

type Category struct {
	Base
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (p *Category) TableName() string {
	return "tb_category"
}

func SelectCategory(db *gorm.DB, ctx context.Context, id int64) (category Category, err error) {
	result := db.WithContext(ctx).Where("id=?", id).First(&category)
	err = result.Error
	return
}

func CreateCategory(db *gorm.DB, ctx context.Context, category *Category) (err error) {
	result := db.WithContext(ctx).Create(&category)
	err = result.Error
	return
}

func DeleteCategory(db *gorm.DB, ctx context.Context, id int64) (err error) {
	result := db.WithContext(ctx).Where("id=?", id).Delete(&Category{})
	err = result.Error
	return
}

func UpdateCategory(db *gorm.DB, ctx context.Context, category *Category) (err error) {
	result := db.WithContext(ctx).Updates(&category)
	err = result.Error
	return
}
