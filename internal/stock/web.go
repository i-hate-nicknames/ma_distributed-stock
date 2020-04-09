package stock

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/order"
	wh "nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/stock/warehouse"
)

type itemsReq struct {
	Items []int64
}

type idReq struct {
	Id uint
}

// StartServer starts a web server that listens to incoming requests and performs
// corresponding actions using available warehouses
func StartServer(ctx context.Context, port string, stock *Stock) {
	r := gin.Default()
	var takeMux sync.Mutex
	r.GET("/hello", func(c *gin.Context) {
		log.Println("Greeting all warehouses")
		wh.GreetWarehouses(stock.Warehouses)
		c.JSON(http.StatusNoContent, gin.H{})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "I am ok",
		})
	})

	r.POST("/take", func(c *gin.Context) {
		// lock take operation since it's not thread-safe
		takeMux.Lock()
		defer takeMux.Unlock()
		take(c, stock)
	})

	// Order management handlers
	r.POST("/submit", func(c *gin.Context) {
		var req itemsReq
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		ID, err := stock.Orders.SubmitOrder(req.Items)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "failed to create an order"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"orderId": ID})
	})

	r.GET("/getStatus", func(c *gin.Context) {
		// todo: this is a GET request, get the id appropriately
		var req idReq
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		ord, ok := stock.Orders.GetOrder(req.Id)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{})
			return
		}
		c.JSON(http.StatusOK, gin.H{"status": ord.Status})

	})

	// todo: currently sets order state to canceled
	// when order scheduler is implemented the status should
	// be pendingCancel that denote that the order is planned to
	// be canceled
	r.POST("/cancel", func(c *gin.Context) {
		var req idReq
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{})
			return
		}
		err = stock.Orders.CancelOrder(req.Id)
		// todo: not sure how to distinguish between not found
		// and failed to update errors here
		if err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusNoContent, gin.H{})
	})
	r.Run(":" + port)
}

func take(c *gin.Context, stock *Stock) {
	var req itemsReq
	err := c.BindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
		return
	}
	log.Printf("Taking items %v", req.Items)
	ord := &order.Order{Items: req.Items}
	err = order.Process(ord, stock.Warehouses)
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
