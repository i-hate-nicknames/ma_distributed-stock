package main

import (
	"log"
	"net"
	"os"
	"time"
)

const invitationTimeout = 500 * time.Millisecond

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Please, provide port to listen on as an argument")
	}
	port := os.Args[1]
	addr := "127.0.0.1:" + port
	go sendInvitations(addr)

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
