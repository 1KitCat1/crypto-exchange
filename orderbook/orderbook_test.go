package orderbook

import (
	"fmt"
	"testing"
)

func TestLimit(t *testing.T) {
	limit := NewLimit(10_000)
	buyOrder := NewOrder(true, 5) // true - buy|bid, false - sell|ask

	limit.AddOrder(buyOrder)
	fmt.Println(limit)
}

func TestOrderbook(t *testing.T) {

}
