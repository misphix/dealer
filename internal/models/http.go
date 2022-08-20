package models

type OrderRequest struct {
	OrderType OrderType `json:"order_type"`
	Quantity  uint      `son:"quantity"`
	PriceType PriceType `json:"price_type"`
	Price     float64   `json:"price"`
}

type CancelOrderRequest struct {
	ID int64 `uri:"id"`
}
