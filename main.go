package main

import (
	"blockchain/peer"
	"fmt"
	"os"
)

func main() {

	if len(os.Args) != 3 {
		fmt.Println("There should be 3 commands: main [host] [port]")
	}

	host := os.Args[1]
	port := os.Args[2]

	fmt.Println("Starting...")
	peer := peer.NewBlockchainPeer(host, port, "Kien Mai, your fellow internet peer")

	peer.Run()

	for {
	}
}
