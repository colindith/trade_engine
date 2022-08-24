package trade

import (
	"sync"
	"time"
)

type Action int

const (
	ACTION_BUY Action = iota
	ACTION_SELL
)

var orderCounter = 1

type Order struct {
	mu sync.RWMutex

	orderID   int
	timestamp int64

	action   Action
	price    int
	quantity int

	remainingQuantity int
}

func (o *Order) OrderID() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.orderID
}

func (o *Order) SetOrderID(value int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.remainingQuantity = value
}

func (o *Order) Timestamp() int64 {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.timestamp
}

func (o *Order) SetTimestamp(timestamp int64) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.timestamp = timestamp
}

func (o *Order) Action() Action {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.action
}

func (o *Order) SetAction(action Action) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.action = action
}

func (o *Order) Price() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.price
}

func (o *Order) SetPrice(price int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.price = price
}

func (o *Order) Quantity() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.quantity
}

func (o *Order) SetQuantity(quantity int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.quantity = quantity
}

func (o *Order) RemainingQuantity() int {
	o.mu.RLock()
	defer o.mu.RUnlock()
	return o.remainingQuantity
}

func (o *Order) SetRemainingQuantity(value int) {
	o.mu.Lock()
	defer o.mu.Unlock()
	o.remainingQuantity = value
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
		orderID:           orderCounter,
		timestamp:         time.Now().UnixNano(),
		action:            action,
		price:             price,
		quantity:          quantity,
		remainingQuantity: quantity,
	}
}