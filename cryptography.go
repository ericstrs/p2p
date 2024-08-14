package peer

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

// generateKey returns a 256-bit key using the given password and salt.
func generateKey(pass []byte, salt []byte) []byte {
	return pbkdf2.Key(pass, salt, 10000, 32, sha256.New)
}

// Encrypt returns ciphertext for a given body of plaintext.
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
		return nil, errors.New("key size is not 256 bits")
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

// Decrypt returns plaintext for a given body of ciphertext.
func Decrypt(ciphertext []byte, password []byte) ([]byte, error) {
	if len(ciphertext) < 32 {
		fmt.Printf("Ciphertext in question: %q\n", ciphertext)
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
