package orderbook

import (
	"fmt"
	"sort"
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

func (orders Orders) Len() int           { return len(orders) }
func (orders Orders) Swap(i, j int)      { orders[i], orders[j] = orders[j], orders[i] }
func (orders Orders) Less(i, j int) bool { return orders[i].Timestamp < orders[j].Timestamp }

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

type Limits []*Limit

type ByBestAsk struct{ Limits }

func (ask ByBestAsk) Len() int           { return len(ask.Limits) }
func (ask ByBestAsk) Swap(i, j int)      { ask.Limits[i], ask.Limits[j] = ask.Limits[j], ask.Limits[i] }
func (ask ByBestAsk) Less(i, j int) bool { return ask.Limits[i].Price < ask.Limits[j].Price }

type ByBestBid struct{ Limits }

func (bid ByBestBid) Len() int           { return len(bid.Limits) }
func (bid ByBestBid) Swap(i, j int)      { bid.Limits[i], bid.Limits[j] = bid.Limits[j], bid.Limits[i] }
func (bid ByBestBid) Less(i, j int) bool { return bid.Limits[i].Price < bid.Limits[j].Price }

func NewLimit(price float64) *Limit {
	return &Limit{
		Price:  price,
		Orders: []*Order{},
	}
}

func (limit *Limit) String() string {
	return fmt.Sprintf("[price: %.2f | volume: %.2f]", limit.Price, limit.Volume)
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
	asks []*Limit
	bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		asks:      []*Limit{},
		bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
	}
}

func (orderbook *Orderbook) PlaceLimitOrder(price float64, order *Order) {
	var limit *Limit

	if order.Bid {
		limit = orderbook.BidLimits[price]
	} else {
		limit = orderbook.AskLimits[price]
	}

	if limit == nil {
		limit = NewLimit(price)
		if order.Bid {
			orderbook.bids = append(orderbook.bids, limit)
			orderbook.BidLimits[price] = limit
		} else {
			orderbook.asks = append(orderbook.asks, limit)
			orderbook.AskLimits[price] = limit
		}
	}

	limit.AddOrder(order)
}

func (orderbook *Orderbook) Asks() []*Limit {
	sort.Sort(ByBestAsk{orderbook.asks})
	return orderbook.asks
}

func (orderbook *Orderbook) Bids() []*Limit {
	sort.Sort(ByBestBid{orderbook.bids})
	return orderbook.bids
}
