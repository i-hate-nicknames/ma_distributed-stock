package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

const discoverTimeout = 500 * time.Millisecond

type RemoteMachines struct {
	mux   sync.Mutex
	items map[string][]int
}

func main() {
	items := make(map[string][]int, 0)
	machines := &RemoteMachines{items: items}
	go discoverMachines(machines)

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		log.Println("Greeting all machines")
		greetMachines(machines)
		c.JSON(http.StatusNoContent, gin.H{})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "I am ok",
		})
	})
	r.GET("/takeSome", func(c *gin.Context) {
		takeItems(machines)
	})
	r.Run(":8001")
}

func greetMachines(machines *RemoteMachines) {
	machines.mux.Lock()
	defer machines.mux.Unlock()
	for addr := range machines.items {
		resp, err := http.Get("http://" + addr + "/hello")
		if err != nil {
			log.Printf("Error greeting machine at %s, error: %s\n", addr, err)
		}
		defer resp.Body.Close()
	}
}

// simulate taking items: send take item requests to all available machines
func takeItems(machines *RemoteMachines) {
	machines.mux.Lock()
	defer machines.mux.Unlock()
	for addr := range machines.items {
		msg := []int{1, 2, 3, 5}
		data, _ := json.Marshal(msg)
		resp, err := http.Post("http://"+addr+"/take", "application/json", bytes.NewReader(data))
		if err != nil {
			log.Printf("Error taking stuff from machine at %s, error: %s\n", addr, err)
		}
		defer resp.Body.Close()
	}
}

func discoverMachines(machines *RemoteMachines) {
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
		addMachine(string(buf[:n]), machines)
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
