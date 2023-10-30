package main

import (
	"crypto-exchange/orderbook"
	"fmt"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	e.Start(":3000")

	var _ = orderbook.Limit{}
	fmt.Println("Check")
}
