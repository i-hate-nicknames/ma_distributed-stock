package stock

import (
	"context"

	wh "nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/web"
)

func StartService() {
	warehouses := make(map[string][]int, 0)
	addressBook := &wh.AddressBook{Warehouses: warehouses}
	ctx, _ := context.WithCancel(context.Background())
	go wh.DiscoverWarehouses(addressBook)
	go web.StartServer(ctx, "8001", addressBook)

	// todo: listen to cancellation signals
	done := make(chan struct{}, 1)
	<-done
}

// todo maybe add method to check which warehouses are still
// alive
// maybe make an infinite loop that will periodically check on every warehouse
// and remove those that are dead

// todo: add orders and order status checking

// todo: add order persistence
