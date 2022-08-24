package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/colindith/trade_engine/trade"
)

func Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func PlaceOrder(c *gin.Context) {
	order := trade.NewEmptyOrder()

	// parsing price query string
	priceStr, ok := c.GetQuery("price")
	if len(priceStr) == 0 || !ok {
		order.Price = 0
	} else {
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid price",
			})
			return
		}
		order.Price = price
	}

	// parsing quantity query string
	quantityStr, ok := c.GetQuery("quantity")
	if len(quantityStr) == 0 || !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid empty quantity",
		})
		return
	}
	quantity, err := strconv.Atoi(quantityStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid quantity " + quantityStr,
		})
		return
	}
	order.Quantity = quantity

	action, _ := c.GetQuery("action")
	if strings.EqualFold(action, "buy") {
		order.Action = trade.ACTION_BUY
	} else if strings.EqualFold(action, "sell") {
		order.Action = trade.ACTION_SELL
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid action",
		})
		return
	}

	success := trade.GetEngine().PlaceOrder(order)

	c.JSON(http.StatusOK, gin.H{
		"success": success,
	})
}

func ListOrders(c *gin.Context) {
	orders := trade.GetEngine().ListOrders()
	// TODO: format the response
	c.JSON(http.StatusOK, gin.H{
		"orders": orders,
	})
}