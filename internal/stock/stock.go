package stock

import (
	"context"
	"log"

	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/order"
	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
)

type Stock struct {
	Warehouses *warehouse.Catalog
	Orders     *order.Registry
}

func (s *Stock) DiscoverWarehouses(ctx context.Context) {
	addresses := make(chan string, 5)
	go warehouse.ListenToInvitations(ctx, addresses)
	for address := range addresses {
		if s.Warehouses.HasWarehouse(address) {
			continue
		}
		s.updateWarehouseItems(ctx, address)
	}
}

func (s *Stock) updateWarehouseItems(ctx context.Context, address string) {
	items, err := warehouse.LoadItems(ctx, address)
	if err != nil {
		s.Warehouses.RemoveWarehouse(address)
		return
	}
	log.Printf("Added warehouse %s with items %v\n", address, items)
	s.Warehouses.AddWarehouse(address, items)
}

func (s *Stock) SumbitOrder(items []int64) (*order.Order, error) {
	// todo: this method is messy because later order processing
	// will happen asynchronously in a different thread, and order
	// submission will just create an order
	ord, err := s.Orders.SubmitOrder(items)
	if err != nil {
		return nil, err
	}
	shipment, err := s.Warehouses.CalculateShipment(ord)
	if err != nil {
		return nil, err
	}
	err = s.Warehouses.ExecuteShipment(shipment)
	if err != nil {
		// executing was not entirely successful
		// collect all the items in the shipment, store
		// them as shipped in the order, and remove those items
		// from items in order
		// set order status to pending

		// we still consider a partial order satisfaction a success
		return ord, nil
	}
	return ord, nil
}

// GreetWarehouses sends greeting to every warehouse to test connection
func (s *Stock) GreetWarehouses() {
	whs := s.Warehouses.GetWarehouses()
	for addr := range whs {
		ctx := context.Background()
		warehouse.GreetWarehouse(ctx, addr)
	}
}
