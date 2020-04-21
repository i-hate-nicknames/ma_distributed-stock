package order

import (
	"context"
	"log"
	"sync"
	"time"

	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
)

var timeout = time.Second * 1

var mux sync.Mutex

// Processor represents an entity that loads orders in bulk
// and tries to ship needed items
type Processor struct {
	reg     *Registry
	catalog *warehouse.Catalog
}

// MakeProcessor returns a new processor with given order registry and warehouse catalog
func MakeProcessor(reg *Registry, catalog *warehouse.Catalog) *Processor {
	return &Processor{reg, catalog}
}

// RunProcessor to process orders that need processing
// orders with status pending will be checked if it's possible
// to ship required items. If it is, the shipment will be executed
// and the order updated with results of the shipment
// This is a blocking call
func (p *Processor) RunProcessor(ctx context.Context) {
	for {
		p.processOrders(ctx)
		time.Sleep(timeout)
	}
}

func (p *Processor) processOrders(ctx context.Context) {
	// prevent simulteneous processing runs
	mux.Lock()
	defer mux.Unlock()
	pending := p.reg.GetOrdersByStatus(StatusPending, -1)
	for _, order := range pending {
		p.process(ctx, order)
	}
}

func (p *Processor) process(ctx context.Context, order *Order) {
	log.Printf("Processing order %d, remaining items: %v", order.ID, order.Items)
	shipment, err := p.catalog.CalculateShipment(order.Items)
	if err != nil {
		log.Println("no viable warehouses to satisfy order, order:", order.ID)
		return
	}
	taken := p.executeShipment(ctx, shipment)
	if len(taken) == 0 {
		log.Println("no items retrieved, order:", order.ID)
		return
	}
	order.AddReadyItems(taken)
}

// executeShipment using given warehouse catalog.
// Request items from every warehouse in the shipment
// Return error if none of the items in the shipment were
// successfuly requested
// Otherwise return a slice of items that have been successfuly retrieved
func (p *Processor) executeShipment(ctx context.Context, shipment warehouse.Inventory) []int64 {
	taken := make([]int64, 0)
	var wg sync.WaitGroup
	wg.Add(len(shipment))
	for address, items := range shipment {
		go func(address string, items []int64) {
			defer wg.Done()
			log.Printf("Requesting %v from %s\n", items, address)
			remaining, err := warehouse.TakeItems(ctx, address, items)
			if err != nil {
				log.Println("Cannot take items from ", address, err.Error())
				return
			}
			taken = append(taken, items...)
			// update item list in the catalog
			p.catalog.SetWarehouse(address, remaining)
		}(address, items)
	}
	wg.Wait()
	return taken
}
