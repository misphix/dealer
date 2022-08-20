package models

import (
	"math"
	"sort"
)

const (
	TOLERANCE = 0.000001
)

type Comparator func(*Order, *Order) bool

var BuyComparator Comparator = func(o1, o2 *Order) bool {
	if int(o1.PriceType) != int(o2.PriceType) {
		return int(o1.PriceType) < int(o2.PriceType)
	}

	diff := o1.Price - o2.Price
	if o1.PriceType == PriceTypeLimit && math.Abs(diff) > TOLERANCE {
		return diff < 0
	}

	return o1.ID > o2.ID
}

var SellComparator Comparator = func(o1, o2 *Order) bool {
	if int(o1.PriceType) != int(o2.PriceType) {
		return int(o1.PriceType) < int(o2.PriceType)
	}

	diff := o1.Price - o2.Price
	if o1.PriceType == PriceTypeLimit && math.Abs(diff) > TOLERANCE {
		return diff > 0
	}

	return o1.ID > o2.ID
}

type OrderBook struct {
	Orders     []*Order
	comparator Comparator
}

func NewOrderBook(comparator Comparator) *OrderBook {
	return &OrderBook{
		comparator: comparator,
	}
}

func (book *OrderBook) AddOrder(order *Order) {
	book.Orders = append(book.Orders, order)
	sort.Slice(book.Orders, func(i, j int) bool {
		return book.comparator(book.Orders[i], book.Orders[j])
	})
}

func (book *OrderBook) Remove(index int) {
	book.Orders = append(book.Orders[:index], book.Orders[index+1:]...)
}
