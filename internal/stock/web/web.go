package web

import (
	"context"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	wh "nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
)

// StartServer starts a web server that listens to incoming requests and performs
// corresponding actions using available warehouses
func StartServer(ctx context.Context, port string, addressBook *wh.AddressBook) {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		log.Println("Greeting all warehouses")
		wh.GreetWarehouses(addressBook)
		c.JSON(http.StatusNoContent, gin.H{})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "I am ok",
		})
	})
	r.GET("/takeSome", func(c *gin.Context) {
		wh.TakeItems(addressBook)
	})
	r.Run(":" + port)
}

// todo: implement taking items
// look which warehouses have required items and can satisfy
// the request
// for now just take as much items as possible from every warehouse
// (eventually port our search algorithm to determine from which warehouses to)
// if something breaks down in the process, ignore the warehouse with error
// if we cannot satisfy the order after exhausting all the warehouses, put the order
// on pending
