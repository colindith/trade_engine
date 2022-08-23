package trade

import "time"

type Action int

const (
	ACTION_BUY Action = iota
	ACTION_SELL
)

var orderCounter = 1

type Order struct {
	orderID int
	timestamp int64

	Action Action
	Price int
	Quantity int
}

func NewEmptyOrder() *Order {
	orderCounter += 1
	return &Order{
		orderID:   orderCounter,
		timestamp: time.Now().UnixNano(),
	}
}

func NewOrder(action Action, price int, quantity int) *Order {
	orderCounter += 1
	return &Order{
		orderID:   orderCounter,
		timestamp: time.Now().UnixNano(),
		Action:    action,
		Price:     price,
		Quantity:  quantity,
	}
}