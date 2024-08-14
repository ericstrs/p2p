# Point-to-Point (P2P) Encrypted Messaging

This project implements a simple Point-to-Point encrypted messaging system using AES encryption. The system includes both server and client. It allows two parties to establish a secure session to send messages to one another. The communication between both parties is encrypted with a password-derived key.

See the [report](./report.md) for a detailed breakdown of how this system works.

## Features

* Point-to-Point communication over TCP.
* AES encryption using CBC mode.
* Password-based key generation with PBKDF2 and SHA-256.
* PKCS#7 padding.

## Usage

### Server

Starting the server:

```
go run main.go server [host] [port]
```

Replace `[host]` and `[port]` with desired host and port.

### Client

Starting the client:

```
go run main.go client [host] [port]
```

Replace `[host]` and `[port]` with desired host and port.

### Password

Both the server and client will prompt for a password. Enter the same password on both sides to establish a secure connection.

## Design considerations

This tool makes the assumption that Alice and Bob share a secret (password), and are therefore able to utilize the benefits of symmetric key cryptography for the Internet transmission.
