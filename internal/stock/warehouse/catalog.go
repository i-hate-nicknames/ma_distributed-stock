package warehouse

import (
	"fmt"
	"sync"
	"time"

	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/order"
)

const discoverTimeout = 500 * time.Millisecond

// Inventory holds a collection of warehouse addresses mapped
// a set of items. This can be used to represent a total inventory of
// a warehouse as well as an inventory of a shipping
type Inventory map[string][]int64

// Catalog maps warehouse address to its items
// Each warehouse data is valid when it's added to the catalog
// but may later get invalid, e.g. when a warehouse goes down or
// somehow loses items
type Catalog struct {
	mux        sync.Mutex
	warehouses Inventory
}

// MakeCatalog makes a single instance of an address book
func MakeCatalog() *Catalog {
	warehouses := make(Inventory, 0)
	return &Catalog{warehouses: warehouses}
}

// popItem takes and removes first item in the warehouse
// in the given inventory. If there is no such warehouse, or
// it has no items, return non-nil error
func popItem(inv Inventory, address string) error {
	wh, ok := inv[address]
	if !ok {
		return fmt.Errorf("warehouse %s not found", address)
	}
	if len(wh) == 0 {
		return fmt.Errorf("warehouse %s is empty", address)
	}
	inv[address] = wh[1:]
	return nil
}

// AddWarehouse located by this address to the list of warehouses
func (c *Catalog) AddWarehouse(address string, items []int64) {
	c.mux.Lock()
	defer c.mux.Unlock()
	c.warehouses[address] = items
}

// HasWarehouse returns true if this catalog has a warehouse under
// given address
func (c *Catalog) HasWarehouse(address string) bool {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.warehouses[address]
	return ok
}

// RemoveWarehouse identified by the address from the catalog
func (c *Catalog) RemoveWarehouse(address string) {
	c.mux.Lock()
	defer c.mux.Unlock()
	_, ok := c.warehouses[address]
	if ok {
		delete(c.warehouses, address)
	}
}

// CalculateShipment orders from a client order and a catalog of warehouses
func (c *Catalog) CalculateShipment(o *order.Order) (Inventory, error) {
	inv := make(Inventory)
	for addr, items := range c.warehouses {
		new := make([]int64, len(items))
		copy(new, items)
		inv[addr] = new
	}
	// for every item in the order, check all warehouses if they have it
	// use the first you come upon
	orders := make(Inventory)
	for _, orderItem := range o.Items {
		// todo: wrap errors properly
		address, err := findItem(inv, orderItem)
		if err != nil {
			return nil, err
		}
		err = popItem(inv, address)
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
func findItem(inv Inventory, item int64) (string, error) {
	for address, wh := range inv {
		if len(wh) > 0 && wh[0] == item {
			return address, nil
		}
	}
	return "", fmt.Errorf("Item %d not found", item)
}

// ExecuteShipment using given warehouse catalog.
// Request items from every warehouse in the shipment
// If at least one of the warehouses failed, return non-nil error
func (c *Catalog) ExecuteShipment(shipment Inventory) error {
	executed := make(Inventory)
	for addr, items := range shipment {
		fmt.Printf("Requesting %v from %s\n", items, addr)
		// todo: perform a grpc take call
		// todo if error, return executed shipping
		executed[addr] = items
	}
	// todo:
	// applying a shipment introduces a consistency issue:
	// when shipment A is executed but haven't updated
	// the catalog yet, and shipment B is calculated on
	// the old inventory. Upon finish shipment A will
	// apply whatever has been executed thus changing inventory,
	// and shipment B may become invalid
	// we can simply allow that to happen and then ensure that
	// shipping operations are performed in a single "transaction",
	// or we can introduce some mechanism of validation, the simplest
	// may be adding "last calculated shipment" integer to catalog
	c.applyShipment(executed)
	return fmt.Errorf("not implemented")
}

func (c *Catalog) applyShipment(shipment Inventory) {
	c.mux.Lock()
	defer c.mux.Unlock()
}
