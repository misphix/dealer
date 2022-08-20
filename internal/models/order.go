package models

import (
	"encoding"
	"encoding/json"

	"gorm.io/gorm/schema"
)

type OrderType int

const (
	OrderTypeBuy OrderType = iota + 1
	OrderTypeSell
)

type PriceType int

const (
	PriceTypeLimit PriceType = iota + 1
	PriceTypeMarket
)

type Order struct {
	ID             int64     `gorm:"primaryKey;column:id" json:"id"`
	OrderType      OrderType `gorm:"column:order_type" json:"order_type"`
	Quantity       uint      `gorm:"column:quantity" json:"quantity"`
	RemainQuantity uint      `gorm:"column:remain_quantity" json:"remain_quantity"`
	PriceType      PriceType `gorm:"column:price_type" json:"price_type"`
	Price          float64   `gorm:"column:price" json:"price"`
	IsCancel       bool      `gorm:"column:is_cancel"`
}

var _ schema.Tabler = (*Order)(nil)
var _ encoding.BinaryMarshaler = (*Order)(nil)
var _ encoding.BinaryUnmarshaler = (*Order)(nil)

func (Order) TableName() string {
	return "order"
}

func (o Order) MarshalBinary() ([]byte, error) {
	return json.Marshal(o)
}

func (o *Order) UnmarshalBinary(data []byte) error {
	return json.Unmarshal(data, &o)
}
