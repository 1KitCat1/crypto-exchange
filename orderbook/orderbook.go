package orderbook

import (
	"fmt"
	"time"
)

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	ID        int64
	UserID    int64
	Size      float64
	Bid       bool // bid or ask
	Limit     *Limit
	Timestamp int64
}

type Orders []*Order

func (order *Order) String() string {
	return fmt.Sprintf("[sizeL %.2f]", order.Size)
}

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

func (limit *Limit) DeleteOrder(order *Order) {
	for i := 0; i < len(limit.Orders); i++ {
		if limit.Orders[i] == order { // TODO: try to do it more efficient
			limit.Orders[i] = limit.Orders[len(limit.Orders)-1]
			limit.Orders = limit.Orders[:len(limit.Orders)-1]
		}
	}
	order.Limit = nil
	limit.Volume -= order.Size
}

type Orderbook struct {
	Asks []*Limit
	Bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		Asks:      []*Limit{},
		Bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (orderbook *Orderbook) PlaceOrder(price float64, order *Order) []Match {
	// TODO: implement matching orders
	//
	if order.Size > 0.0 {
		orderbook.add(price, order)
	}
	return []Match{}
}

func (orderbook *Orderbook) add(price float64, order *Order) {
	var limit *Limit

	if order.Bid {
		limit = orderbook.BidLimits[price]
	} else {
		limit = orderbook.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		if order.Bid {
			orderbook.Bids = append(orderbook.Bids, limit)
			orderbook.BidLimits[price] = limit
		} else {
			orderbook.Asks = append(orderbook.Asks, limit)
			orderbook.AskLimits[price] = limit
		}
	}
}

func (o Orders) Len() int { return len(o) }
