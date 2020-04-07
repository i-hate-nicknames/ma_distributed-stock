package main

import (
	"context"

	wh "nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/web"
)

func main() {
	catalog := wh.MakeCatalog()
	ctx, cancel := context.WithCancel(context.Background())
	go wh.DiscoverWarehouses(catalog)
	go web.StartServer(ctx, "8001", catalog)

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
