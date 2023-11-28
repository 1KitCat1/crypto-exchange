package main

import "crypto-exchange/orderbook"

const MarketETH orderbook.Market = "ETH"

const (
	MarketOrder OrderType = "MARKET"
	LimitOrder  OrderType = "LIMIT"
)
