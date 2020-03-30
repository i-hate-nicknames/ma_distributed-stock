package main

import (
	"log"
	"os"

	wh "nvm.ga/mastersofcode/golang_2019/stock_distributed/internal/warehouse"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Please, provide port to listen on as an argument")
	}
	port := os.Args[1]

	wh.StartWarehouse(port)
}
