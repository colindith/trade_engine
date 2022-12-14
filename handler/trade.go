package handler

import (
	"net/http"
	"strconv"
	"strings"
	"time"

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
		order.SetPrice(0)  // Use market price if input price is not specified. Else use limit price
	} else {
		price, err := strconv.Atoi(priceStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "invalid price",
			})
			return
		}
		order.SetPrice(price)
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
	order.SetQuantity(quantity)
	order.SetRemainingQuantity(quantity)

	action, _ := c.GetQuery("action")
	if strings.EqualFold(action, "buy") {
		order.SetAction(trade.ACTION_BUY)
	} else if strings.EqualFold(action, "sell") {
		order.SetAction(trade.ACTION_SELL)
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "invalid action",
		})
		return
	}

	success := trade.GetEngine().PlaceOrder(order)

	c.JSON(http.StatusOK, gin.H{
		"success": success,
		"order_id": order.OrderID(),
	})
}

type ListOrderResp struct {
	OrderID int           `json:"order_id,omitempty"`
	Timestamp time.Time   `json:"timestamp,omitempty"`
	Action string         `json:"action,omitempty"`
	Price int             `json:"price,omitempty"`
	Quantity int          `json:"quantity"`
	RemainingQuantity int `json:"remaining_quantity"`
	Status string         `json:"status,omitempty"`
}

func ListOrders(c *gin.Context) {
	orders := trade.GetEngine().ListOrders()
	res := make([]*ListOrderResp, 0, len(orders))
	for _, order := range orders {
		res = append(res, &ListOrderResp{
			OrderID:           order.OrderID(),
			Timestamp:         time.Unix(order.Timestamp() / 1e9, 0),
			Action:            order.GetActionRep(),
			Price:             order.Price(),
			Quantity:          order.Quantity(),
			RemainingQuantity: order.RemainingQuantity(),
			Status:            order.GetStatus(),
		})
	}
	c.JSON(http.StatusOK, gin.H{
		"orders": res,
	})
}