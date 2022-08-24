# trade_engine

### Usage
Start the server:
```bash
> go run main/main.go
```
Place a new order:
```bash
> curl http://localhost:8080/api/place_order?price=13&quantity=13&action=buy
{"order_id": 1,"success": true}
> curl http://localhost:8080/api/place_order?price=13&quantity=10&action=sell
{"order_id": 2,"success": true}
```
List all the orders:
```bash
> curl http://localhost:8080/api/list_orders
{
  "orders": [
    {
      "order_id": 1,
      "timestamp": "2022-08-24T19:28:18+08:00",
      "action": "buy",
      "price": 13,
      "quantity": 13,
      "remaining_quantity": 3,
      "status": "ongoing"
    },
    {
      "order_id": 2,
      "timestamp": "2022-08-24T19:28:20+08:00",
      "action": "sell",
      "price": 13,
      "quantity": 10,
      "remaining_quantity": 0,
      "status": "completed"
    }
  ]
}
```
The 10 out of the 13 in the first order were done with the second one. 3 didn't match any other sell order. So the status of the first order is still `ongoing`.
The second order finished all the transaction. The remaining quantity is 0 thus the status is `completed`.

### API Design
`/api/place_order`

The api for a user to create a new buy/sell order with price, quantity.

 - Request Payload:

|Parameter|type|Description|
|----|----|----|
|Price|int|The limit price of the order. If not specified, will use the market price.|
|quantity|int|The quantity of the order. The Engine will try to find other orders that match the price. The transaction can happen between multiple orders if one order can not fulfill the quantity.|
|action|string|Action can be either `buy` or `sell`|

 - Response

|Parameter|type|Description|
|----|----|----|
|success|bool|True if the action success. This field only indicates that the order is recorded by the system. It doesn't guarantee the transaction to succeed. The transaction will happen asynchronously.|
|order_id|int|The order_id generated by the system.|
|error|string|The error message returned by the system, if any.|

`/api/list_orders`

This api list the details of all the orders recorded by the system.

- Request Payload:

N/A

- Response:

|Parameter|type|Description|
|----|----|----|
|orders|Order list||

 - Order:

|Parameter|type|Description|
|----|----|----|
|order_id|int|The id of the order.|
|timestamp|string|The timestamp when the order be placed.|
|action|string|Can be either `buy` or `sell`.|
|price|int|The limit price of the order.|
|quantity|int|The user specified quantity of the order.|
|remaining_quantity|int|The amount that still pending trading. 0 means order is completed.|
|status|string|The current status of the order. It can be either `ongoing` or `completed`|

### Engine Design
In the Trade Engine, two kinds of hash map are used to store the order data.
They are `sellMap` and `buyMap`.
The `sellMap` only store the selling orders while the `buyMap` only store the buying orders.
The key of the maps will be the price of the orders.
So that we can find all the orders of some specific price.
The structure of the map are shown below:
```javascript
sellMap
{
    100: [selling_order_1, selling_order_2], 
    101: [selling_order_3],
}
// The above sellMap means we have 2 selling order at price 100 and 1 selling order at price 101
buyMap
{
    99: [buying_order_1, buying_order_2],
}
// The above buyMap means we have 2 buying order at price 99
```

Based on these two maps, once there is a new buying order coming into the Engine, it will check if there is any selling order who has the same price in the `sellMap`.
If it is, the Engine will let the transaction happen between these two orders.

If the quantity of these two orders are different, let's say, the quantity of the buying order is more than the quantity of the selling order,
then the remaining buying order will be written back the buyMap and waiting for next selling order that has the same price. 

Once an order is completed, it will be popped from the front of the list in the map. In order to do so, we use the linked list to store the orders.

The map in Go is not thread-safe, so whenever we try access the `sellMap` or `buyMap`, it is necessary to acquire the mutex lock first.
```go
type Engine struct {
	mapLock sync.RWMutex
	...
}
```
Similarly, the `Order` object in the map can be accessed by multiple goroutines concurrently, so it also has a mutex lock in each of the individual object.
```go
type Order struct {
    mu sync.RWMutex
    ...
}
```