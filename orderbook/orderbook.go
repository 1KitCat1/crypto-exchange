package orderbook

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"time"
)

type Market string

type Match struct {
	Ask        *Order
	Bid        *Order
	SizeFilled float64
	Price      float64
}

type Order struct {
	ID        int64
	UserID    int64
	Market    Market
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

func NewOrder(is_bid bool, size float64, market Market, userID int64) *Order {
	return &Order{
		ID:        rand.Int63n(int64(1) << 62),
		UserID:    userID,
		Size:      size,
		Market:    market,
		Bid:       is_bid,
		Timestamp: time.Now().UnixNano(),
	}
}

func (order *Order) IsFilled() bool {
	return order.Size == 0.0
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
func (bid ByBestBid) Less(i, j int) bool { return bid.Limits[i].Price > bid.Limits[j].Price }

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

func (limit *Limit) Fill(order *Order) []Match {
	matches := []Match{}
	filledOrders := []*Order{}

	for _, orderMatched := range limit.Orders {
		match := limit.fillOrder(order, orderMatched)
		matches = append(matches, match)

		limit.Volume -= match.SizeFilled

		if orderMatched.IsFilled() {
			filledOrders = append(filledOrders, orderMatched)
		}

		if order.IsFilled() {
			break
		}
	}

	for _, order := range filledOrders {
		limit.DeleteOrder(order)
	}
	return matches
}

func (limit *Limit) fillOrder(firstOrder, secondOrder *Order) Match {
	bid := firstOrder
	ask := secondOrder

	if secondOrder.Bid {
		bid = secondOrder
		ask = firstOrder
	}

	filled := math.Min(bid.Size, ask.Size)

	bid.Size -= filled
	ask.Size -= filled

	return Match{
		Bid:        bid,
		Ask:        ask,
		SizeFilled: filled,
		Price:      limit.Price,
	}
}

type Orderbook struct {
	asks []*Limit
	bids []*Limit

	AskLimits map[float64]*Limit
	BidLimits map[float64]*Limit
	Orders    map[int64]*Order
}

func NewOrderbook() *Orderbook {
	return &Orderbook{
		asks:      []*Limit{},
		bids:      []*Limit{},
		AskLimits: make(map[float64]*Limit),
		BidLimits: make(map[float64]*Limit),
		Orders:    make(map[int64]*Order),
	}
}

func (orderbook *Orderbook) PlaceMarketOrder(order *Order) []Match {
	matches := []Match{}

	if order.Bid {
		if order.Size > orderbook.AsksTotalVolume() {
			return nil
		}
		for _, limit := range orderbook.Asks() {
			limitMatches := limit.Fill(order)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				orderbook.clearLimit(false, limit)
			}
		}
	} else {
		if order.Size > orderbook.BidsTotalVolume() {
			return nil
		}
		for _, limit := range orderbook.Bids() {
			limitMatches := limit.Fill(order)
			matches = append(matches, limitMatches...)

			if len(limit.Orders) == 0 {
				orderbook.clearLimit(true, limit)
			}
		}
	}

	return matches
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
	orderbook.Orders[order.ID] = order
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

func (orderbook *Orderbook) clearLimit(is_bid bool, limit *Limit) {
	if is_bid {
		delete(orderbook.BidLimits, limit.Price)
		for i := 0; i < len(orderbook.bids); i++ {
			if orderbook.bids[i] == limit {
				orderbook.bids[i] = orderbook.bids[len(orderbook.bids)-1]
				orderbook.bids = orderbook.bids[:len(orderbook.bids)-1]
			}
		}
	} else {
		delete(orderbook.AskLimits, limit.Price)
		for i := 0; i < len(orderbook.asks); i++ {
			if orderbook.asks[i] == limit {
				orderbook.asks[i] = orderbook.asks[len(orderbook.asks)-1]
				orderbook.asks = orderbook.asks[:len(orderbook.asks)-1]
			}
		}
	}
}

func (orderbook *Orderbook) CancelOrder(order *Order) {
	limit := order.Limit
	limit.DeleteOrder(order)
	if len(limit.Orders) == 0 {
		orderbook.clearLimit(order.Bid, limit)
	}
	delete(orderbook.Orders, order.ID)
}

func (orderbook *Orderbook) BidsTotalVolume() float64 {
	totalVolume := 0.0
	// TODO: optimize
	for i := 0; i < len(orderbook.bids); i++ {
		totalVolume += orderbook.bids[i].Volume
	}

	return totalVolume
}

func (orderbook *Orderbook) AsksTotalVolume() float64 {
	totalVolume := 0.0
	// TODO: optimize
	for i := 0; i < len(orderbook.asks); i++ {
		totalVolume += orderbook.asks[i].Volume
	}

	return totalVolume
}
