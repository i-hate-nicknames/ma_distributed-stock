package warehouse

import (
	"fmt"
	"sync"
	"time"

	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/order"
)

const discoverTimeout = 500 * time.Millisecond

// todo: maybe create a global type for map[string][]int64?
// this is getting on me nerves

// Catalog maps warehouse address to its items
// Each warehouse data is valid when it's added to the catalog
// but may later get invalid, e.g. when a warehouse goes down or
// somehow loses items
type Catalog struct {
	mux        sync.Mutex
	warehouses map[string][]int64
}

// Shipment holds addresses of warehouses from a catalog mapped to
// list of items to request
type Shipment map[string][]int64

// MakeCatalog makes a single instance of an address book
func MakeCatalog() *Catalog {
	warehouses := make(map[string][]int64, 0)
	return &Catalog{warehouses: warehouses}
}

// GetWarehouses returns all warehouses in this catalog
func (c *Catalog) GetWarehouses() map[string][]int64 {
	return c.warehouses
}

// Copy returs a deep copy of this catalog, which may
// be safely modified
func (c *Catalog) Copy() *Catalog {
	res := MakeCatalog()
	for addr, items := range c.warehouses {
		new := make([]int64, len(items))
		copy(new, items)
		res.warehouses[addr] = new
	}
	return res
}

func (c *Catalog) ApplyShipment(shipment map[string][]int64) {

}

// this doesn't perform the request but only removes item
// from the catalog
func (c *Catalog) PopItem(address string) error {
	wh, ok := c.warehouses[address]
	if !ok {
		return fmt.Errorf("warehouse %s not found", address)
	}
	if len(wh) == 0 {
		return fmt.Errorf("warehouse %s is empty", address)
	}
	c.warehouses[address] = wh[1:]
	return nil
}

// AddWarehouse located by this address to the list of warehouses
func (c *Catalog) AddWarehouse(address string, items []int64) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.warehouses[address] = items
}

// AddWarehouse located by this address to the list of warehouses
func (c *Catalog) HasWarehouse(address string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.warehouses[address]
	return ok
}

// AddWarehouse located by this address to the list of warehouses
func (c *Catalog) RemoveWarehouse(address string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.warehouses[address]
	if ok {
		delete(c.warehouses, address)
	}
}

// Calculate shipment orders from a client order and a catalog of warehouses
func (c *Catalog) CalculateShipment(o *order.Order) (Shipment, error) {
	// since we need to mutate passed catalog for calculation, make a copy
	calculationCatalog := c.Copy()
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
func findItem(catalog *Catalog, item int64) (string, error) {
	for address, wh := range catalog.GetWarehouses() {
		if len(wh) > 0 && wh[0] == item {
			return address, nil
		}
	}
	return "", fmt.Errorf("Item %d not found", item)
}

// ExecuteShipment using given warehouse catalog.
// Request items from every warehouse in the shipment
// If at least one of the warehouses failed, return non-nil error
// Return performed shipment which only includes warehouses from which
// items have been successfuly taken
// If all warehouses succeed, returned shipment is identical to the passed
func (c *Catalog) ExecuteShipment(s Shipment) (Shipment, error) {
	executed := make(map[string][]int64)
	for addr, items := range s {
		fmt.Printf("Requesting %v from %s\n", items, addr)
		// todo: perform a grpc take call
		// todo if error, return executed shipping
		executed[addr] = items
	}
	return executed, fmt.Errorf("not implemented")
}
