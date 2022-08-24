package trade

import (
	"container/list"
	"log"
	"sync"

	"github.com/colindith/trade_engine/util"
)

type Engine struct {
	orderQueue chan *Order

	buyMap map[int]*list.List
	sellMap map[int]*list.List
	mapLock sync.RWMutex

	exit chan interface{}
}

var engineObj = &Engine{}

// GetEngine get the current Engine object
func GetEngine() *Engine {
	return engineObj
}

// PlaceOrder add a new order into the order queue
func (e *Engine) PlaceOrder(order *Order) bool {
	if e.orderQueue == nil {
		log.Print("[ERROR] trade engine not init")
		return false
	}
	e.orderQueue <- order
	return true
}

// ListOrders list all the ongoing orders
func (e *Engine) ListOrders() []*Order {
	e.mapLock.RLock()
	defer e.mapLock.RUnlock()
	res := make([]*Order, 0, len(e.sellMap) + len(e.buyMap))
	for _, orderList := range e.sellMap {
		for e := orderList.Front(); e != nil; e = e.Next() {
			if order, ok := e.Value.(*Order); ok {
				res = append(res, order)
			}
		}
	}

	for _, orderList := range e.buyMap {
		for e := orderList.Front(); e != nil; e = e.Next() {
			if order, ok := e.Value.(*Order); ok {
				res = append(res, order)
			}
		}
	}

	return res
}

// StartEngine starts the engine
func (e *Engine) StartEngine() {
	e.orderQueue = make(chan *Order, 100)

	e.buyMap = make(map[int]*list.List)
	e.sellMap = make(map[int]*list.List)

	e.exit = make(chan interface{})

	go func(){
		for {
			select {
			case <- e.exit:
				return
			case o := <- e.orderQueue:
				e.processOrder(o)
			}
		}
		log.Print("[INFO] Quit trade engine")
	}()
}

func (e *Engine) processOrder(o *Order) {
	e.mapLock.Lock()
	defer e.mapLock.Unlock()
	if o.Action == ACTION_BUY {
		processOrderInner(o, e.buyMap, e.sellMap)
	} else if o.Action == ACTION_SELL {
		// invert sellMap and buyMap
		processOrderInner(o, e.sellMap, e.buyMap)
	} else {
		log.Printf("[ERROR] Invalid action: %v", o.Action)
	}
}

func processOrderInner(o *Order, orderMap map[int]*list.List, counterOrderMap map[int]*list.List) {
	if orderList, ok := counterOrderMap[o.Price]; ok {
		for orderList.Len() != 0 && o.Quantity > 0 {
			if order, ok := orderList.Front().Value.(*Order); ok {
				successfulQuantity := util.Min(order.Quantity, o.Quantity)
				order.Quantity -= successfulQuantity
				o.Quantity -= successfulQuantity
				log.Printf("[DEBUG] order_id: %v successfully trade %v at price %v", order.OrderID, successfulQuantity, order.Price)
				log.Printf("[DEBUG] order_id: %v successfully trade %v at price %v", o.OrderID, successfulQuantity, o.Price)
				if order.Quantity == 0 {
					orderList.Remove(orderList.Front())
					log.Printf("[DEBUG] order_id: %v is done", order.OrderID)
				}
			}
		}
	}

	if o.Quantity == 0 {
		log.Printf("[DEBUG] order_id: %v is done", o.OrderID)
		return
	}
	// remaining quantity not finished
	// put it in buy map
	orderList, ok := orderMap[o.Price]
	if !ok {
		orderList = list.New()
		orderMap[o.Price] = orderList
	}
	orderList.PushBack(o)
}