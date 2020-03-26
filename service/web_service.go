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

type Warehouses struct {
	mux   sync.Mutex
	items map[string][]int
}

func main() {
	items := make(map[string][]int, 0)
	warehouses := &Warehouses{items: items}
	go discoverWarehouses(warehouses)

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		log.Println("Greeting all warehouses")
		greetMachines(warehouses)
		c.JSON(http.StatusNoContent, gin.H{})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "I am ok",
		})
	})
	r.GET("/takeSome", func(c *gin.Context) {
		takeItems(warehouses)
	})
	r.Run(":8001")
}

func greetMachines(warehouses *Warehouses) {
	warehouses.mux.Lock()
	defer warehouses.mux.Unlock()
	for addr := range warehouses.items {
		resp, err := http.Get("http://" + addr + "/hello")
		if err != nil {
			log.Printf("Error greeting warehouse at %s, error: %s\n", addr, err)
		}
		defer resp.Body.Close()
	}
}

// simulate taking items: send take item requests to all available warehouses
func takeItems(warehouses *Warehouses) {
	warehouses.mux.Lock()
	defer warehouses.mux.Unlock()
	for addr := range warehouses.items {
		msg := []int{1, 2, 3, 5}
		data, _ := json.Marshal(msg)
		resp, err := http.Post("http://"+addr+"/take", "application/json", bytes.NewReader(data))
		if err != nil {
			log.Printf("Error taking stuff from warehouse at %s, error: %s\n", addr, err)
		}
		defer resp.Body.Close()
	}
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
