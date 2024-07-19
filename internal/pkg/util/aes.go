package util

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

// DecryptAES decrypt data with AES privateKey
func DecryptAES(data, key string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher([]byte(key[:32]))
	if err != nil {
		return nil, err
	}

	// Create secret IV
	iv := [16]byte{}
	mode := cipher.NewCBCDecrypter(block, iv[:])
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	if len(decrypted) == 0 {
		return nil, fmt.Errorf("decrypted data is empty")
	}

	decrypted = PKCS7UnPadding(decrypted)
	return decrypted, nil
}

// EncryptAES encrypt data with AES privateKey
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
