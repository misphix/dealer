package service

import (
	"context"
	"dealer/internal/dao"
	"dealer/internal/models"
	"errors"

	"gorm.io/gorm"
)

type DealerInterface interface {
	ProcessOrder(context.Context, *models.Order) ([]*models.Deal, error)
}

type Dealer struct {
	db               *gorm.DB
	orderDAO         dao.OrderInterface
	dealDAO          dao.DealInterface
	buyBook          *models.OrderBook
	sellBook         *models.OrderBook
	lastTradingPrice float64
}

var _ (DealerInterface) = (*Dealer)(nil)

func NewDealer(db *gorm.DB, orderDAO dao.OrderInterface, dealDAO dao.DealInterface) *Dealer {
	return &Dealer{
		db:       db,
		orderDAO: orderDAO,
		dealDAO:  dealDAO,
		buyBook:  models.NewOrderBook(models.BuyComparator),
		sellBook: models.NewOrderBook(models.SellComparator),
	}
}

func (d *Dealer) ProcessOrder(ctx context.Context, order *models.Order) ([]*models.Deal, error) {
	order.RemainQuantity = order.Quantity
	if err := d.orderDAO.Insert(ctx, d.db, order); err != nil {
		return nil, err
	}

	switch order.OrderType {
	case models.OrderTypeBuy:
		return d.processOrder(ctx, order, d.sellBook, d.buyBook)
	case models.OrderTypeSell:
		return d.processOrder(ctx, order, d.buyBook, d.sellBook)
	default:
		return nil, errors.New("invalid order type")
	}
}

func (d *Dealer) processOrder(ctx context.Context, takerOrder *models.Order, makerBook, takerBook *models.OrderBook) ([]*models.Deal, error) {
	var deals []*models.Deal
	var updateOrders []*models.Order
	for i := len(makerBook.Orders) - 1; i >= 0; i-- {
		makerOrder := makerBook.Orders[i]

		var price float64
		switch {
		case makerOrder.PriceType == models.PriceTypeLimit:
			price = makerOrder.Price
		case makerOrder.PriceType == models.PriceTypeMarket:
			price = d.lastTradingPrice
		}

		if takerOrder.PriceType == models.PriceTypeLimit {
			switch {
			case takerOrder.OrderType == models.OrderTypeBuy && takerOrder.Price < price:
				break
			case takerOrder.OrderType == models.OrderTypeSell && takerOrder.Price > price:
				break
			}
		}

		var quantity uint
		if takerOrder.RemainQuantity > makerOrder.RemainQuantity {
			quantity = makerOrder.RemainQuantity
		} else {
			quantity = takerOrder.RemainQuantity
		}

		d.lastTradingPrice = price
		deal := &models.Deal{
			TakerOrderID: takerOrder.ID,
			MakerOrderID: makerOrder.ID,
			Quantity:     quantity,
			Price:        price,
		}
		deals = append(deals, deal)

		takerOrder.RemainQuantity -= quantity
		makerOrder.RemainQuantity -= quantity
		updateOrders = append(updateOrders, makerOrder)
		if makerOrder.RemainQuantity == 0 {
			makerBook.Remove(i)
		}

		if takerOrder.RemainQuantity == 0 {
			break
		}
	}

	updateOrders = append(updateOrders, takerOrder)
	if takerOrder.RemainQuantity > 0 {
		takerBook.AddOrder(takerOrder)
	}

	if err := d.recordDeal(ctx, deals, updateOrders); err != nil {
		return nil, err
	}

	return deals, nil
}

func (d *Dealer) recordDeal(ctx context.Context, deals []*models.Deal, orders []*models.Order) error {
	tx := d.db.Begin()
	if err := d.orderDAO.BulkUpdate(ctx, tx, orders); err != nil {
		tx.Rollback()
		return err
	}

	if err := d.dealDAO.Insert(ctx, tx, deals); err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}