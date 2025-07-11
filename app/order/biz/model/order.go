package model

import (
	"context"
	"time"

	"gorm.io/gorm"
)

const (
	OrderStatusUnpaid = iota
	OrderStatusPaid
	OrderStatusCanceled
)

type Order struct {
	OrderID       string  `gorm:"primarykey;type:varchar(64)"`
	UserID        int32   `gorm:"not null;type:int;index:idx_user_id"`
	TotalCost     float32 `gorm:"not null;type:decimal(10,2)"`
	Name          string  `gorm:"not null;type:varchar(64)"`
	PhoneNumber   string  `gorm:"not null;type:char(11)"`
	Province      string  `gorm:"not null;type:varchar(64)"`
	City          string  `gorm:"not null;type:varchar(64)"`
	Region        string  `gorm:"not null;type:varchar(64)"`
	DetailAddress string  `gorm:"not null;type:varchar(255)"`
	Status        int32   `gorm:"not null;type:int;default:0"`
	CreatedAt     time.Time
	OrderItems    []OrderItem `gorm:"foreignKey:OrderID;references:OrderID"`
}

type OrderInfo struct {
	Order
	orderItems []OrderItem
}

func (o *Order) TableName() string {
	return "tb_order"
}

func CreateOrder(ctx context.Context, db *gorm.DB, order *Order) error {
	return db.WithContext(ctx).Create(order).Error
}

func GetOrdersByUserId(ctx context.Context, db *gorm.DB, userId int32) (orderList []Order, err error) {
	err = db.WithContext(ctx).Model(&Order{}).
		Where(&Order{UserID: userId}).
		Preload("OrderItems").
		Find(&orderList).Error
	return
}

func CheckOrderUnpaid(ctx context.Context, db *gorm.DB, orderId string) (bool, error) {
	var order Order
	err := db.WithContext(ctx).Model(&Order{}).Select("status").
		Where(&Order{OrderID: orderId}).First(&order).Error
	if err != nil {
		return false, err
	} else {
		return order.Status == OrderStatusUnpaid, nil
	}
}

func CancelOrder(ctx context.Context, db *gorm.DB, orderId string) (int64, error) {
	tx := db.WithContext(ctx).Model(&Order{}).
		Where(&Order{OrderID: orderId, Status: OrderStatusUnpaid}).
		Update("status", OrderStatusCanceled)
	return tx.RowsAffected, tx.Error
}

func CancelOrderList(ctx context.Context, db *gorm.DB, orderIdList []string) (int64, error) {
	tx := db.WithContext(ctx).Model(&Order{}).
		Where("order_id IN ? AND status = ?", orderIdList, OrderStatusUnpaid).
		Update("status", OrderStatusCanceled)
	return tx.RowsAffected, tx.Error
}

func GetOverdueOrder(ctx context.Context, db *gorm.DB, placeTime time.Time) (orderIdList []string, err error) {
	err = db.WithContext(ctx).Model(&Order{}).
		Where("status = ? and created_at < ?", OrderStatusUnpaid, placeTime).
		Order("created_at ASC").
		Pluck("order_id", &orderIdList).Error
	return
}

func MarkOrderPaid(ctx context.Context, db *gorm.DB, orderId string) (int64, error) {
	tx := db.WithContext(ctx).Model(&Order{}).
		Where(&Order{OrderID: orderId, Status: OrderStatusUnpaid}).
		Update("status", OrderStatusPaid)
	return tx.RowsAffected, tx.Error
}

func GetOrder(ctx context.Context, db *gorm.DB, orderId string) (order *Order, err error) {
	err = db.WithContext(ctx).Model(&Order{}).
		Where(&Order{OrderID: orderId}).
		Preload("OrderItems").
		First(&order).Error
	return
}
