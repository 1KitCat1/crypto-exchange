package orderbook

import (
	"fmt"
	"reflect"
	"testing"
)

const (
	BID = true
	ASK = false
)

func assert(t *testing.T, a, b any) {
	if !reflect.DeepEqual(a, b) {
		t.Errorf("%+v != %+v", a, b)
	}
}

func TestLimit(t *testing.T) {
	limit := NewLimit(10_000)
	buyOrderA := NewOrder(true, 5) // true - buy|bid, false - sell|ask
	buyOrderB := NewOrder(true, 7)
	buyOrderC := NewOrder(true, 11)

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

	orderbook.PlaceLimitOrder(10_000, NewOrder(ASK, 10))
	assert(t, len(orderbook.asks), 1)
	orderbook.PlaceLimitOrder(11_000, NewOrder(ASK, 15))
	assert(t, len(orderbook.asks), 2)

	orderbook.PlaceLimitOrder(9_000, NewOrder(BID, 8))
	assert(t, len(orderbook.bids), 1)
}

func TestPlaceMarketOrder(t *testing.T) {
	orderbook := NewOrderbook()

	sellOrder := NewOrder(ASK, 10)
	orderbook.PlaceLimitOrder(100, sellOrder)

	buyOrder := NewOrder(BID, 5)
	matches := orderbook.PlaceMarketOrder(buyOrder)

	assert(t, len(matches), 1)
	assert(t, len(orderbook.asks), 1)
	assert(t, orderbook.AsksTotalVolume(), 5.0)
	assert(t, orderbook.BidTotalVolume(), 0.0)

	fmt.Println(matches)
}
