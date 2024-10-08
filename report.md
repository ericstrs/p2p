# Secure Instant Point-to-Point Messaging

This secure instant point-to-point (P2P) messaging tool is suited for two parties where it is assumed they share a password. The password is used the tool to correctly encrypt and decrypted messages shared between them. The assumption of a shared secret (password) between the two parties makes it so that the benefits of symmetric key cryptography can be utilized. In this case, each message during Internet transmission is encrypted using the advanced encryption standard (AES) with a 256-bit key. Taking a broad overview, the tool works as follows: establish a connection between both parties, prompt for a message, generate a new key to encrypt the message, transmit the message, and decrypt the message upon reception.

The main function first checks if the user inputs a valid role. The roles consist of server and client. The corresponding code goes as follows:

```
func main() {
  // omitted error handling
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
    // omitted error handling
  }
}
```

Functions `peer.Server` and `peer.Client` are reliant upon on a goroutine to concurrently run the `handleMessages` function to keep receiving messages from the client:

```
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

    // Decrypt the ciphertext
    plaintext, err := Decrypt(ciphertext, []byte(password))
    if err != nil {
      log.Fatalf("Error decrypting ciphertext: %v", err)
    }

    fmt.Printf("Decrypted plaintext: %s\n\n", plaintext)
  }
}
```

The `handleMessages` function takes in a connection and password as input. Waits and reads an encrypted message from the connection. The ciphertext and the password are then passed to the `Decrypt` function which returns the plaintext.

The `handleMessages` function waits until the connection receives an incoming ciphertext. Using the ciphertext and the password, a call is made to the `Decrypt` function to get the respective plaintext.

The `peer.Server` function handles the server role:

```
func Server(host, port string) {
  // Listen for connection
  addr := host + ":" + port
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

  // Prompt user for password
  var password string
  fmt.Print("Enter password: ")
  fmt.Scan(&password)

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
  }

  if err := s.Err(); err != nil {
    fmt.Fprintln(os.Stderr, "Error:", err)
  }
}
```

Firstly, `Server` listens on the `port`. It then accepts a connection from the client. Once this connection has been established, the user is prompted for the shared password. A call to the goroutine `handleMessages` function is made to handle reading messages while the main thread can focus on prompting and sending messages to the client. The `handleMessages` function continually prompts the user for a message, encrypts it, and transmits it to the other party.

The `peer.Client` function handles the client role:

```
func Client(host, port string) {
  // Connect to server
  addr := host + ":" + port
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
```

A connection is first made to the server. The user is prompted for the shared password. The goroutine `handleMessages` function is once again used read incoming messages from the other party. The user is then continually prompted for a message, encrypts it, and then transmits the message to the other party.

We now turn to the cryptographic part of this tool. Since we assume a shared secret between the two parties, we are able to generate the same private key for both parties. This is accomplished using a key derivation function (KDF) that takes as input the passphrase and produces a key. The KDF used in this tool is the Password-Based Key Derivation Function 2 (PBKKDF2). The PBKKDF2 takes in a passphrase and a salt to derive a secure encryption key. The `generateKey` function is shown below:

```
// generateKey returns a 256-bit key using the given password and salt.
func generateKey(pass []byte, salt []byte) []byte {
  return pbkdf2.Key(pass, salt, 10000, 32, sha256.New)
}
```

This uses the `sha256.New` function to create an instance of the SHA-256 hash function.

The `Encrypt` functions goes as follows:

```
func Encrypt(plaintext []byte, password []byte) ([]byte, error) {
  // Generate new salt
  salt := make([]byte, 16)
  if _, err := io.ReadFull(rand.Reader, salt); err != nil {
    return nil, err
  }

  // Generate new key
  key := generateKey(password, salt)

  // Check key size
  if len(key) != 32 {
    return nil, errors.New("key size is 256 bits")
  }

  // Create a AES cipher block using the key
  block, err := aes.NewCipher(key)
  if err != nil {
    return nil, err
  }

  // Generate a new IV
  iv := make([]byte, 16)
  _, err = io.ReadFull(rand.Reader, iv)
  if err != nil {
    return nil, err
  }

  // Encrypt the plaintext
  mode := cipher.NewCBCEncrypter(block, iv)
  paddedPlaintext := pkcs7Padding(plaintext, aes.BlockSize)
  ciphertext := make([]byte, len(paddedPlaintext))
  mode.CryptBlocks(ciphertext, paddedPlaintext)

  output := make([]byte, 0, len(salt)+len(iv)+len(ciphertext))
  output = append(output, salt...)
  output = append(output, iv...)
  output = append(output, ciphertext...)

  return output, nil
}
```

