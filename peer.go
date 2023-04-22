package peer

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func Server(port string) {
	// Listen for connection
	addr := "localhost:" + port
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Printf("Error: Could not listen on %s: %v\n", addr, err)
		return
	}
	defer listener.Close()

	fmt.Printf("Listening for connection on %s\n", addr)

	// Accpet connection
	conn, err := listener.Accept()
	if err != nil {
		log.Printf("Error: Could not accept connection from client: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Println("Accepted connection")

	// Read messages
	go handleMessages(conn)

	// Send messages
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		m := s.Text()
		fmt.Fprintf(conn, "Alice: %s\n", m)
	}

	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func Client(port string) {
	// Connect
	addr := "localhost:" + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Error: Could connect to %s: %v\n", addr, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", addr)

	// Read messages
	go handleMessages(conn)

	// Send messages
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		m := s.Text()
		fmt.Fprintf(conn, "Bob: %s\n", m)
	}

	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func handleMessages(conn net.Conn) {
	s := bufio.NewScanner(conn)
	for s.Scan() {
		m := s.Text()
		fmt.Printf("%s\n", m)
	}

	if err := s.Err(); err != nil {
		log.Printf("Error reading: %v\n", err)
	}
}
