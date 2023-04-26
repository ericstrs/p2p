package peer

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"log"
	"net"
	"os"
)

func Server(port string) {

	// Prompt user for password
	var password string
	fmt.Print("Enter password: ")
	fmt.Scan(&password)

	// Generate random salt
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

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

	fmt.Println("Accepted connection")

	// Read messages
	go handleMessages(conn, []byte(password))

	// Send messages
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		m := s.Text()

		// Encrypt message
		ciphertext, err := Encrypt([]byte(m), []byte(password))
		if err != nil {
			log.Printf("%v\n", err)
			return
		}

		fmt.Printf("Sending encrypetd message: %q\n", ciphertext)

		// Send encrypted message
		fmt.Fprintf(conn, "%s\n", ciphertext)
	}

	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func Client(port string) {

	// Prompt user for password
	var password string
	fmt.Print("Enter password: ")
	fmt.Scan(&password)

	// Generate random salt
	salt := make([]byte, 32)
	_, err := rand.Read(salt)
	if err != nil {
		log.Printf("%v\n", err)
		return
	}

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
	go handleMessages(conn, []byte(password))

	// Send messages
	s := bufio.NewScanner(os.Stdin)
	for s.Scan() {
		m := s.Text()

		// Encrypt message
		ciphertext, err := Encrypt([]byte(m), []byte(password))
		if err != nil {
			log.Printf("%v\n", err)
			return
		}

		fmt.Printf("Sending encrypetd message: %q\n", ciphertext)

		// Send encrypted message
		fmt.Fprintf(conn, "%s\n", ciphertext)
	}

	if err := s.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
	}
}

func handleMessages(conn net.Conn, password []byte) {
	s := bufio.NewScanner(conn)
	for s.Scan() {
		c := s.Text()
		fmt.Printf("Received encrypetd message: %q\n", c)
		// Decrypt the message
		plaintext, err := Decrypt([]byte(c), []byte(password))

		if err != nil {
			panic(err)
		}

		fmt.Printf("%s\n", plaintext)
	}

	if err := s.Err(); err != nil {
		log.Printf("Error reading: %v\n", err)
	}
}
