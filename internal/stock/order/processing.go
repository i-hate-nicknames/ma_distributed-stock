package order

import (
	"fmt"

	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
)

// shipmentOrders holds addresses of warehouses from a catalog mapped to
// list of items to request
type shipmentOrders map[string][]int64

// Calculate shipment orders from a client order and a catalog of warehouses
func CalculateOrders(o *Order, catalog *warehouse.Catalog) (shipmentOrders, error) {
	// since we need to mutate passed catalog for calculation, make a copy
	calculationCatalog := catalog.Copy()
	// for every item in the order, check all warehouses if they have it
	// use the first you come upon
	orders := make(map[string][]int64)
	for _, orderItem := range o.Items {
		// todo: wrap errors properly
		address, err := findItem(calculationCatalog, orderItem)
		if err != nil {
			return nil, err
		}
		err = calculationCatalog.PopItem(address)
		if err != nil {
			return nil, err
		}
		if _, ok := orders[address]; !ok {
			orders[address] = make([]int64, 0)
		}
		orders[address] = append(orders[address], orderItem)
	}
	return orders, nil
}

// find a warehouse that has given item on the top of its queue
// return address of that warehouse, or error if item cannot be found
func findItem(catalog *warehouse.Catalog, item int64) (string, error) {
	for address, wh := range catalog.GetWarehouses() {
		if len(wh) > 0 && wh[0] == item {
			return address, nil
		}
	}
	return "", fmt.Errorf("Item %d not found", item)
}

// Execute given orders with warehouse catalog. Request items from every warehouse
// via grpc
// If at least one of the warehouses failed, return non-nil error
// Return an updated request orders that only includes warehouses that have been
// successfuly queried, so that in case of a partial result a client can return
// items back to warehouses
// In case all warehouses succeeded both initial and result request orderss are identical
func ExecuteOrders(catalog *warehouse.Catalog, orders shipmentOrders) (shipmentOrders, error) {
	executed := make(map[string][]int64)
	for addr, items := range orders {
		fmt.Printf("Requesting %v from %s\n", items, addr)
		// todo: perform a grpc take call
		// todo if error, return executed orders
		executed[addr] = items
	}
	return executed, fmt.Errorf("not implemented")
}
