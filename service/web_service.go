package main

import (
	"log"
	"net"
)

// RemoteMachine represents remote stock machine
type RemoteMachine struct {
	address string
	conn    net.Conn
	items   []int
}

func main() {
	// todo: add web service handlers here
	discoverStocks()
}

func discoverStocks() {
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
		log.Println("Reding next portion")
		_, _, err := conn.ReadFrom(buf[:])
		if err != nil {
			log.Println(err)
		}
		addMachine(buf)
	}
}

func addMachine(buf []byte) {
	// todo: actually add machine
	log.Println("client send us: " + string(buf))
}
