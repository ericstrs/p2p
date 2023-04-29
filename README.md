# Secure Instant Point-to-Point (P2P) Messaging

This project is a secure instant messaging tool written in Go. The system supports the following functions

* Alice and Bob can use the tool to send messages to each other.
* Alice and Bob share the same password, they must use the password to set up the tool to correctly encrypt and decrypt messages shared between each other.
* Each message during Internet transmission must be encrypted using a key no less than 56 bits.

## Design considerations

This makes the assumption that Alice and Bob share a secret (password), and are therefor able to utilize the benefits of symmetric key cryptography for the Internet transmission.
