package peer

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"log"
	"testing"
)

func TestGenerateKey(t *testing.T) {
	password := []byte("Charlie")

	// Generate a random salt
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		t.Fatal(err)
	}

	// Create keys
	aliceKey := generateKey(password, salt)
	bobKey := generateKey(password, salt)

	// Check to see if keys are the same
	if !bytes.Equal(aliceKey, bobKey) {
		log.Println("Alice and Bob have different keys")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	plaintext := []byte("secret message")

	// Encrypt the plaintext
	encodedText, err := Encrypt(plaintext, []byte("password"))

	if err != nil {
		t.Fatalf("Encrypt failed: %v", err)
	}

	// Decrypt the encoded text
	decodedText, err := Decrypt(encodedText, []byte("password"))

	if err != nil {
		t.Fatalf("Decrypt failed: %v", err)
	}

	// Check if decoded text matches plaintext
	if !bytes.Equal(plaintext, decodedText) {
		fmt.Printf("Got: %s\nExpected: %d\n", decodedText, plaintext)
		t.Fatalf("Decoded text doesn't match plaintext")
	}
}

func TestPaddingUnpadding(t *testing.T) {
	plaintext := []byte("secret message")
	blockSize := 16

	// Test padding
	padded := pkcs7Padding(plaintext, blockSize)
	expectedPadding := 16 - len(plaintext)%blockSize
	expectedPaddingLen := len(plaintext) + expectedPadding

	if len(padded) != expectedPaddingLen {
		t.Errorf("padding lenght was: %d\nexpected: %d", len(padded), expectedPaddingLen)
	}
	lastByte := padded[len(padded)-1]
	for i := 1; i <= int(lastByte); i++ {
		if padded[len(padded)-i] != lastByte {
			t.Errorf("padding byte %d was %d, expected %d", i, padded[len(padded)-i], lastByte)
		}
	}

	// Test unpadded
	unpadded, err := pkcs7Unpadding(padded)
	if err != nil {
		t.Errorf("unexpected error during unpadding: %v", err)
	}

	if !bytes.Equal(unpadded, plaintext) {
		t.Errorf("unpadded data was %q, expected %q", unpadded, plaintext)
	}
}
