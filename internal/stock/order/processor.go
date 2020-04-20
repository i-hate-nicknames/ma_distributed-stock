package order

import (
	"context"
	"log"
	"sync"
	"time"
)

var timeout = time.Second * 1

var mux sync.Mutex

func RunProcessor(ctx context.Context, reg *Registry) {
	for {
		processOrders(reg)
		time.Sleep(timeout)
	}
}

func processOrders(reg *Registry) {
	// prevent simulteneous processing runs
	mux.Lock()
	defer mux.Unlock()
	pending := reg.GetOrdersByStatus(StatusPending, -1)
	for _, order := range pending {
		process(order)
	}
}

func process(order *Order) {
	log.Printf("Processing order %d, remaining items: %v", order.ID, order.Items)
	// simulate work for now
	time.Sleep(2 * time.Second)
}
