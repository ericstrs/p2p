package main

import (
	"log"
	"os"

	"github.com/ericstrs/p2p/peer"
)

func main() {
	if len(os.Args) != 3 {
		log.Println("Usage: ./p2p <role> <port>")
		return
	}

	role := os.Args[1]
	port := os.Args[2]

	if role == "server" {
		peer.Server(port)
	} else if role == "client" {
		peer.Client(port)
	} else {
		log.Println("Invalid role specified. Valid roles are either 'server' or 'client'.")
	}
}
