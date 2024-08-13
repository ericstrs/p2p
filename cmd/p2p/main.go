package main

import (
	"log"
	"os"
	"strings"

	"github.com/ericstrs/p2p/peer"
)

func main() {
	if len(os.Args) != 3 {
		log.Println("Usage: p2p <role> <port>")
		return
	}

	role := os.Args[1]
	port := os.Args[2]

	switch strings.ToLower(role) {
	case "server":
		peer.Server(port)
	case "client":
		peer.Client(port)
	default:
		log.Fatal("Invalid role specified. Valid roles are either 'server' or 'client'.")
	}
}
