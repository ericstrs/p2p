package peer

import (
	"fmt"
	"log"
	"net"
)

func Server(port string) {
	// Listen for connection
	addr := "localhost:" + port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Error: Could not listen on %s: %v\n", addr, err)
		return
	}
	fmt.Printf("Listening for connection on %s\n", addr)

	defer listener.Close()

	// Accpet connection
	conn, err := listener.Accept()
	if err != nil {
		log.Printf("Error: Could not accept connection from client: %v\n", err)
		return
	}
	fmt.Println("Accepted connection")

	defer conn.Close()

	// Communication can now take place
	fmt.Fprintln(conn, "Hello Charlie.")
	var msg string
	fmt.Fscanln(conn, &msg)
	fmt.Printf("Charlie says: %s\n", msg)
}

func Client(port string) {
	addr := "localhost:" + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Error: Could connect to %s: %v\n", addr, err)
		return
	}
	fmt.Printf("Connected to %s\n", addr)

	defer conn.Close()

	// Communcation can now take place
	fmt.Fprintln(conn, "Hello Sierra")
	var msg string
	fmt.Fscanln(conn, &msg)
	fmt.Printf("Sierra says: %s\n", msg)
}
