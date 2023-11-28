package main

import "crypto-exchange/orderbook"

type Exchange struct {
	orderbooks map[orderbook.Market]*orderbook.Orderbook
}

func NewExchange() *Exchange {
	orderbooks := make(map[orderbook.Market]*orderbook.Orderbook)

	for _, market := range SUPPORTED_MARKETS {
		orderbooks[market] = orderbook.NewOrderbook()
	}

	return &Exchange{
		orderbooks: orderbooks,
	}
}
