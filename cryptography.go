package peer

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

func generateKey(password []byte, salt []byte) []byte {
	key := pbkdf2.Key(password, salt, 10000, 32, sha256.New)

	return key
}

func Encrypt(plaintext []byte, password []byte) (string, error) {
	// Generate new salt
	salt := make([]byte, 16)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	// Generate new key
	key := generateKey(password, salt)
	fmt.Printf("Enc's key: %x\n", key)

	// Check key size
	if len(key) != 32 {
		return "", errors.New("key size if ont 256 bits")
	}

	// Create a AES cipher block using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// Generate a new IV
	iv := make([]byte, 16+len(salt))
	_, err = io.ReadFull(rand.Reader, iv[:16])
	if err != nil {
		return "", err
	}
	copy(iv[16:], salt)

	// Encrypt the plaintext
	mode := cipher.NewCBCEncrypter(block, iv[:16])
	paddedPlaintext := pkcs7Padding(plaintext, len(key))

	fmt.Printf("\n")
	fmt.Printf("Padded plaintext size: %d\n", len(paddedPlaintext))
	ciphertext := make([]byte, len(paddedPlaintext))
	mode.CryptBlocks(ciphertext, paddedPlaintext)
	fmt.Printf("Encrypted text size: %d\n", len(ciphertext))
	fmt.Printf("Padding size:: %d\n", len(paddedPlaintext)-len(plaintext))
	fmt.Printf("\n")

	saltIvCipher := make([]byte, len(iv)+len(ciphertext))
	copy(saltIvCipher, iv)
	copy(saltIvCipher[len(iv):], ciphertext)

	// Encode the ciphertext
	encodedText := base64.StdEncoding.EncodeToString(saltIvCipher)

	return encodedText, nil
}

func Decrypt(encodedText string, password []byte) ([]byte, error) {
	saltIvCipher, err := base64.StdEncoding.DecodeString(encodedText)
	if err != nil {
		return nil, err
	}

	//fmt.Println("Start of decrypt:", saltIvCipher)
	salt := saltIvCipher[16:32]
	//fmt.Println("salt:", salt)
	key := generateKey(password, salt)
	fmt.Printf("Dec's key: %x\n", key)

	// Check key size
	if len(key) != 32 {
		return nil, errors.New("key size is not 256 bits")
	}

	// Check length of encrypted message
	if len(saltIvCipher) < aes.BlockSize {
		return nil, errors.New("encrypted text length is invalid")
	}

	// Extract the IV from the ciphertext
	iv := saltIvCipher[:16]
	encryptedText := saltIvCipher[16:]

	// Create AES cipher block using the key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Decrypt the ciphertext using AES CBC mode
	mode := cipher.NewCBCDecrypter(block, iv)
	decryptedText := make([]byte, len(encryptedText))
	mode.CryptBlocks(decryptedText, encryptedText)

	// Remove the padding
	unpaddedText, err := pkcs7Unpadding(decryptedText)
	if err != nil {
		return nil, err
	}

	return unpaddedText, nil
}

/*
func pkcs7Padding(plaintext []byte, blockSize int) []byte {
	padding := blockSize - len(plaintext)%blockSize
	padText := make([]byte, len(plaintext)+padding)
	copy(padText, plaintext)
	for i := len(plaintext); i < len(padText); i++ {
		padText[i] = byte(padding)
	}
	fmt.Println("paddedText siz:", len(padText))
	fmt.Println("last byte of padded text", padText[len(padText)-1])

	return padText
}*/

func pkcs7Padding(plaintext []byte, paddingSize int) []byte {
	padSize := paddingSize - (len(plaintext) % paddingSize)
	padding := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(plaintext, padding...)
}

func pkcs7Unpadding(plaintext []byte) ([]byte, error) {
	paddingSize := int(plaintext[len(plaintext)-1])
	if paddingSize > len(plaintext) {
		fmt.Printf("paddingsize= %d, length %d\n", paddingSize, len(plaintext))
		return nil, errors.New("invalid padding")
	}

	return plaintext[:len(plaintext)-paddingSize], nil
}

/*
func pkcs7Unpadding(plaintext []byte) ([]byte, error) {
	length := len(plaintext)
	unpadding := int(plaintext[length-1])
	if unpadding > length {
		fmt.Printf("unpadding = %d, length %d\n", unpadding, length)
		return nil, errors.New("invalid padding")
	}
	return plaintext[:(length - unpadding)], nil
}
*/
