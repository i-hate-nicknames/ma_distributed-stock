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

func LoadInventory(ctx context.Context, address string) ([]int64, error) {
	conn, client := getClient(address)
	defer conn.Close()
	rs, err := client.GetItems(ctx, &api.Empty{})
	if err != nil {
		return nil, fmt.Errorf("load inventory from %s: %v", address, err)
	}
	return rs.GetItems(), nil
}

// todo: separate catalogue manipulation code from warehouse network requests
// and from warehouse discovery

// GetInvitations starts listening for invitation messages that active
// warehouses send over UDP. Whenever it gets new warehouse address it sends it
// to the addresses channel
func GetInvitations(ctx context.Context, addresses chan<- string) {
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

// TakeItems simulates taking items: send take item requests to all available warehouses
func TakeItems(catalog *Catalog) {
	catalog.mux.Lock()
	defer catalog.mux.Unlock()
	toTake := []int64{1, 2}
	for addr := range catalog.warehouses {
		log.Printf("Taking %v from %s\n", toTake, addr)
		ctx := context.Background()
		err := doTakeItems(ctx, addr, toTake)
		if err != nil {
			log.Printf(err.Error())
		}
	}
}

// GreetWarehouses sends greeting to every warehouse to test connection
func GreetWarehouses(catalog *Catalog) {
	catalog.mux.Lock()
	defer catalog.mux.Unlock()
	for addr := range catalog.warehouses {
		ctx := context.Background()
		doHello(ctx, addr)
	}
}

func doTakeItems(ctx context.Context, address string, items []int64) error {
	conn, client := getClient(address)
	defer conn.Close()
	_, err := client.TakeItems(ctx, &api.ItemList{Items: items})
	if err != nil {
		return fmt.Errorf("take items from %s: %v", address, err)
	}
	return nil
}
