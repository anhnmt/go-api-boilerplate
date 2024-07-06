package cryptoutils

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
)

// DecryptAES decrypt data with AES key
func DecryptAES(data, key string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher([]byte(key[:32]))
	if err != nil {
		return "", err
	}

	// Create secret IV
	iv := [16]byte{}
	mode := cipher.NewCBCDecrypter(block, iv[:])
	decrypted := make([]byte, len(ciphertext))
	mode.CryptBlocks(decrypted, ciphertext)

	decrypted = PKCS7UnPadding(decrypted)
	return string(decrypted), nil
}
