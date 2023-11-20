package main

import "crypto-exchange/orderbook"

func getOrderView(order *orderbook.Order) *OrderView {
	return &OrderView{
		ID:        order.ID,
		Price:     order.Limit.Price,
		Size:      order.Size,
		Bid:       order.Bid,
		Timestamp: order.Timestamp,
	}
}

func getMatchView(match *orderbook.Match) *MatchView {
	return &MatchView{
		IDAsk: match.Ask.ID,
		IDBid: match.Bid.ID,
		Size:  match.SizeFilled,
		Price: match.Price,
	}
}
