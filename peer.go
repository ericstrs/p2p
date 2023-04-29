package peer

import (
	"bufio"
	"fmt"
	"io"
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

	// Accept connection
	conn, err := listener.Accept()
	if err != nil {
		log.Printf("Error: Could not accept connection from client: %v\n", err)
		return
	}
	defer conn.Close()

	fmt.Printf("Accepted connection\n")

	// Prompt user for password
	var password string
	fmt.Print("Enter password: ")
	fmt.Scan(&password)
	fmt.Println()

	// Read messages
	go handleMessages(conn, []byte(password))

	// Get user input
	s := bufio.NewScanner(os.Stdin)
	// Send messages
	for {
		if !s.Scan() {
			break
		}
		m := s.Text()

		// Encrypt message
		ciphertext, err := Encrypt([]byte(m), []byte(password))
		if err != nil {
			log.Printf("%v\n", err)
			return
		}

		// Send encrypted message
		fmt.Fprintf(conn, "%s\n", ciphertext)
		fmt.Printf("Sent ciphertext: %q\n\n", ciphertext)
	}

	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func Client(port string) {

	// Connect to server
	addr := "localhost:" + port
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("Error: Could connect to %s: %v\n", addr, err)
		return
	}
	defer conn.Close()

	fmt.Printf("Connected to %s\n", addr)

	// Prompt user for password
	var password string
	fmt.Print("Enter password: ")
	fmt.Scan(&password)
	fmt.Println()

	// Read messages
	go handleMessages(conn, []byte(password))

	// Send messages
	s := bufio.NewScanner(os.Stdin)
	for {
		if !s.Scan() {
			break
		}
		m := s.Text()

		// Encrypt message
		ciphertext, err := Encrypt([]byte(m), []byte(password))
		if err != nil {
			log.Printf("%v\n", err)
			return
		}

		// Send encrypted message
		fmt.Fprintf(conn, "%s\n", ciphertext)
		fmt.Printf("Sent ciphertext: %q\n\n", ciphertext)
	}

	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func handleMessages(conn net.Conn, password []byte) {
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)

		if err != nil {
			if err == io.EOF {
				// Connection was closed by other party
				fmt.Println("Connection was closed")
				break
			}
			log.Fatalf("Error receiving ciphertext %v", err)
		}

		// Remove newline
		ciphertext := buffer[:n-1]

		fmt.Printf("Received ciphertext: %q\n", ciphertext)

		// Decrypt the ciphertext
		plaintext, err := Decrypt(ciphertext, []byte(password))
		if err != nil {
			log.Fatalf("Error decrypting ciphertext: %v", err)
		}

		fmt.Printf("Decrypted plaintext: %s\n\n", plaintext)
	}
}
