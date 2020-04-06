package warehouse

import (
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
	Warehouses map[string][]int
}

// Listen to active warehouses over UDP
// and add new warehouses
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
	items := make([]int, 0)
	addresBook.Warehouses[address] = items
	log.Printf("Added new warehouse by the address: %s\n", address)
	// todo: maybe perform grpc items call synchronously?
	// the only point of doing it async if we add new warehouses very often
	// and don't want this process to be blocked - highly unlikely
	go updateWarehouseItems(address, addresBook)
}

// send request to the given warehouse and add
func updateWarehouseItems(address string, addresBook *AddressBook) {
	var items []int
	// todo: perform grpc call here
	addresBook.Mux.Lock()
	defer addresBook.Mux.Unlock()
	addresBook.Warehouses[address] = items
	log.Printf("Updated items for %s, items: %v\n", address, items)
}
