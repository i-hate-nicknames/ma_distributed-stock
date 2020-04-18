package main

import (
	"context"

	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock"
	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/order"
	wh "nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
)

func main() {
	catalog := wh.MakeCatalog()
	ctx, cancel := context.WithCancel(context.Background())
	st := &stock.Stock{Catalog: catalog, Orders: order.MakeRegistry()}
	go st.DiscoverWarehouses(ctx)
	go stock.StartWebServer(ctx, "8001", st)

	// todo: listen to cancellation signals
	done := make(chan struct{}, 1)
	<-done
	cancel()
}

// todo maybe add method to check which warehouses are still
// alive
// maybe make an infinite loop that will periodically check on every warehouse
// and remove those that are dead

// todo: add orders and order status checking

// todo: add order persistence
