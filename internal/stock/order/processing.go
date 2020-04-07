package order

import "nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"

type Order struct {
	items []int64
}

func Process(o *Order, catalog *warehouse.Catalog) error {
	plan, err := calculatePlan(o, catalog)
	if err != nil {
		return err
	}
	return executePlan(catalog, plan)
}

// requestPlan holds addresses of warehouses from a catalog mapped to
// list of items to request
type requestPlan struct {
	map[string][]int64
}

func calculatePlan(o *Order, catalog *warehouse.Catalog) (*requestPlan, error) {
	return nil, fmt.Errorf("not implemented")
}

func executePlan(catalog *warehouse.Catalog, plan *requestPlan) error {
	return fmt.Errorf("not implemented")
}
