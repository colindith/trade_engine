package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"github.com/colindith/trade_engine/handler"
	"github.com/colindith/trade_engine/trade"
)

func main() {
	trade.GetEngine().StartEngine()
	defer trade.GetEngine().Stop()

	r := gin.Default()
	r.GET("/api/ping", handler.Ping)
	r.GET("/api/place_order", handler.PlaceOrder)
	r.GET("/api/list_orders", handler.ListOrders)

	err := r.Run() // listen and serve on 0.0.0.0:8080
	if err != nil {
		log.Print("[ERROR] server ends with error: " + err.Error())
	}
}

