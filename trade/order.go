package trade

import "time"

type Action int

const (
	ACTION_BUY Action = iota
	ACTION_SELL
)

var orderCounter = 1

type Order struct {
	OrderID   int
	Timestamp int64

	Action   Action
	Price    int
	Quantity int

	RemainingQuantity int
}

func NewEmptyOrder() *Order {
	orderCounter += 1
	return &Order{
		OrderID:   orderCounter,
		Timestamp: time.Now().UnixNano(),
	}
}

func NewOrder(action Action, price int, quantity int) *Order {
	orderCounter += 1
	return &Order{
		OrderID:           orderCounter,
		Timestamp:         time.Now().UnixNano(),
		Action:            action,
		Price:             price,
		Quantity:          quantity,
		RemainingQuantity: quantity,
	}
}