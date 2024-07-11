package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// EncryptAES encrypt data with AES key
func EncryptAES(plaintext []byte, key string) (string, error) {
	if len(plaintext) == 0 {
		return "", fmt.Errorf("plaintext is empty")
	}

	plaintext = PKCS7Padding(plaintext)
	ciphertext := make([]byte, len(plaintext))

	block, err := aes.NewCipher([]byte(key[:32]))
	if err != nil {
		return "", err
	}

	// Create secret IV
	iv := [16]byte{}
	mode := cipher.NewCBCEncrypter(block, iv[:])
	mode.CryptBlocks(ciphertext, plaintext)

	encrypted := base64.StdEncoding.EncodeToString(ciphertext)
	return encrypted, nil
}
