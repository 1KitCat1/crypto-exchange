package orderbook

import "time"

type Order struct {
	ID        int64
	UserID    int64
	Size      float64
	Bid       bool // limit or market
	Limit     *Limit
	Timestamp int64
}

type Orders []*Order

func NewOrder(is_bid bool, size float64) *Order {
	return &Order{
		Size:      size,
		Bid:       is_bid,
		Timestamp: time.Now().UnixNano(),
	}
}

type Limit struct { // group of orders at a certain price limit
	Price  float64
	Orders Orders
	Volume float64
}

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func (limit *Limit) AddOrder(order *Order) {
	order.Limit = limit
	limit.Orders = append(limit.Orders, order)
	limit.Volume += order.Size
}

type Orderbook struct {
	Asks []*Limit
	Bids []*Limit
}

func (o Orders) Len() int { return len(o) }
