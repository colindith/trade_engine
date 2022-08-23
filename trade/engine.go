//go:generate mockgen -source=engine.go
package trade

// Engine define the interface of the Engine
type Engine interface {
	PlaceOrder(order Order)
	ListOrders() []*Order
	StartEngine()
	StopEngine()
}

type defaultEngineImp struct {}

var engineObj Engine = &defaultEngineImp{}

// GetEngine get the current Engine object
func GetEngine() Engine {
	return engineObj
}

// SetEngine set the current Engine object
func SetEngine(t Engine) {
	engineObj = t
}

// PlaceOrder see the Engine PlaceOrder method
func (t *defaultEngineImp) PlaceOrder(order Order) {

}

// ListOrders see the Engine ListOrders method
func (t *defaultEngineImp) ListOrders() []*Order {
	return []*Order{}
}

// StartEngine see the Engine StartEngine method
func (t *defaultEngineImp) StartEngine() {

}

// StopEngine see the Engine StopEngine method
func (t *defaultEngineImp) StopEngine() {

}