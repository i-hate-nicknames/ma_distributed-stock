package order

import "sync"

const (
	// StatusPending means the order is being processed or is paused
	StatusPending = "pending"
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
	ID     uint
	Items  []int64
	Status string
	mux    sync.Mutex
}

// MakeOrder creates new order with the given items, and assingns it an id
func MakeOrder(items []int64) *Order {
	orderCounter++
	return &Order{ID: orderCounter, Items: items, Status: StatusPending}
}

// GetOrder gets order by id if it's present in the system.
// The second result denotes whether the order was found
func (or *Registry) GetOrder(orderID uint) (*Order, bool) {
	or.mux.Lock()
	defer or.mux.Unlock()
	val, ok := or.orders[orderID]
	return val, ok
}

// SubmitOrder creates a new order and adds it to the registry
// returns an assigned order id for the new order
func (or *Registry) SubmitOrder(items []int64) uint {
	or.mux.Lock()
	defer or.mux.Unlock()
	order := MakeOrder(items)
	or.orders[order.ID] = order
	return order.ID
}
