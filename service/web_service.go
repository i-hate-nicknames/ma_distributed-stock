package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
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
