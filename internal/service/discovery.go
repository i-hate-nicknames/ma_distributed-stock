package service

import (
	"encoding/json"
	"log"
	"net"
	"sync"
	"time"
)

const discoverTimeout = 500 * time.Millisecond

// Warehouses map warehouse address to
// the items that it
type Warehouses struct {
	mux   sync.Mutex
	items map[string][]int
}

// Listen to active warehouses over UDP
// and add new warehouses
func discoverWarehouses(warehouses *Warehouses) {
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
		addWarehouse(string(buf[:n]), warehouses)
		time.Sleep(discoverTimeout)
	}
}

// Add a warehouse located by this address to the list of warehouses
func addWarehouse(address string, warehouses *Warehouses) {
	warehouses.mux.Lock()
	defer warehouses.mux.Unlock()
	if _, ok := warehouses.items[address]; ok {
		// warehouse is already added
		return
	}
	items := make([]int, 0)
	warehouses.items[address] = items
	log.Printf("Added new warehouse by the address: %s\n", address)
	go updateWarehouseItems(address, warehouses)
}

// send request to the given warehouse and add
func updateWarehouseItems(address string, warehouses *Warehouses) {
	var items []int
	body, err := callWarehouse(address, "getItems")
	if err != nil {
		log.Println(err)
		return
	}
	err = json.Unmarshal(body, &items)
	if err != nil {
		log.Printf("Malformed response from %s: %s", address, err)
		return
	}
	warehouses.mux.Lock()
	defer warehouses.mux.Unlock()
	warehouses.items[address] = items
	log.Printf("Updated items for %s, items: %v\n", address, items)
}
