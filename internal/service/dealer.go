package service

import (
	"context"
	"dealer/internal/dao"
	"dealer/internal/models"
	"errors"

	"gorm.io/gorm"
)

type DealerInterface interface {
	ProcessOrder(context.Context, *models.Order) error
}

type Dealer struct {
	db               *gorm.DB
	orderDAO         dao.OrderInterface
	dealDAO          dao.DealInterface
	buyBook          OrderBookInterface
	sellBook         OrderBookInterface
	lastTradingPrice float64
}

var _ (DealerInterface) = (*Dealer)(nil)

func NewDealer(db *gorm.DB, orderDAO dao.OrderInterface, dealDAO dao.DealInterface) *Dealer {
	return &Dealer{
		db:       db,
		orderDAO: orderDAO,
		dealDAO:  dealDAO,
		buyBook:  NewOrderBook(BuyComparator),
		sellBook: NewOrderBook(SellComparator),
	}
}

func (d *Dealer) ProcessOrder(ctx context.Context, order *models.Order) error {
	if order.IsCancel {
		d.buyBook.RemoveOrder(order.ID)
		d.sellBook.RemoveOrder(order.ID)
		return nil
	}

	switch order.OrderType {
	case models.OrderTypeBuy:
		return d.processOrder(ctx, order, d.sellBook, d.buyBook)
	case models.OrderTypeSell:
		return d.processOrder(ctx, order, d.buyBook, d.sellBook)
	default:
		return errors.New("invalid order type")
	}
}

func (d *Dealer) processOrder(ctx context.Context, takerOrder *models.Order, makerBook, takerBook OrderBookInterface) error {
	var deals []*models.Deal
	var updateOrders []*models.Order
	for {
		makerOrder := makerBook.Peek()
		if makerOrder == nil {
			break
		}

		var price float64
		switch {
		case makerOrder.PriceType == models.PriceTypeLimit:
			price = makerOrder.Price
		case makerOrder.PriceType == models.PriceTypeMarket:
			price = d.lastTradingPrice
		}

		if !isPriceMatch(takerOrder, price) {
			break
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
			makerBook.Dequeue()
		}

		if takerOrder.RemainQuantity == 0 {
			break
		}
	}

	updateOrders = append(updateOrders, takerOrder)
	if takerOrder.RemainQuantity > 0 {
		takerBook.AddOrder(takerOrder)
	}

	return d.recordDeal(ctx, deals, updateOrders)
}

func isPriceMatch(takerOrder *models.Order, price float64) bool {
	if takerOrder.PriceType == models.PriceTypeLimit {
		switch {
		case takerOrder.OrderType == models.OrderTypeBuy && takerOrder.Price < price:
			return false
		case takerOrder.OrderType == models.OrderTypeSell && takerOrder.Price > price:
			return false
		}
	}

	return true
}

func (d *Dealer) recordDeal(ctx context.Context, deals []*models.Deal, orders []*models.Order) error {
	tx := d.db.Begin()
	if len(orders) != 0 {
		if err := d.orderDAO.BulkUpdate(ctx, tx, orders); err != nil {
			tx.Rollback()
			return err
		}
	}

	if len(deals) != 0 {
		if err := d.dealDAO.Insert(ctx, tx, deals); err != nil {
			tx.Rollback()
			return err
		}
	}

	return tx.Commit().Error
}
