package warehouse

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	api "nvm.ga/mastersofcode/golang_2019/stock_distributed/api"
)

const invitationTimeout = 500 * time.Millisecond

var itemList = &api.ItemList{}

func StartWarehouse(port string, items []int64) {
	itemList.Items = items
	addr := "127.0.0.1:" + port
	go sendInvitations(addr)
	startGrpcServer(port)
}

// continuously broadcast invitation message over UDP
// with address to connect
func sendInvitations(myAddr string) {
	for {
		time.Sleep(invitationTimeout)
		con, _ := net.Dial("udp", "127.0.0.1:3000")
		buf := []byte(myAddr)
		_, err := con.Write(buf)
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

func (service *whService) PutItems(ctx context.Context, itemList *api.ItemList) (*api.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "Not implemented")
}

func (service *whService) TakeItems(ctx context.Context, itemList *api.ItemList) (*api.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "Not implemented")
}

func (service *whService) GetItems(ctx context.Context, empty *api.Empty) (*api.ItemList, error) {
	return nil, status.Errorf(codes.Unimplemented, "Not implemented")
}

func (service *whService) Hello(ctx context.Context, text *api.Text) (*api.Text, error) {
	txt := text.GetText()
	return &api.Text{Text: fmt.Sprintln("Hello there", txt)}, nil
}
