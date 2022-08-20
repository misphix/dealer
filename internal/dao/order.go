package dao

import (
	"dealer/internal/models"

	"golang.org/x/net/context"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type OrderInterface interface {
	Insert(context.Context, *gorm.DB, *models.Order) error
	List(context.Context, *gorm.DB, *models.Order) ([]*models.Order, error)
	TakeAndLock(context.Context, *gorm.DB, *models.Order) (*models.Order, error)
	Update(context.Context, *gorm.DB, *models.Order, *models.Order) error
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

func (o *Order) List(ctx context.Context, tx *gorm.DB, cond *models.Order) ([]*models.Order, error) {
	var orders []*models.Order
	if err := tx.WithContext(ctx).Find(&orders, cond).Error; err != nil {
		return nil, err
	}

	return orders, nil
}

func (d *Order) TakeAndLock(ctx context.Context, tx *gorm.DB, cond *models.Order) (*models.Order, error) {
	var order *models.Order
	cmd := tx.WithContext(ctx).Clauses(clause.Locking{Strength: "UPDATE"}).Take(&order, cond)
	if err := cmd.Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, err
	}

	return order, nil
}

func (d *Order) Update(ctx context.Context, tx *gorm.DB, data, cond *models.Order) error {
	return tx.WithContext(ctx).Where(cond).Updates(data).Error
}

func (d *Order) BulkUpdate(ctx context.Context, tx *gorm.DB, orders []*models.Order) error {
	return tx.WithContext(ctx).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"remain_quantity"}),
		}).Create(&orders).
		Error
}
