package main

import (
	"log"
	"net"
)

func main() {

}

func searchService() {
	addr := "127.0.0.1:8050"
	go sendInvitations(addr)
	// todo setup an http server listening for requests
	// from stock center
	// todo: maybe once center has discovered us reduce the
	// rate of send invitations
}

// continuously broadcast invitation message over UDP
// with address to connect
func sendInvitations(myAddr string) {
	// todo: add timer and loop
	con, _ := net.Dial("udp", "127.0.0.1:3000")
	buf := []byte("Hello zerver:D\n")
	log.Println("Sending stuff :DDD")
	_, err := con.Write(buf)
	if err != nil {
		log.Println(err)
	}
}
