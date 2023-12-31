package main

import "crypto-exchange/orderbook"

type Exchange struct {
	orderbooks map[orderbook.Market]*orderbook.Orderbook
	orders     map[int64]*orderbook.Order
	userOrders map[int64]map[int64]struct{}
}

func NewExchange() *Exchange {
	orderbooks := make(map[orderbook.Market]*orderbook.Orderbook)

	for _, market := range SUPPORTED_MARKETS {
		orderbooks[market] = orderbook.NewOrderbook()
	}

	return &Exchange{
		orderbooks: orderbooks,
		orders:     make(map[int64]*orderbook.Order),
		userOrders: make(map[int64]map[int64]struct{}),
	}
}
