package stock

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type itemsReq struct {
	Items []int64
}

type idReq struct {
	Id uint
}

// StartWebServer starts a web server that listens to incoming requests and performs
// corresponding actions using available warehouses
func StartWebServer(ctx context.Context, port string, stock *Stock) {
	r := gin.Default()
	r.GET("/sendGreetings", func(c *gin.Context) {
		log.Println("Greeting all warehouses")
		stock.GreetWarehouses()
		c.JSON(http.StatusNoContent, gin.H{})
	})

	// Order management handlers
	r.POST("/submit", func(c *gin.Context) {
		var req itemsReq
		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		order, err := stock.SubmitOrder(req.Items)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"orderId": order.ID})
	})

	r.GET("/order/:id", func(c *gin.Context) {
		idParam := c.Param("id")
		orderID, err := strconv.Atoi(idParam)
		if err != nil || orderID <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a positive integer"})
			return
		}
		ord, ok := stock.Orders.GetOrder(uint(orderID))
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
