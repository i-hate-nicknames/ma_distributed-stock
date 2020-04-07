package order

import (
	"fmt"

	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
)

// Order represents a user order for a set of items
type Order struct {
	Items []int64
}

// Process order: try finding order items in the warehouses in the catalog
// if the order is satisfiable then carry out the order by requesting items
// from the warehouses
// otherwise, return an error
// This function is not thread safe
func Process(o *Order, catalog *warehouse.Catalog) error {
	// todo: we might want to lock the catalog here to ensure
	// that only one order is processed at a time
	// but probably is not worth it since order processing will
	// most likely be directed through a scheduler that runs in a single thread
	plan, err := calculatePlan(o, catalog)
	if err != nil {
		return err
	}
	return executePlan(catalog, plan)
}

// requestPlan holds addresses of warehouses from a catalog mapped to
// list of items to request
type requestPlan map[string][]int64

func calculatePlan(o *Order, catalog *warehouse.Catalog) (requestPlan, error) {
	// for every item in the order, check all warehouses if they have it
	// use the first you come upon
	plan := make(map[string][]int64)

	for _, orderItem := range o.Items {
		// todo: wrap errors properly
		address, err := findItem(catalog, orderItem)
		if err != nil {
			return nil, err
		}
		err = catalog.PopItem(address)
		if err != nil {
			return nil, err
		}
		if _, ok := plan[address]; !ok {
			plan[address] = make([]int64, 0)
		}
		plan[address] = append(plan[address], orderItem)
	}
	return plan, nil
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

func executePlan(catalog *warehouse.Catalog, plan requestPlan) error {
	// for every warehouse in the plan, request planned items from the wh
	// if error happens along any, return error
	// for now this assumes that items are destroyed.
	// todo: in future we need to keep track of the taken items,
	// and either send cancellation commands to warehouses or to
	// collect taken items, most likely both
	for addr, items := range plan {
		fmt.Printf("Requesting %v from %s\n", items, addr)
		// todo: perform a grpc take call
	}
	return fmt.Errorf("not implemented")
}
