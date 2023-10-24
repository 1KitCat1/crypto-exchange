package orderbook

import (
	"fmt"
	"testing"
)

const (
	BID = true
	ASK = false
)

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

func TestOrderbook(t *testing.T) {
	orderbook := NewOrderbook()

	orderbook.PlaceOrder(34_000, NewOrder(BID, 10))
	orderbook.PlaceOrder(34_000, NewOrder(BID, 100))

	if orderbook.BidLimits[34_000].Volume != 100+10 {
		t.Fail()
	}

	fmt.Println(orderbook.Bids)
}
