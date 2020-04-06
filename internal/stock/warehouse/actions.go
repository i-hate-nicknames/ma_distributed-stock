package warehouse

import (
	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
	api "nvm.ga/mastersofcode/golang_2019/stock_distributed/api"
)

// GetClient attempts to dial the specified address flag and returns a service
// client and its underlying connection. If it is unable to make a connection,
// it dies.
func getClient(address string) (*grpc.ClientConn, api.WarehouseServiceClient) {
	conn, err := grpc.Dial(address, grpc.WithTimeout(5*time.Second), grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	return conn, api.NewWarehouseServiceClient(conn)
}

// doHello is a basic wrapper around the corresponding service's RPC.
// It parses the provided arguments, calls the service, and prints the
// response. If any errors are encountered, it dies.
func doHello(ctx context.Context, address string) {
	conn, client := getClient(address)
	defer conn.Close()
	rs, err := client.Hello(ctx, &api.Text{Text: "hey"})
	if err != nil {
		log.Fatalf("Hello: %v", err)
	}
	fmt.Printf("Warehouse replied to our greeting: %s\n", rs.GetText())
}

// todo: add other actions like taking items or getting item list

func doGetItems(ctx context.Context, address string) ([]int64, error) {
	conn, client := getClient(address)
	defer conn.Close()
	rs, err := client.GetItems(ctx, &api.Empty{})
	if err != nil {
		return nil, fmt.Errorf("get items from %s: %v", address, err)
	}
	return rs.GetItems(), nil
}
