package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net"
	"net/http"
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
	resp, err := http.Get("http://" + address + "/getItems")
	if err != nil {
		log.Printf("Error getting warehouse items, address: %s, error: %s\n", address, err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error reading response data: %s\n", err)
		return
	}
	var items []int
	err = json.Unmarshal(body, &items)
	if err != nil {
		log.Printf("Malformed response: %s\n", err)
		return
	}
	warehouses.mux.Lock()
	defer warehouses.mux.Unlock()
	warehouses.items[address] = items
	log.Printf("Updated items for %s, items: %v\n", address, items)
}
