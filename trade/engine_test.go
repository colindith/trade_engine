package trade

import (
	"container/list"
	"reflect"
	"testing"
)

func TestEngine_processOrderInner(t *testing.T) {
	type args struct {
		o               *Order
		orderMap        map[int]*list.List
		counterOrderMap map[int]*list.List
	}
	tests := []struct {
		name                    string
		args                    args
		expectedOrderMap        map[int]*list.List
		expectedCounterOrderMap map[int]*list.List
	}{
		{
			name:                    "new_order_completed_immediately",
			args:                    args{
				o:               &Order{
					orderID:           4,
					timestamp:         1661319309,
					action:            0,
					price:             11,
					quantity:          100,
					remainingQuantity: 100,
				},
				orderMap:        map[int]*list.List{},
				counterOrderMap: map[int]*list.List{
					11: orderSliceToList([]*Order{
						{
							orderID:           1,
							timestamp:         1661319306,
							action:            1,
							price:             11,
							quantity:          50,
							remainingQuantity: 40,
						},
						{
							orderID:           2,
							timestamp:         1661319307,
							action:            1,
							price:             11,
							quantity:          40,
							remainingQuantity: 40,
						},
						{
							orderID:           3,
							timestamp:         1661319308,
							action:            1,
							price:             11,
							quantity:          40,
							remainingQuantity: 40,
						},
					}),
				},
			},
			expectedOrderMap:        map[int]*list.List{},
			expectedCounterOrderMap: map[int]*list.List{
				11: orderSliceToList([]*Order{
					{
						orderID:           3,
						timestamp:         1661319308,
						action:            1,
						price:             11,
						quantity:          40,
						remainingQuantity: 20,
					},
				}),
			},
		},
		{
			name:                    "new_order_not_completed",
			args:                    args{
				o:               &Order{
					orderID:           4,
					timestamp:         1661319309,
					action:            1,
					price:             11,
					quantity:          150,
					remainingQuantity: 150,
				},
				orderMap:        map[int]*list.List{},
				counterOrderMap: map[int]*list.List{
					11: orderSliceToList([]*Order{
						{
							orderID:           1,
							timestamp:         1661319306,
							action:            0,
							price:             11,
							quantity:          50,
							remainingQuantity: 40,
						},
						{
							orderID:           2,
							timestamp:         1661319307,
							action:            0,
							price:             11,
							quantity:          40,
							remainingQuantity: 40,
						},
						{
							orderID:           3,
							timestamp:         1661319308,
							action:            0,
							price:             11,
							quantity:          40,
							remainingQuantity: 40,
						},
					}),
				},
			},
			expectedOrderMap:        map[int]*list.List{
				11: orderSliceToList([]*Order{
					{
						orderID:           4,
						timestamp:         1661319309,
						action:            1,
						price:             11,
						quantity:          150,
						remainingQuantity: 30,
					},
				}),
			},
			expectedCounterOrderMap: map[int]*list.List{
				11: orderSliceToList([]*Order{}),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Engine{}
			e.processOrderInner(tt.args.o, tt.args.orderMap, tt.args.counterOrderMap)
			if !reflect.DeepEqual(tt.args.orderMap, tt.expectedOrderMap) {
				t.Errorf("orderMap not currect. got = %v, want %v", tt.args.orderMap, tt.expectedOrderMap)
			}
			if !reflect.DeepEqual(tt.args.counterOrderMap, tt.expectedCounterOrderMap) {
				t.Errorf("counterOrderMap not currect. got = %v, want %v", tt.args.counterOrderMap, tt.expectedCounterOrderMap)
			}
		})
	}
}

func orderSliceToList(orders []*Order) *list.List {
	l := list.New()
	for i := range orders {
		l.PushBack(orders[i])
	}
	return l
}