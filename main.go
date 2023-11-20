package main

import (
	"crypto-exchange/orderbook"
	"fmt"

	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()
	e.HTTPErrorHandler = httpErrorHandler
	exchange := NewExchange()

	e.GET("/book/:market", exchange.handleGetBook)
	e.GET("/book/volume/:market", exchange.handleGetVolume)
	e.POST("/order", exchange.handlePlaceOrder)
	e.DELETE("/order/:id", exchange.handleCancelOrder)
	e.Start(":3000")

	var _ = orderbook.Limit{}
	fmt.Println("Check")
}
