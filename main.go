package main

import (
	"crypto-exchange/orderbook"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	exchange := NewExchange()

	e.GET("/book/:market", exchange.handleGetBook)
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
	orderbooks := make(map[Market]*orderbook.Orderbook)
	orderbooks[MarketETH] = orderbook.NewOrderbook()

	return &Exchange{
		orderbooks: orderbooks,
	}
}

type OrderType string

const (
	MarketOrder OrderType = "MARKET"
	LimitOrder  OrderType = "LIMIT"
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

	if placeOrderRequest.Type == LimitOrder {
		ob.PlaceLimitOrder(placeOrderRequest.Price, order)
		return context.JSON(200, map[string]any{"msg": "Limit order placed"})
	} else {
		matches := ob.PlaceMarketOrder(order)
		return context.JSON(200, map[string]any{
			"msg":     "Market order has been executed",
			"matches": len(matches),
		})
	}

}

type Order struct {
	Price     float64
	Size      float64
	Bid       bool
	Timestamp int64
}

type OrderbookData struct {
	BidsTotalVolume float64
	AsksTotalVolume float64
	Asks            []*Order
	Bids            []*Order
}

func (exchange *Exchange) handleGetBook(context echo.Context) error {
	market := Market(context.Param("market"))
	ob, ok := exchange.orderbooks[market]

	if !ok {
		return context.JSON(http.StatusBadRequest, map[string]any{"msg": "market not found"})
	}

	orderbookData := OrderbookData{
		BidsTotalVolume: ob.BidsTotalVolume(),
		AsksTotalVolume: ob.AsksTotalVolume(),
		Asks:            []*Order{},
		Bids:            []*Order{},
	}

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			order := Order{
				Price:     order.Limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookData.Asks = append(orderbookData.Asks, &order)
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			order := Order{
				Price:     order.Limit.Price,
				Size:      order.Size,
				Bid:       order.Bid,
				Timestamp: order.Timestamp,
			}
			orderbookData.Bids = append(orderbookData.Bids, &order)
		}
	}

	return context.JSON(http.StatusOK, orderbookData)
}
