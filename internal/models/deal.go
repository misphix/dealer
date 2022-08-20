package models

import "gorm.io/gorm/schema"

type Deal struct {
	ID           int64   `gorm:"primaryKey;column:id" json:"id"`
	TakerOrderID int64   `gorm:"column:taker_order_id"`
	MakerOrderID int64   `gorm:"column:maker_order_id"`
	Quantity     uint    `gorm:"column:quantity"`
	Price        float64 `gorm:"column:price"`
}

var _ schema.Tabler = (*Deal)(nil)

func (Deal) TableName() string {
	return "deal"
}
