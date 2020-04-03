package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	api "nvm.ga/mastersofcode/golang_2019/stock_distributed/api"
)

func StartService() {
	warehouses := make(map[string][]int, 0)
	addressBook := &AddressBook{warehouses: warehouses}
	go discoverWarehouses(addressBook)

	r := gin.Default()
	r.GET("/hello", func(c *gin.Context) {
		log.Println("Greeting all warehouses")
		greetWarehouses(addressBook)
		c.JSON(http.StatusNoContent, gin.H{})
	})
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "I am ok",
		})
	})
	r.GET("/takeSome", func(c *gin.Context) {
		takeItems(addressBook)
	})
	r.Run(":8001")
}

// Send greeting to every warehouse to test connection
func greetWarehouses(addressBook *AddressBook) {
	addressBook.mux.Lock()
	defer addressBook.mux.Unlock()
	for addr := range addressBook.warehouses {
		ctx := context.Background()
		doHello(ctx, addr)
	}
}

// todo maybe add method to check which warehouses are still
// alive
// maybe make an infinite loop that will periodically check on every warehouse
// and remove those that are dead

// simulate taking items: send take item requests to all available warehouses
func takeItems(addressBook *AddressBook) {
	addressBook.mux.Lock()
	defer addressBook.mux.Unlock()
	for addr := range addressBook.warehouses {
		// todo: move this to requests
		log.Println("Taking items from", addr)
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

// todo: move grpc client code to a separate file

// GetClient attempts to dial the specified address flag and returns a service
// client and its underlying connection. If it is unable to make a connection,
// it dies.
func GetClient(address string) (*grpc.ClientConn, api.WarehouseServiceClient) {
	conn, err := grpc.Dial(address, grpc.WithTimeout(5*time.Second), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn, api.NewWarehouseServiceClient(conn)
}

// doList is a basic wrapper around the corresponding book service's RPC.
// It parses the provided arguments, calls the service, and prints the
// response. If any errors are encountered, it dies.
func doHello(ctx context.Context, address string) {
	conn, client := GetClient(address)
	defer conn.Close()
	rs, err := client.Hello(ctx, &api.Text{Text: "hey"})
	if err != nil {
		log.Fatalf("List books: %v", err)
	}
	fmt.Printf("Server replied to our greeting: %s\n", rs.GetText())
}