The `encrypt` function creates a salt to generate a new key. It should be noted that this design results in each new message using a newly generated key. Updating the key for every new messages means that, in the case of a key being compromised, the rest of the messages will be made unavailable to the attacker. The fact that the users are able to build the same key--that is, they don't actually have to transmit the key--makes generating a key for every message feasible. Next, a AES block cipher is created using the key, followed by the creation of an Initialization Vector (IV). The IV ensures that if the same message is sent repeatedly, the ciphertext will differ for each message. This eliminates potential patterns that would be visible to a keen attacker. Cipher Block Chaining (CBC) mode of operation is desirable in this use case as it establishes confidentially and message integrity.

PKCS #7 is used to pad the plaintext to a multiple of 16 bytes. PKCS7 padding works as follows: if message is not a multiple of 16 bytes long, count the number $n$ of bytes it would take to pad to a multiple of 16; append $n$ bytes. If the number of padding bytes is 0, then we don't need to do anything since it is already a multiple of 16. The corresponding code is shown below:

```
// pkcs7Padding returns the input plaintext padded with PKCS#7 padding.
func pkcs7Padding(plaintext []byte, paddingSize int) []byte {
  // Add padding to the end of the plaintext to make its length a
  // multiple of the block size.
  padSize := paddingSize - (len(plaintext) % paddingSize)

  // Create a list of length `paddingSize` that consists of values all
  // equal to `paddingSize`.
  padding := bytes.Repeat([]byte{byte(padSize)}, padSize)
  return append(plaintext, padding...)
}
```

Once the plaintext has been appropriately padded, we encrypt a series of blocks of data using CBC mode. Then, we prepend the salt and the IV to the ciphertext before returning.

The `Decrypt` function works in a similar fashion:

```
func Decrypt(ciphertext []byte, password []byte) ([]byte, error) {
  if len(ciphertext) < 32 {
    return nil, errors.New("ciphertext too short")
  }

  // Extract the salt and IV from ciphertext
  s := ciphertext[:16]
  iv := ciphertext[16:32]
  ciphertext = ciphertext[32:]

  // Generate key from password and salt
  key := generateKey(password, s)

  // Check key size
  if len(key) != 32 {
    return nil, errors.New("key size is not 256 bits")
  }

  // Create AES cipher block using the key
  block, err := aes.NewCipher(key)
  if err != nil {
    return nil, err
  }

  // Decrypt the ciphertext using AES CBC mode
  mode := cipher.NewCBCDecrypter(block, iv)
  plaintext := make([]byte, len(ciphertext))
  mode.CryptBlocks(plaintext, ciphertext)

  // Remove the padding
  unpaddedPlaintext, err := pkcs7Unpadding(plaintext)
  if err != nil {
    return nil, err
  }

  return unpaddedPlaintext, nil
}
```

Here, we first extract the salt and IV from the encrypted message which leaves us with just the ciphertext. Like in the `Encrypt` function, we generate the key using the shared password and salt. We create an AES cipher block using the key and specify the mode to be CBC. `CryptBlocks` is applied to decrypt the ciphertext accordingly. At this point, we are just left with the plaintext and some padded information. We are able to get rid of the extra data using the `pkcs7Unpadding` function. We get the last byte which represents the number of bytes that needed to be padded. We use this value to return a slice that removes the padding.

```
// pkcs7Unpadding returns a byte slice that leaves out the PKCS#7
// padding from the given plaintext.
func pkcs7Unpadding(plaintext []byte) ([]byte, error) {
  paddingSize := int(plaintext[len(plaintext)-1]) // bytes that were added
  if paddingSize > len(plaintext) {
    fmt.Printf("paddingsize= %d, length %d\n", paddingSize, len(plaintext))
    return nil, errors.New("invalid padding")
  }
  return plaintext[:len(plaintext)-paddingSize], nil
}
```

Once we have successfully unpadded we are left with the original plaintext. The plaintext is returned and printed to the user.
