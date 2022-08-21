package service

import (
	"dealer/internal/models"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAddOrder(t *testing.T) {
	tests := []struct {
		name      string
		orderBook *OrderBook
		orders    []*models.Order
		expected  []*models.Order
	}{
		{
			name:      "Sell limit price lower price first",
			orderBook: NewOrderBook(SellComparator),
			orders: []*models.Order{
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: models.PriceTypeLimit,
					Price:     3,
				},
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     1,
				},
			},
			expected: []*models.Order{
				{
					ID:        2,
					PriceType: models.PriceTypeLimit,
					Price:     3,
				},
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     1,
				},
			},
		},
		{
			name:      "Sell market price market price first",
			orderBook: NewOrderBook(SellComparator),
			orders: []*models.Order{
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: models.PriceTypeMarket,
				},
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     1,
				},
			},
			expected: []*models.Order{
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     1,
				},
				{
					ID:        2,
					PriceType: models.PriceTypeMarket,
				},
			},
		},
		{
			name:      "Sell limit price eariler first",
			orderBook: NewOrderBook(SellComparator),
			orders: []*models.Order{
				{
					ID:        2,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
			},
			expected: []*models.Order{
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
			},
		},
		{
			name:      "Sell market price eariler first",
			orderBook: NewOrderBook(SellComparator),
			orders: []*models.Order{
				{
					ID:        2,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
			},
			expected: []*models.Order{
				{
					ID:        3,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
			},
		},
		{
			name:      "Buy limit price higher price first",
			orderBook: NewOrderBook(BuyComparator),
			orders: []*models.Order{
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: models.PriceTypeLimit,
					Price:     3,
				},
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     1,
				},
			},
			expected: []*models.Order{
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     1,
				},
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				}, {
					ID:        2,
					PriceType: models.PriceTypeLimit,
					Price:     3,
				},
			},
		},
		{
			name:      "Buy limit price eariler first",
			orderBook: NewOrderBook(BuyComparator),
			orders: []*models.Order{
				{
					ID:        2,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
			},
			expected: []*models.Order{
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
			},
		},
		{
			name:      "Buy market price eariler first",
			orderBook: NewOrderBook(BuyComparator),
			orders: []*models.Order{
				{
					ID:        2,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
			},
			expected: []*models.Order{
				{
					ID:        3,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: models.PriceTypeMarket,
					Price:     2,
				},
			},
		},
		{
			name:      "Buy market price market price first",
			orderBook: NewOrderBook(BuyComparator),
			orders: []*models.Order{
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: models.PriceTypeMarket,
				},
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     1,
				},
			},
			expected: []*models.Order{
				{
					ID:        3,
					PriceType: models.PriceTypeLimit,
					Price:     1,
				},
				{
					ID:        1,
					PriceType: models.PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: models.PriceTypeMarket,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, order := range test.orders {
				test.orderBook.AddOrder(order)
			}

			assert.Equal(t, test.expected, test.orderBook.orders)
		})
	}
}

func TestPeek(t *testing.T) {
	tests := []struct {
		name      string
		orderBook *OrderBook
		expected  *models.Order
	}{
		{
			name: "Peek order",
			orderBook: &OrderBook{
				orders: []*models.Order{
					{ID: 1},
					{ID: 2},
				},
			},
			expected: &models.Order{ID: 2},
		},
		{
			name:      "Peek order nil",
			orderBook: &OrderBook{},
			expected:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.orderBook.Peek()
			assert.Equal(t, test.expected, actual)
		})
	}
}

func TestDequeue(t *testing.T) {
	tests := []struct {
		name      string
		orderBook *OrderBook
		length    int
		expected  *models.Order
	}{
		{
			name: "Dequeue order",
			orderBook: &OrderBook{
				orders: []*models.Order{
					{ID: 1},
					{ID: 2},
				},
			},
			length:   1,
			expected: &models.Order{ID: 2},
		},
		{
			name:      "Dequeue order nil",
			orderBook: &OrderBook{},
			length:    0,
			expected:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual := test.orderBook.Dequeue()
			assert.Equal(t, test.expected, actual)
			assert.Len(t, test.orderBook.orders, test.length)
		})
	}
}

func TestRemoveOrder(t *testing.T) {
	tests := []struct {
		name      string
		orderBook *OrderBook
		expected  []*models.Order
	}{
		{
			name: "Remove order",
			orderBook: &OrderBook{
				orders: []*models.Order{
					{ID: 1},
					{ID: 2},
				},
			},
			expected: []*models.Order{
				{ID: 2},
			},
		},
		{
			name:      "Remove no order",
			orderBook: &OrderBook{},
			expected:  nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			test.orderBook.RemoveOrder(1)
			assert.Equal(t, test.expected, test.orderBook.orders)
		})
	}
}
