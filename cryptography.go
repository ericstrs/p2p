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

func generateKey(password []byte, salt []byte) []byte {
	key := pbkdf2.Key(password, salt, 10000, 32, sha256.New)

	return key
}

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

func Decrypt(ciphertext []byte, password []byte) ([]byte, error) {
	// Extract the salt and IV from ciphertext
	salt := ciphertext[:16]
	iv := ciphertext[16:32]
	ciphertext = ciphertext[32:]

	// Generate key from password and salt
	key := generateKey(password, salt)

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

func pkcs7Padding(plaintext []byte, paddingSize int) []byte {
	padSize := paddingSize - (len(plaintext) % paddingSize)
	padding := bytes.Repeat([]byte{byte(padSize)}, padSize)
	return append(plaintext, padding...)
}

/*
func pkcs7Padding(plaintext []byte, blockSize int) []byte {
	r := len(plaintext) % blockSize
	pl := blockSize - r
	for i := 0; i < pl; i++ {
		plaintext = append(plaintext, byte(pl))
	}
	return plaintext
}

func pkcs7Unpadding(plaintext []byte) ([]byte, error) {
	if plaintext == nil || len(plaintext) == 0 {
		return nil, nil
	}

	pl := int(plaintext[len(plaintext)-1])

	err := checkPadding(plaintext, pl)
	if err != nil {
		return nil, err
	}

	return plaintext[:len(plaintext)-pl], nil
}

func checkPadding(plaintext []byte, paddingLen int) error {
	if len(plaintext) < paddingLen {
		return errors.New("invalid padding length of plaintext smaller than padding length")
	}
	p := plaintext[len(plaintext)-paddingLen:]
	for _, pc := range p {
		if uint(pc) != uint(len(p)) {
			fmt.Println("trouble number:", pc)
			return errors.New("invalid padding one of the padding values does not represent the padded value")
		}
	}
	return nil
}
*/

func pkcs7Unpadding(plaintext []byte) ([]byte, error) {
	paddingSize := int(plaintext[len(plaintext)-1])
	if paddingSize > len(plaintext) {
		fmt.Printf("paddingsize= %d, length %d\n", paddingSize, len(plaintext))
		return nil, errors.New("invalid padding")
	}

	return plaintext[:len(plaintext)-paddingSize], nil
}
