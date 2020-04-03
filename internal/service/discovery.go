package service

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
	mux        sync.Mutex
	warehouses map[string][]int
}

// Listen to active warehouses over UDP
// and add new warehouses
func discoverWarehouses(addressBook *AddressBook) {
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
	addresBook.mux.Lock()
	defer addresBook.mux.Unlock()
	if _, ok := addresBook.warehouses[address]; ok {
		// warehouse is already added
		return
	}
	items := make([]int, 0)
	addresBook.warehouses[address] = items
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
	addresBook.mux.Lock()
	defer addresBook.mux.Unlock()
	addresBook.warehouses[address] = items
	log.Printf("Updated items for %s, items: %v\n", address, items)
}
