package model

import (
	"context"

	"gorm.io/gorm"
)

type Payment struct {
	Base
	OrderID int64   `json:"order_id"`
	Amount  float32 `json:"amount"`
}

func Create(db *gorm.DB, ctx context.Context, paymentLog *Payment) error {
	result := db.WithContext(ctx).Create(paymentLog).Error
	return result
}
