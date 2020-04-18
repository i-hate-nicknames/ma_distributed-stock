package stock

import (
	"context"
	"log"

	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/order"
	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
)

// Stock represents a central facade with which clients interact via
// order submission, and which contacts remote warehouses for items
type Stock struct {
	Catalog *warehouse.Catalog
	Orders  *order.Registry
}

// DiscoverWarehouses listens to warehouse invitations and adds
// newly discovered warehouses to the stock catalog
func (s *Stock) DiscoverWarehouses(ctx context.Context) {
	addresses := make(chan string, 5)
	go warehouse.ListenToInvitations(ctx, addresses)
	for address := range addresses {
		if s.Catalog.HasWarehouse(address) {
			continue
		}
		s.addWarehouse(ctx, address)
	}
}

// todo maybe add method to check which warehouses are still
// alive
// maybe make an infinite loop that will periodically check on every warehouse
// and remove those that are dead

// query warehouse located at the given address for its items
// and add it to the catalog,
func (s *Stock) addWarehouse(ctx context.Context, address string) {
	items, err := warehouse.LoadItems(ctx, address)
	if err != nil {
		s.Catalog.RemoveWarehouse(address)
		return
	}
	log.Printf("Added warehouse %s with items %v\n", address, items)
	s.Catalog.AddWarehouse(address, items)
}

// SubmitOrder creates a new order for given items and immediately
// tries to ship items
func (s *Stock) SubmitOrder(items []int64) (*order.Order, error) {
	// todo: later order processing will happen asynchronously
	// in a different thread, and order, submission will just create an order
	ord, err := s.Orders.SubmitOrder(items)
	if err != nil {
		return nil, err
	}
	shipment, err := s.Catalog.CalculateShipment(ord)
	// todo: just set the order to pending tbh
	if err != nil {
		log.Println("cannot calculate shipping", err)
		return ord, nil
	}
	ctx := context.Background()
	taken, err := s.Catalog.ExecuteShipment(ctx, shipment)
	if err != nil {
		log.Println("no items retrieved", err)
		return ord, nil
	}
	ord.AddReadyItems(taken)
	return ord, nil
}

// GreetWarehouses sends greeting to every warehouse to test connection
func (s *Stock) GreetWarehouses() {
	whs := s.Catalog.GetWarehouses()
	for addr := range whs {
		ctx := context.Background()
		reply, err := warehouse.GreetWarehouse(ctx, addr)
		if err != nil {
			log.Printf("Error when greeting a warehouse %s: %v\n", addr, err)
		} else {
			log.Printf("Warehouse %s replied: %s\n", addr, reply)
		}
	}
}
