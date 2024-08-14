package main

import (
	"log"
	"os"
	"strings"

	"github.com/ericstrs/p2p/peer"
)

const usage = "Usage: p2p <role> [host] [port]"

func main() {
	if len(os.Args) < 2 {
		log.Println("Error: role is required")
		log.Println(usage)
		os.Exit(1)
	}

	role := os.Args[1]

	host := "localhost"
	port := "8080"

	if len(os.Args) > 2 {
		host = os.Args[2]
	}
	if len(os.Args) > 3 {
		port = os.Args[3]
	}

	switch strings.ToLower(role) {
	case "server":
		peer.Server(host, port)
	case "client":
		peer.Client(host, port)
	default:
		log.Println("Error: Invalid role specified. Valid roles are either 'server' or 'client'.")
		log.Println(usage)
		os.Exit(1)
	}
}
