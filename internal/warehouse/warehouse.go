package warehouse

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	api "nvm.ga/mastersofcode/golang_2019/stock_distributed/api"
)

const invitationTimeout = 1 * time.Millisecond

// todo: instead of itemlist store items internally
// in a LIFO machine
var myItems = &api.ItemList{}
var mux sync.Mutex

// StartWarehouse starts a single warehouse instance. This warehouse will
// continuously broadcast invitations via UDP
// Also it starts a grpc server for interaction with items in this warehouse
func StartWarehouse(port string, items []int64) {
	myItems.Items = items
	addr := "127.0.0.1:" + port
	go sendInvitations(addr)
	startGrpcServer(port)
}

// continuously broadcast invitation message over UDP
// with address to connect
func sendInvitations(myAddr string) {
	serverAddr, err := net.ResolveUDPAddr("udp", "255.255.255.255:3000")
	if err != nil {
		log.Fatal(err)
	}
	localAddr, err := net.ResolveUDPAddr("udp", ":0")
	if err != nil {
		log.Fatal(err)
	}
	conn, err := net.DialUDP("udp", localAddr, serverAddr)
	if err != nil {
		log.Fatal(err)
	}
	msg := []byte(myAddr)
	for {
		time.Sleep(invitationTimeout)
		_, err := conn.Write(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func startGrpcServer(port string) {
	addr := fmt.Sprintf("localhost:%s", port)
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	server := grpc.NewServer(opts...)
	fmt.Println("Running server on", addr)
	service := whService{}
	api.RegisterWarehouseServiceServer(server, &service)
	server.Serve(lis)
}

// empty type that is used to implement grpc server interface
type whService struct{}

// PutItems adds given items to this warehouse in the order left to right, so the last
// item on the list will end up being on top
func (service *whService) PutItems(ctx context.Context, itemList *api.ItemList) (*api.ItemList, error) {
	return nil, status.Errorf(codes.Unimplemented, "Not implemented")
}

// TakeItems retrieves multiple items from this warehouse, while removing them locally
// If items in the request cannot be satisfied in the order they are provided an error is returned
func (service *whService) TakeItems(ctx context.Context, itemList *api.ItemList) (*api.ItemList, error) {
	mux.Lock()
	defer mux.Unlock()
	log.Printf("Processing items request %v, having %v", itemList.Items, myItems.Items)
	requestItems := itemList.GetItems()
	if len(requestItems) > len(myItems.Items) {
		log.Println("Too few items in the warehouse")
		return nil, status.Errorf(codes.FailedPrecondition, "Too few items in the warehouse")
	}
	for idx, item := range requestItems {
		if item != myItems.Items[idx] {
			errorText := fmt.Sprintf("Items differ at %dth element: %v, %v", idx, requestItems, myItems.Items)
			return nil, status.Errorf(codes.FailedPrecondition, errorText)
		}
	}
	myItems.Items = myItems.Items[len(requestItems):]
	return myItems, nil
}

// GetItems returns what items are available in this warehouse in the order they must be requested
func (service *whService) GetItems(ctx context.Context, empty *api.Empty) (*api.ItemList, error) {
	return myItems, nil
}

// Hello is a test method to check that grpc works properly
func (service *whService) Hello(ctx context.Context, text *api.Text) (*api.Text, error) {
	txt := text.GetText()
	return &api.Text{Text: fmt.Sprintln("Hello there", txt)}, nil
}
