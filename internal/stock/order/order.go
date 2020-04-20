package order

import (
	"errors"
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
// Items represent items that need to be fetched from remote warehouses
// ReadyItems represent items that have been fetched from remote warehouses
// and are ready to be delivered to the client
// Orders is created with empty ReadyItems and when all the Items move to
// ReadyItems the order is considered to be completed
type Order struct {
	ID         uint
	Items      []int64
	ReadyItems []int64
	Status     string
	mux        sync.Mutex
}

var NotFoundError = errors.New("order not found")

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
func (or *Registry) GetOrder(orderID uint) (*Order, error) {
	or.mux.Lock()
	defer or.mux.Unlock()
	order, ok := or.orders[orderID]
	if !ok {
		return nil, NotFoundError
	}
	return order, nil
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

// CancelOrder identified by the given orderID
// This function only updates the status and doesn't do
// anything with the ready items
// Return NotFoundError when there is no such order in the system
func (or *Registry) CancelOrder(orderID uint) error {
	order, err := or.GetOrder(orderID)
	if err != nil {
		return err
	}
	order.mux.Lock()
	defer order.mux.Unlock()
	// todo: currently sets order state to canceled
	// when order scheduler is implemented the status should
	// be pendingCancel that denote that the order is planned to
	// be canceled
	order.Status = StatusCanceled
	return nil
}

// AddReadyItems adds items to the list of ready for delivery
// items of this order
func (o *Order) AddReadyItems(items []int64) {
	o.mux.Lock()
	defer o.mux.Unlock()
	o.Items = o.Items[len(items):]
	o.ReadyItems = append(o.ReadyItems, items...)
	if len(o.Items) == 0 {
		o.Status = StatusCompleted
	}
}
