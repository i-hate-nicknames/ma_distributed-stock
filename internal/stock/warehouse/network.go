package warehouse

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	api "nvm.ga/mastersofcode/golang_2019/stock_distributed/api"
)

// ListenToInvitations starts listening for invitation messages that active
// warehouses send over UDP. Whenever it gets new warehouse address it sends it
// to the addresses channel
func ListenToInvitations(ctx context.Context, addresses chan<- string) {
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
		addresses <- string(buf[:n])
		time.Sleep(discoverTimeout)
	}
}

// getGrpcClient attempts to dial the specified address flag and returns a service
// client and its underlying connection. If it is unable to make a connection,
// it dies.
func getGrpcClient(address string) (*grpc.ClientConn, api.WarehouseServiceClient, error) {
	conn, err := grpc.Dial(address, grpc.WithTimeout(5*time.Second), grpc.WithInsecure())
	if err != nil {
		return nil, nil, fmt.Errorf("did not connect: %v", err)
	}
	return conn, api.NewWarehouseServiceClient(conn), nil
}

// LoadItems returns sequence of items available for retrieving in warehouse
// by the specified address
func LoadItems(ctx context.Context, address string) ([]int64, error) {
	conn, client, err := getGrpcClient(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	rs, err := client.GetItems(ctx, &api.Empty{})
	if err != nil {
		return nil, fmt.Errorf("load inventory from %s: %v", address, err)
	}
	return rs.GetItems(), nil
}

// TakeItems orders warehouse to ship given items. If the order cannot be
// fulfilled precisely, the the operation fails altogether
func TakeItems(ctx context.Context, address string, items []int64) ([]int64, error) {
	conn, client, err := getGrpcClient(address)
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	response, err := client.TakeItems(ctx, &api.ItemList{Items: items})
	if err != nil {
		return nil, fmt.Errorf("take items from %s: %v", address, err)
	}
	return response.Items, nil
}

// GreetWarehouse is a basic wrapper around the corresponding service's RPC.
// It parses the provided arguments, calls the service, and prints the
// response. If any errors are encountered, it dies.
func GreetWarehouse(ctx context.Context, address string) (string, error) {
	conn, client, err := getGrpcClient(address)
	if err != nil {
		return "", err
	}
	defer conn.Close()
	response, err := client.Hello(ctx, &api.Text{Text: "hey"})
	if err != nil {
		return "", err
	}
	return response.Text, nil
}
