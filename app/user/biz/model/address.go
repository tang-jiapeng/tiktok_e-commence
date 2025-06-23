package model

import (
	"context"

	"gorm.io/gorm"
)

const (
	AddressDefaultStatusNo  = 0
	AddressDefaultStatusYes = 1
)

type Address struct {
	Base
	UserId        int32  `gorm:"not null;type:int;index:idx_user_id_default_status"`
	Name          string `gorm:"not null;type:varchar(64);comment:收件人姓名"`
	PhoneNumber   string `gorm:"not null;type:varchar(64);comment:收件人手机号"`
	DefaultStatus int8   `gorm:"not null;type:int(1);default:0;index:idx_user_id_default_status;comment:是否默认地址，0-否，1-是"`
	Province      string `gorm:"not null;type:varchar(64);comment:省"`
	City          string `gorm:"not null;type:varchar(64);comment:市"`
	Region        string `gorm:"not null;type:varchar(64);comment:区"`
	DetailAddress string `gorm:"not null;type:varchar(256);comment:详细地址"`
}

type AddressInfo struct {
	Name          string `json:"name"`
	PhoneNumber   string `json:"phone_number"`
	DefaultStatus int8   `json:"default_status"`
	Province      string `json:"province"`
	City          string `json:"city"`
	Region        string `json:"region"`
	DetailAddress string `json:"detail_address"`
}

func (a Address) TableName() string {
	return "tb_receive_address"
}

func CreateAddress(db *gorm.DB, ctx context.Context, address *Address) (err error) {
	result := db.WithContext(ctx).Create(address)
	err = result.Error
	return
}

func ExistDefaultAddress(db *gorm.DB, ctx context.Context, userId int32) (address Address, err error) {
	result := db.WithContext(ctx).
		Where(&Address{UserId: userId, DefaultStatus: AddressDefaultStatusYes}).First(&address)
	err = result.Error
	return
}

func UpdateAddress(db *gorm.DB, ctx context.Context, address Address) (err error) {
	result := db.WithContext(ctx).Save(&address)
	err = result.Error
	return
}

func GetAdressList(db *gorm.DB, ctx context.Context, userId int32) (addressList []AddressInfo, err error) {
	err = db.WithContext(ctx).Model(Address{}).
		Select("name, phone_number, default_status, province, city, region, detail_address").
		Where(Address{UserId: userId}).Find(&addressList).Error
	return
}
