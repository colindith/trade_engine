package trade

import (
	"container/list"
	"log"
	"sync"

	"github.com/colindith/trade_engine/util"
)

const MAX_INT = int(^uint(0) >> 1)

type Engine struct {
	orderQueue chan *Order

	buyMap map[int]*list.List   // all the ongoing "BUY" orders. The key represent the price of the orders
	sellMap map[int]*list.List  // all the ongoing "SELL" orders. The key represent the price of the orders
	mapLock sync.RWMutex

	orders []*Order  // records all the ongoing/completed orders. Only used for displaying all orders
	ordersLock sync.RWMutex

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

	e.ordersLock.Lock()
	defer e.ordersLock.Unlock()
	e.orders = append(e.orders, order)

	e.orderQueue <- order
	return true
}

// ListOrders list all the ongoing orders
func (e *Engine) ListOrders() []*Order {
	e.ordersLock.RLock()
	defer e.ordersLock.RUnlock()
	res := make([]*Order, 0, len(e.orders))
	for _, order := range e.orders {
		res = append(res, order)
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
				log.Print("[INFO] Quit trade engine")
				return
			case o := <- e.orderQueue:
				e.processOrder(o)
			}
		}
	}()
}

func (e *Engine) processOrder(o *Order) {
	e.mapLock.Lock()
	defer e.mapLock.Unlock()
	if o.Action() == ACTION_BUY {
		e.processOrderInner(o, e.buyMap, e.sellMap)
	} else if o.Action() == ACTION_SELL {
		// invert sellMap and buyMap
		e.processOrderInner(o, e.sellMap, e.buyMap)
	} else {
		log.Printf("[ERROR] Invalid action: %v", o.Action())
	}
}

func (e *Engine) processOrderInner(o *Order, orderMap map[int]*list.List, counterOrderMap map[int]*list.List) {
	if o.Price() == 0 {
		// set market price
		o.SetPrice(e.getMarketPrice(o.action))
	}
	counterOrderList, ok := counterOrderMap[o.Price()]
	if ok {
		for counterOrderList.Len() != 0 && o.RemainingQuantity() > 0 {
			order, ok := counterOrderList.Front().Value.(*Order)
			if !ok {
				log.Print("[ERROR] Order assertion error")
				continue
			}

			func() {
				order.mu.Lock()
				defer order.mu.Unlock()
				o.mu.Lock()
				defer o.mu.Unlock()

				successfulQuantity := util.Min(order.remainingQuantity, o.remainingQuantity)
				order.remainingQuantity -= successfulQuantity
				o.remainingQuantity -= successfulQuantity
				log.Printf("[DEBUG] order_id: %v successfully %v %v at price %v", order.orderID, order.GetActionRep(), successfulQuantity, order.price)
				log.Printf("[DEBUG] order_id: %v successfully %v %v at price %v", o.orderID, o.GetActionRep(), successfulQuantity, o.price)
				if order.remainingQuantity == 0 {
					counterOrderList.Remove(counterOrderList.Front())
					log.Printf("[DEBUG] order_id: %v is done", order.orderID)
				}
			}()
		}
	}

	if o.RemainingQuantity() == 0 {
		log.Printf("[DEBUG] order_id: %v is done", o.OrderID())
		return
	}
	// remaining quantity not finished
	// put it in the map
	orderList, ok := orderMap[o.Price()]
	if !ok {
		orderList = list.New()
		orderMap[o.Price()] = orderList
	}
	orderList.PushBack(o)
}

func (e *Engine) getMarketPrice(action Action) int {
	if action == ACTION_BUY {
		// try get max from sell_map
		maxPrice := -MAX_INT-1
		for price, orderList := range e.sellMap {
			if price > maxPrice && orderList.Len() != 0 {
				maxPrice = price
			}
		}
		if maxPrice != -MAX_INT-1 {
			return maxPrice
		}
	} else {
		// try get min from buy_map
		minPrice := MAX_INT
		for price, orderList := range e.buyMap {
			if price < minPrice && orderList.Len() != 0 {
				minPrice = price
			}
		}
		if minPrice != MAX_INT {
			return minPrice
		}
	}
	return 0
}

func (e *Engine) Stop() {
	e.exit <- true
}