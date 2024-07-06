package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// EncryptAES encrypt data with AES key
func EncryptAES(data, key string) (string, error) {
	plaintext := PKCS7Padding([]byte(data))
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
