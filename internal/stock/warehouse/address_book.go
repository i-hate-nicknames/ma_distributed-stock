package warehouse

import (
	"context"
	"log"
	"net"
	"sync"
	"time"
)

const discoverTimeout = 500 * time.Millisecond

// AddressBook map warehouse address to
// the items that it has
type AddressBook struct {
	Mux        sync.Mutex
	Warehouses map[string][]int64
}

// MakeAddressBook makes a single instance of an address book
func MakeAddressBook() *AddressBook {
	warehouses := make(map[string][]int64, 0)
	return &AddressBook{Warehouses: warehouses}
}

// DiscoverWarehouses starts listening for invitation messages that active
// warehouses send over UDP. It then adds new warehouses to the address book
// This is a blocking operation that blocks indefinitely
// todo: add context for cancelation
func DiscoverWarehouses(addressBook *AddressBook) {
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
		addWarehouse(string(buf[:n]), addressBook)
		time.Sleep(discoverTimeout)
	}
}

// Add a warehouse located by this address to the list of warehouses
func addWarehouse(address string, addresBook *AddressBook) {
	addresBook.Mux.Lock()
	defer addresBook.Mux.Unlock()
	if _, ok := addresBook.Warehouses[address]; ok {
		// warehouse is already added
		return
	}
	items := make([]int64, 0)
	addresBook.Warehouses[address] = items
	log.Printf("Added new warehouse by the address: %s\n", address)
	go updateWarehouseItems(address, addresBook)
}

// send request to the given warehouse and add
func updateWarehouseItems(address string, addresBook *AddressBook) {
	ctx := context.Background()
	items, err := doGetItems(ctx, address)
	addresBook.Mux.Lock()
	defer addresBook.Mux.Unlock()
	if err != nil {
		log.Println(err)
		delete(addresBook.Warehouses, address)
		return
	}
	addresBook.Warehouses[address] = items
	log.Printf("Updated items for %s, items: %v\n", address, items)
}

// TakeItems simulates taking items: send take item requests to all available warehouses
func TakeItems(addressBook *AddressBook) {
	addressBook.Mux.Lock()
	defer addressBook.Mux.Unlock()
	toTake := []int64{1, 2}
	for addr := range addressBook.Warehouses {
		log.Printf("Taking %v from %s\n", toTake, addr)
		ctx := context.Background()
		err := doTakeItems(ctx, addr, toTake)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

// GreetWarehouses sends greeting to every warehouse to test connection
func GreetWarehouses(addressBook *AddressBook) {
	addressBook.Mux.Lock()
	defer addressBook.Mux.Unlock()
	for addr := range addressBook.Warehouses {
		ctx := context.Background()
		doHello(ctx, addr)
	}
}
