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
