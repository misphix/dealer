package dao

import (
	"dealer/internal/models"

	"golang.org/x/net/context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderInterface interface {
	Insert(context.Context, *gorm.DB, *models.Order) error
	Update(context.Context, *gorm.DB, *models.Order) error
	BulkUpdate(context.Context, *gorm.DB, []*models.Order) error
}

type Order struct{}

var _ OrderInterface = (*Order)(nil)

func NewOrder() *Order {
	return &Order{}
}

func (o *Order) Insert(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	if order == nil {
		return nil
	}

	return tx.WithContext(ctx).Create(&order).Error
}

func (d *Order) Update(ctx context.Context, tx *gorm.DB, order *models.Order) error {
	return tx.WithContext(ctx).Updates(&order).Error
}

func (d *Order) BulkUpdate(ctx context.Context, tx *gorm.DB, orders []*models.Order) error {
	return tx.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"remain_quantity"}),
		}).Create(&orders).
		Error
}
