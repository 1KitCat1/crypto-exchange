package main

import (
	"crypto-exchange/orderbook"
	"encoding/json"
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	exchange := NewExchange()

	e.GET("/book/", exchange.handleGetBook)
	e.POST("/order", exchange.handlePlaceOrder)

	e.Start(":3000")

	var _ = orderbook.Limit{}
	fmt.Println("Check")
}

type Market string

const (
	MarketETH Market = "ETH"
)

type Exchange struct {
	orderbooks map[Market]*orderbook.Orderbook
}

func NewExchange() *Exchange {
	return &Exchange{
		orderbooks: make(map[Market]*orderbook.Orderbook),
	}
}

type OrderType bool

const (
	MarketOrder OrderType = true
	LimitOrder  OrderType = false
)

type PlaceOrderRequest struct {
	Type   OrderType // limit / market
	Bid    bool      // bid / ask
	Size   float64
	Price  float64
	Market Market
}

func (exchange *Exchange) handlePlaceOrder(context echo.Context) error {
	var placeOrderRequest PlaceOrderRequest

	if err := json.NewDecoder(context.Request().Body).Decode(&placeOrderRequest); err != nil {
		return err
	}

	market := Market(placeOrderRequest.Market)
	ob := exchange.orderbooks[market]
	order := orderbook.NewOrder(placeOrderRequest.Bid, placeOrderRequest.Size)
	ob.PlaceLimitOrder(order.Limit.Price, order)

	return context.JSON(200, map[string]any{"msg": "Order placed"})
}

func (exchange *Exchange) handleGetBook(context echo.Context) error {

}
