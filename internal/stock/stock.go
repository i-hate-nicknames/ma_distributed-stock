package stock

import (
	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/order"
	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
)

type Stock struct {
	Warehouses *warehouse.Catalog
	Orders     *order.Registry
}
