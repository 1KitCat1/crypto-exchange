package main

import "crypto-exchange/orderbook"

const AUTH_SERVICE_URL = "a"
const AUTH_SERVICE_ENDPOINT = "b"

const MarketETH orderbook.Market = "ETH"

const (
	MarketOrder OrderType = "MARKET"
	LimitOrder  OrderType = "LIMIT"
)
