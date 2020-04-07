package warehouse

import (
	"context"
	"log"
	"net"
	"sync"
	"time"
)

const discoverTimeout = 500 * time.Millisecond

// Catalog maps warehouse address to its items
// Each warehouse data is valid when it's added to the catalog
// but may later get invalid, e.g. when a warehouse goes down or
// somehow loses items
type Catalog struct {
	mux        sync.Mutex
	warehouses map[string][]int64
}

// MakeCatalog makes a single instance of an address book
func MakeCatalog() *Catalog {
	warehouses := make(map[string][]int64, 0)
	return &Catalog{warehouses: warehouses}
}

// DiscoverWarehouses starts listening for invitation messages that active
// warehouses send over UDP. It then adds new warehouses to the address book
// This is a blocking operation that blocks indefinitely
// todo: add context for cancelation
func DiscoverWarehouses(catalog *Catalog) {
	conn, err := net.ListenUDP("udp", &net.UDPAddr{
		Port: 3000,
		IP:   net.ParseIP("0.0.0.0"),
	})
	if err != nil {
		log.Println(err)
	}
	defer conn.Close()

	buf := make([]byte, 1024)
	for {
		n, _, err := conn.ReadFrom(buf[:])
		if err != nil {
			log.Println(err)
		}
		addWarehouse(string(buf[:n]), catalog)
		time.Sleep(discoverTimeout)
	}
}

// Add a warehouse located by this address to the list of warehouses
func addWarehouse(address string, catalog *Catalog) {
	catalog.mux.Lock()
	defer catalog.mux.Unlock()
	if _, ok := catalog.warehouses[address]; ok {
		// warehouse is already added
		return
	}
	items := make([]int64, 0)
	catalog.warehouses[address] = items
	log.Printf("Added new warehouse by the address: %s\n", address)
	go updateWarehouseItems(address, catalog)
}

// send request to the given warehouse and add
func updateWarehouseItems(address string, catalog *Catalog) {
	ctx := context.Background()
	items, err := doGetItems(ctx, address)
	catalog.mux.Lock()
	defer catalog.mux.Unlock()
	if err != nil {
		log.Println(err)
		delete(catalog.warehouses, address)
		return
	}
	catalog.warehouses[address] = items
	log.Printf("Updated items for %s, items: %v\n", address, items)
}

// TakeItems simulates taking items: send take item requests to all available warehouses
func TakeItems(catalog *Catalog) {
	catalog.mux.Lock()
	defer catalog.mux.Unlock()
	toTake := []int64{1, 2}
	for addr := range catalog.warehouses {
		log.Printf("Taking %v from %s\n", toTake, addr)
		ctx := context.Background()
		err := doTakeItems(ctx, addr, toTake)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

// GreetWarehouses sends greeting to every warehouse to test connection
func GreetWarehouses(catalog *Catalog) {
	catalog.mux.Lock()
	defer catalog.mux.Unlock()
	for addr := range catalog.warehouses {
		ctx := context.Background()
		doHello(ctx, addr)
	}
}
