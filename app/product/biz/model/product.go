package model

type Product struct {
	Base
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Picture     string  `json:"picture"`
	Price       string  `json:"price"`
	Stock       int64   `json:"stock"`
	Sale        float32 `json:"sale"`
	PublicState int8    `json:"public_state"`
	LockStock   int64   `json:"lock_stock"`
}

func (u Product) TableName() string {
	return "tb_product"
}
