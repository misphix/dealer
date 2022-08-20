package dao

import (
	"dealer/internal/models"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type DealInterface interface {
	Insert(context.Context, *gorm.DB, []*models.Deal) error
	List(context.Context, *gorm.DB, *models.Deal) ([]*models.Deal, error)
}

type Deal struct {
}

var _ DealInterface = (*Deal)(nil)

func NewDeal() *Deal {
	return &Deal{}
}

func (d *Deal) Insert(ctx context.Context, tx *gorm.DB, deals []*models.Deal) error {
	if len(deals) == 0 {
		return nil
	}

	return tx.WithContext(ctx).Create(&deals).Error
}

func (d *Deal) List(ctx context.Context, tx *gorm.DB, cond *models.Deal) ([]*models.Deal, error) {
	var deals []*models.Deal
	if err := tx.WithContext(ctx).Find(&deals, cond).Error; err != nil {
		return nil, err
	}

	return deals, nil
}
