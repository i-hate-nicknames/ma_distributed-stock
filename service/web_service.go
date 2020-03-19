package main

import (
	"log"
	"net"
	"sync"
	"time"
)

const discoverTimeout = 500 * time.Millisecond

type RemoteMachines struct {
	mux   sync.Mutex
	items map[string][]int
}

func main() {
	items := make(map[string][]int, 0)
	machines := &RemoteMachines{items: items}
	// todo: add web service handlers here

	discoverStocks(machines)
}

func discoverStocks(machines *RemoteMachines) {
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
		log.Println("Reding next portion")
		_, _, err := conn.ReadFrom(buf[:])
		if err != nil {
			log.Println(err)
		}
		addMachine(string(buf), machines)
		time.Sleep(discoverTimeout)
	}
}

func addMachine(address string, machines *RemoteMachines) {
	machines.mux.Lock()
	defer machines.mux.Unlock()
	if _, ok := machines.items[address]; ok {
		// machine is already added
		return
	}
	// todo: parse message and take items from there
	items := make([]int, 0)
	machines.items[address] = items
	log.Printf("Added new machine with the following items %v\n", items)
}
