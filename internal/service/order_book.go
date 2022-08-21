package service

import (
	"dealer/internal/models"
	"math"
	"sort"
)

const (
	TOLERANCE = 0.000001
)

type Comparator func(*models.Order, *models.Order) bool

var BuyComparator Comparator = func(o1, o2 *models.Order) bool {
	if int(o1.PriceType) != int(o2.PriceType) {
		return int(o1.PriceType) < int(o2.PriceType)
	}

	diff := o1.Price - o2.Price
	if o1.PriceType == models.PriceTypeLimit && math.Abs(diff) > TOLERANCE {
		return diff < 0
	}

	return o1.ID > o2.ID
}

var SellComparator Comparator = func(o1, o2 *models.Order) bool {
	if int(o1.PriceType) != int(o2.PriceType) {
		return int(o1.PriceType) < int(o2.PriceType)
	}

	diff := o1.Price - o2.Price
	if o1.PriceType == models.PriceTypeLimit && math.Abs(diff) > TOLERANCE {
		return diff > 0
	}

	return o1.ID > o2.ID
}

type OrderBookInterface interface {
	AddOrder(*models.Order)
	Peek() *models.Order
	Dequeue() *models.Order
	RemoveOrder(int64)
}

type OrderBook struct {
	orders     []*models.Order
	comparator Comparator
}

var _ OrderBookInterface = (*OrderBook)(nil)

func NewOrderBook(comparator Comparator) *OrderBook {
	return &OrderBook{
		comparator: comparator,
	}
}

func (book *OrderBook) AddOrder(order *models.Order) {
	book.orders = append(book.orders, order)
	sort.Slice(book.orders, func(i, j int) bool {
		return book.comparator(book.orders[i], book.orders[j])
	})
}

func (book *OrderBook) Peek() *models.Order {
	length := len(book.orders)
	if length == 0 {
		return nil
	}

	return book.orders[length-1]
}

func (book *OrderBook) Dequeue() *models.Order {
	length := len(book.orders)
	if length == 0 {
		return nil
	}

	order := book.orders[length-1]
	book.remove(length - 1)
	return order
}

func (book *OrderBook) RemoveOrder(orderID int64) {
	for i, order := range book.orders {
		if orderID == order.ID {
			book.remove(i)
			break
		}
	}
}

func (book *OrderBook) remove(index int) {
	book.orders = append(book.orders[:index], book.orders[index+1:]...)
}
