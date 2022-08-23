package trade



type Action int

const (
	ACTION_BUY Action = iota
	ACTION_SELL
)

type Order struct {
	action Action
	price int64
	quantity int32
	limitPrice bool
}

//func NewOrder() {
//	return
//}