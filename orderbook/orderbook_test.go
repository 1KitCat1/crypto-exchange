package orderbook

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	BID       = true
	ASK       = false
	MarketETH = "ETH"
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestLimit(t *testing.T) {
	limit := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5, MarketETH, 1111) // true - buy|bid, false - sell|ask
	buyOrderB := NewOrder(true, 7, MarketETH, 1111)
	buyOrderC := NewOrder(true, 11, MarketETH, 1111)

	limit.AddOrder(buyOrderA)
	limit.AddOrder(buyOrderB)
	limit.AddOrder(buyOrderC)

	if limit.Volume != 5+7+11 {
		t.Fail()
	}

	limit.DeleteOrder(buyOrderB)

	if limit.Volume != 5+11 {
		t.Fail()
	}
	fmt.Println(limit)
}

func TestPlaceLimitOrder(t *testing.T) {
	orderbook := NewOrderbook()

	orderbook.PlaceLimitOrder(10_000, NewOrder(ASK, 10, MarketETH, 1111))
	assert(t, len(orderbook.asks), 1)
	orderbook.PlaceLimitOrder(11_000, NewOrder(ASK, 15, MarketETH, 1111))
	assert(t, len(orderbook.asks), 2)

	orderbook.PlaceLimitOrder(9_000, NewOrder(BID, 8, MarketETH, 1111))
	assert(t, len(orderbook.bids), 1)
	assert(t, len(orderbook.Orders), 3)
}

func TestPlaceMarketOrder(t *testing.T) {
	orderbook := NewOrderbook()

	sellOrder := NewOrder(ASK, 10, MarketETH, 1111)
	orderbook.PlaceLimitOrder(100, sellOrder)

	buyOrder := NewOrder(BID, 5, MarketETH, 1111)
	matches := orderbook.PlaceMarketOrder(buyOrder)

	assert(t, len(matches), 1)
	assert(t, len(orderbook.asks), 1)
	assert(t, orderbook.AsksTotalVolume(), 5.0)
	assert(t, orderbook.BidsTotalVolume(), 0.0)
	assert(t, matches[0].Ask, sellOrder)
	assert(t, matches[0].Bid, buyOrder)
	assert(t, matches[0].SizeFilled, 5.0)
	assert(t, matches[0].Price, 100.0)
	assert(t, buyOrder.IsFilled(), true)

	fmt.Println(matches)
}

func TestPlaceMarketOrderMultipleMatches(t *testing.T) {
	orderbook := NewOrderbook()

	orderbook.PlaceLimitOrder(9_000, NewOrder(BID, 12, MarketETH, 1111))
	orderbook.PlaceLimitOrder(10_000, NewOrder(BID, 7, MarketETH, 1111))
	orderbook.PlaceLimitOrder(9_000, NewOrder(BID, 8, MarketETH, 1111))

	assert(t, orderbook.BidsTotalVolume(), 7+8+12.0)

	sellOrder := NewOrder(ASK, 20, MarketETH, 1111)
	matches := orderbook.PlaceMarketOrder(sellOrder)
	assert(t, len(matches), 3)
	assert(t, len(orderbook.asks), 0)
	assert(t, len(orderbook.bids), 1)
	assert(t, orderbook.BidsTotalVolume(), 7.0)

	fmt.Println(matches)

}

func TestPlaceMarketOrderMultipleOrdersOnLimit(t *testing.T) {
	orderbook := NewOrderbook()

	orderbook.PlaceLimitOrder(5_000, NewOrder(BID, 1, MarketETH, 1111))
	orderbook.PlaceLimitOrder(5_000, NewOrder(BID, 10, MarketETH, 1111))
	orderbook.PlaceLimitOrder(9_000, NewOrder(BID, 8, MarketETH, 1111))
	orderbook.PlaceLimitOrder(10_000, NewOrder(BID, 5, MarketETH, 1111))

	assert(t, orderbook.BidsTotalVolume(), 24.00)

	sellOrder := NewOrder(ASK, 20, MarketETH, 1111)
	matches := orderbook.PlaceMarketOrder(sellOrder)
	assert(t, len(matches), 4)
	assert(t, len(orderbook.asks), 0)
	assert(t, len(orderbook.bids), 1)
	assert(t, orderbook.BidsTotalVolume(), 4.0)

	fmt.Println(matches)
}

func TestPlaceMarketOrderOrderTiming(t *testing.T) {
	orderbook := NewOrderbook()

	orderbook.PlaceLimitOrder(5_000, NewOrder(BID, 10, MarketETH, 1111))
	orderbook.PlaceLimitOrder(5_000, NewOrder(BID, 1, MarketETH, 1111))
	orderbook.PlaceLimitOrder(9_000, NewOrder(BID, 8, MarketETH, 1111))
	orderbook.PlaceLimitOrder(10_000, NewOrder(BID, 5, MarketETH, 1111))

	assert(t, orderbook.BidsTotalVolume(), 24.00)

	sellOrder := NewOrder(ASK, 20, MarketETH, 1111)
	matches := orderbook.PlaceMarketOrder(sellOrder)
	assert(t, len(matches), 3)
	assert(t, len(orderbook.asks), 0)
	assert(t, len(orderbook.bids), 1)
	assert(t, orderbook.BidsTotalVolume(), 4.0)

	fmt.Println(matches)

}

func TestCancelOrder(t *testing.T) {
	orderbook := NewOrderbook()

	order := NewOrder(BID, 10, MarketETH, 1111)
	orderbook.PlaceLimitOrder(10_000, NewOrder(BID, 15, MarketETH, 1111))
	orderbook.PlaceLimitOrder(9_000, order)

	assert(t, len(orderbook.bids), 2)
	assert(t, orderbook.BidsTotalVolume(), 25.0)
	orderbook.CancelOrder(order)

	_, ok := orderbook.Orders[order.ID]
	assert(t, ok, false)

	assert(t, orderbook.BidsTotalVolume(), 15.0)
	assert(t, len(orderbook.bids), 1)
}
