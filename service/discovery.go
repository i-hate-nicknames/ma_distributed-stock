package main

import (
	"log"
	"net"
	"sync"
	"time"
)

const discoverTimeout = 500 * time.Millisecond

type Warehouses struct {
	mux   sync.Mutex
	items map[string][]int
}

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

func addWarehouse(address string, warehouses *Warehouses) {
	warehouses.mux.Lock()
	defer warehouses.mux.Unlock()
	if _, ok := warehouses.items[address]; ok {
		// warehouse is already added
		return
	}
	// todo: parse message and take items from there
	items := make([]int, 0)
	warehouses.items[address] = items
	log.Printf("Added new warehouse with the following items %v\n", items)
}
