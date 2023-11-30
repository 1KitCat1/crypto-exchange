package main

import (
	"crypto-exchange/orderbook"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

func httpErrorHandler(err error, context echo.Context) {
	fmt.Println(err)
}

type OrderType string

type PlaceOrderRequest struct {
	Type   OrderType // limit / market
	Bid    bool      // bid / ask
	Size   float64
	Price  float64
	UserID int64
	Market orderbook.Market
}

func (exchange *Exchange) handlePlaceOrder(context echo.Context) error {
	var placeOrderRequest PlaceOrderRequest

	if err := json.NewDecoder(context.Request().Body).Decode(&placeOrderRequest); err != nil {
		return err
	}

	market := orderbook.Market(placeOrderRequest.Market)
	ob := exchange.orderbooks[market]
	order := orderbook.NewOrder(placeOrderRequest.Bid, placeOrderRequest.Size, placeOrderRequest.Market, placeOrderRequest.UserID)
	exchange.orders[order.ID] = order

	if _, exists := exchange.userOrders[order.UserID]; !exists {
		exchange.userOrders[order.UserID] = make(map[int64]struct{})
	}
	exchange.userOrders[order.UserID][order.ID] = struct{}{}

	if placeOrderRequest.Type == LimitOrder {
		ob.PlaceLimitOrder(placeOrderRequest.Price, order)
		return context.JSON(200, map[string]any{"msg": "Limit order placed"})
	} else {
		matches := ob.PlaceMarketOrder(order)

		if matches == nil {
			return context.JSON(http.StatusConflict, map[string]any{
				"msg":    "Orderbook cannot fullfill requested marked order",
				"reason": "Not enough volume",
			})
		}

		// This is done to avoid recursion
		matchedOrders := make([]*MatchView, len(matches))

		for i := 0; i < len(matchedOrders); i++ {
			matchedOrders[i] = getMatchView(&matches[i])
		}
		return context.JSON(200, map[string]any{
			"msg":     "Market order has been executed",
			"matches": matchedOrders,
		})
	}
}

func (exchange *Exchange) handleGetOrder(context echo.Context) error {
	orderIdStr := context.Param("orderId")
	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)

	if err != nil {
		panic("Unable to parse order identfier") // TODO: replace with invalid request code
	}

	order, ok := exchange.orders[orderId]

	if ok {
		orderView := getOrderView(order)
		return context.JSON(http.StatusOK, orderView)
	} else {
		return context.JSON(http.StatusBadRequest, map[string]any{"msg": "Order is not in the exchange memory"})
	}
}

func (exchange *Exchange) handleGetUserOrders(context echo.Context) error {
	userIdStr := context.Param("userId")
	userId, err := strconv.ParseInt(userIdStr, 10, 64)

	if err != nil {
		return context.JSON(http.StatusBadRequest, map[string]any{"msg": "Unable to parse userID"})
	}

	orderIds := exchange.userOrders[userId]

	return context.JSON(http.StatusOK, orderIds)
}

func (exchange *Exchange) handleCancelOrder(context echo.Context) error {
	orderIdStr := context.Param("id")
	orderId, err := strconv.ParseInt(orderIdStr, 10, 64)
	if err != nil {
		panic("Unable to parse order identifier")
	}
	orderbook := exchange.orderbooks[MarketETH]
	order := orderbook.Orders[orderId]
	orderbook.CancelOrder(order)
	return context.JSON(200, map[string]any{
		"msg": "Order has been canceled",
	})
}

func (exchange *Exchange) handleGetBook(context echo.Context) error {
	market := orderbook.Market(context.Param("market"))
	ob, ok := exchange.orderbooks[market]

	if !ok {
		return context.JSON(http.StatusBadRequest, map[string]any{"msg": "market not found"})
	}

	orderbookData := OrderbookData{
		BidsTotalVolume: ob.BidsTotalVolume(),
		AsksTotalVolume: ob.AsksTotalVolume(),
		Asks:            []*OrderView{},
		Bids:            []*OrderView{},
	}

	for _, limit := range ob.Asks() {
		for _, order := range limit.Orders {
			orderbookData.Asks = append(orderbookData.Asks, getOrderView(order))
		}
	}

	for _, limit := range ob.Bids() {
		for _, order := range limit.Orders {
			orderbookData.Bids = append(orderbookData.Bids, getOrderView(order))
		}
	}

	return context.JSON(http.StatusOK, orderbookData)
}

func (exchange *Exchange) handleGetVolume(context echo.Context) error {
	market := orderbook.Market(context.Param("market"))
	ob, ok := exchange.orderbooks[market]

	if !ok {
		return context.JSON(http.StatusBadRequest, map[string]any{
			"msg": "Unable to find market with the provided identifier",
		})
	}

	return context.JSON(http.StatusOK, map[string]any{
		"Bid": ob.BidsTotalVolume(),
		"Ask": ob.AsksTotalVolume(),
	})
}
