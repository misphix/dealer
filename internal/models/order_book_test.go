package models

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_AddOrder(t *testing.T) {
	tests := []struct {
		name      string
		orderBook *OrderBook
		orders    []*Order
		expected  []*Order
	}{
		{
			name:      "Sell limit price lower price first",
			orderBook: NewOrderBook(SellComparator),
			orders: []*Order{
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: PriceTypeLimit,
					Price:     3,
				},
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     1,
				},
			},
			expected: []*Order{
				{
					ID:        2,
					PriceType: PriceTypeLimit,
					Price:     3,
				},
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     1,
				},
			},
		},
		{
			name:      "Sell market price market price first",
			orderBook: NewOrderBook(SellComparator),
			orders: []*Order{
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: PriceTypeMarket,
				},
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     1,
				},
			},
			expected: []*Order{
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     1,
				},
				{
					ID:        2,
					PriceType: PriceTypeMarket,
				},
			},
		},
		{
			name:      "Sell limit price eariler first",
			orderBook: NewOrderBook(SellComparator),
			orders: []*Order{
				{
					ID:        2,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
			},
			expected: []*Order{
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
			},
		},
		{
			name:      "Sell market price eariler first",
			orderBook: NewOrderBook(SellComparator),
			orders: []*Order{
				{
					ID:        2,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
			},
			expected: []*Order{
				{
					ID:        3,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
			},
		},
		{
			name:      "Buy limit price higher price first",
			orderBook: NewOrderBook(BuyComparator),
			orders: []*Order{
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: PriceTypeLimit,
					Price:     3,
				},
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     1,
				},
			},
			expected: []*Order{
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     1,
				},
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				}, {
					ID:        2,
					PriceType: PriceTypeLimit,
					Price:     3,
				},
			},
		},
		{
			name:      "Buy limit price eariler first",
			orderBook: NewOrderBook(BuyComparator),
			orders: []*Order{
				{
					ID:        2,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
			},
			expected: []*Order{
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
			},
		},
		{
			name:      "Buy market price eariler first",
			orderBook: NewOrderBook(BuyComparator),
			orders: []*Order{
				{
					ID:        2,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        3,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
			},
			expected: []*Order{
				{
					ID:        3,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
				{
					ID:        1,
					PriceType: PriceTypeMarket,
					Price:     2,
				},
			},
		},
		{
			name:      "Buy market price market price first",
			orderBook: NewOrderBook(BuyComparator),
			orders: []*Order{
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: PriceTypeMarket,
				},
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     1,
				},
			},
			expected: []*Order{
				{
					ID:        3,
					PriceType: PriceTypeLimit,
					Price:     1,
				},
				{
					ID:        1,
					PriceType: PriceTypeLimit,
					Price:     2,
				},
				{
					ID:        2,
					PriceType: PriceTypeMarket,
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, order := range test.orders {
				test.orderBook.AddOrder(order)
			}

			assert.Equal(t, test.expected, test.orderBook.Orders)
		})
	}
}
