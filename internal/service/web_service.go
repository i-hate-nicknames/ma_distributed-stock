package service

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func StartService() {
	items := make(map[string][]int, 0)
	warehouses := &Warehouses{items: items}
	go discoverWarehouses(warehouses)

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		log.Println("Greeting all warehouses")
		greetWarehouses(warehouses)
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

// Send greeting to every warehouse to test connection
func greetWarehouses(warehouses *Warehouses) {
	warehouses.mux.Lock()
	defer warehouses.mux.Unlock()
	for addr := range warehouses.items {
		callWarehouse(addr, "hello")
	}
}

// todo maybe add method to check which warehouses are still
// alive
// maybe make an infinite loop that will periodically check on every warehouse
// and remove those that are dead

// simulate taking items: send take item requests to all available warehouses
func takeItems(warehouses *Warehouses) {
	warehouses.mux.Lock()
	defer warehouses.mux.Unlock()
	for addr := range warehouses.items {
		// todo: move this to requests
		msg := []int{1, 2, 3, 5}
		data, _ := json.Marshal(msg)
		resp, err := http.Post("http://"+addr+"/take", "application/json", bytes.NewReader(data))
		if err != nil {
			log.Printf("Error taking stuff from warehouse at %s, error: %s\n", addr, err)
		}
		defer resp.Body.Close()
	}
}

// todo: implement taking items
// look which warehouses have required items and can satisfy
// the request
// for now just take as much items as possible from every warehouse
// (eventually port our search algorithm to determine from which warehouses to)

// if something breaks down in the process, ignore the warehouse with error

// if we cannot satisfy the order after exhausting all the warehouses, put the order
// on pending

// todo: add orders and order status checking

// todo: add order persistence
