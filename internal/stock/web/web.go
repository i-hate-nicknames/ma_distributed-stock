package web

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/order"
	wh "nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
)

type takeItemsReq struct {
	items []int64
}

// StartServer starts a web server that listens to incoming requests and performs
// corresponding actions using available warehouses
func StartServer(ctx context.Context, port string, catalog *wh.Catalog) {
	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		log.Println("Greeting all warehouses")
		wh.GreetWarehouses(catalog)
		c.JSON(http.StatusNoContent, gin.H{})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "I am ok",
		})
	})
	r.GET("/takeSome", func(c *gin.Context) {
		wh.TakeItems(catalog)
	})
	r.POST("/take", func(c *gin.Context) {
		take(c, catalog)
	})
	r.Run(":" + port)
}

func take(c *gin.Context, catalog *wh.Catalog) {
	var req takeItemsReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	log.Printf("Taking items %v", req.items)
	ord := &order.Order{Items: req.items}
	err = order.Process(ord, catalog)
	if err != nil {
		msg := fmt.Sprintf("Cannot execute the order: %s", err.Error())
		c.JSON(http.StatusNotFound, gin.H{"message": msg})
		return
	}
	c.JSON(http.StatusOK, gin.H{"status": "order successfuly executed"})
}

// todo: implement taking items
// look which warehouses have required items and can satisfy
// the request
// for now just take as much items as possible from every warehouse
// (eventually port our search algorithm to determine from which warehouses to)
// if something breaks down in the process, ignore the warehouse with error
// if we cannot satisfy the order after exhausting all the warehouses, put the order
// on pending
