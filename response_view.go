package main

import "crypto-exchange/orderbook"

type OrderbookData struct {
	BidsTotalVolume float64
	AsksTotalVolume float64
	Asks            []*OrderView
	Bids            []*OrderView
}

type OrderView struct {
	ID        int64
	Price     float64
	Size      float64
	UserID    int64
	Bid       bool
	Timestamp int64
}

type MatchView struct {
	IDBid int64
	IDAsk int64
	Size  float64
	Price float64
}

func getOrderView(order *orderbook.Order) *OrderView {
	return &OrderView{
		ID:        order.ID,
		UserID:    order.UserID,
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
