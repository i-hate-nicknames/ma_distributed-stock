package order

import (
	"fmt"
	"sync"
)

// todo: add order persistence

const (
	// StatusProcessing means the order is being processed or is paused
	StatusProcessing = "processing"
	// StatusCanceled means the order is cancelled
	StatusCanceled = "canceled"
	// StatusCompleted means the order is completed, items are shipped
	StatusCompleted = "completed"
)

var orderCounter uint = 0

// Registry holds a registry of all orders in the system
type Registry struct {
	orders map[uint]*Order
	mux    sync.Mutex
}

// Order represents a user order for a set of items
type Order struct {
	ID         uint
	Items      []int64
	ReadyItems []int64
	Status     string
	mux        sync.Mutex
}

// MakeRegistry creates an empty order registry
func MakeRegistry() *Registry {
	orders := make(map[uint]*Order)
	return &Registry{orders: orders}
}

// MakeOrder creates new order with the given items, and assingns it an id
func MakeOrder(items []int64) *Order {
	orderCounter++
	return &Order{ID: orderCounter, Items: items, Status: StatusProcessing}
}

// GetOrder gets order by id if it's present in the system.
// The second result denotes whether the order was found
func (or *Registry) GetOrder(orderID uint) (*Order, bool) {
	or.mux.Lock()
	defer or.mux.Unlock()
	order, ok := or.orders[orderID]
	return order, ok
}

// SubmitOrder creates a new order and adds it to the registry
// returns an assigned order id for the new order or error
func (or *Registry) SubmitOrder(items []int64) (*Order, error) {
	or.mux.Lock()
	defer or.mux.Unlock()
	order := MakeOrder(items)
	or.orders[order.ID] = order
	return order, nil
}

func (or *Registry) CancelOrder(orderID uint) error {
	order, ok := or.GetOrder(orderID)
	if !ok {
		return fmt.Errorf("order %d not found", orderID)
	}
	order.mux.Lock()
	defer order.mux.Unlock()
	order.Status = StatusCanceled
	return nil
}

func (o *Order) AddReadyItems(taken []int64) {
	o.mux.Lock()
	defer o.mux.Unlock()
	o.Items = o.Items[len(taken):]
	o.ReadyItems = append(o.ReadyItems, taken...)
	if len(o.Items) == 0 {
		o.Status = StatusCompleted
	}
}
