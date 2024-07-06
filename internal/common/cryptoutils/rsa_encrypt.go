package cryptoutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// EncryptRSA encrypt data with RSA public key
func EncryptRSA(data, key string) (string, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the public key")
	}

	rsaPublicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPublicKey, []byte(data))
	if err != nil {
		return "", err
	}

	encrypted := base64.StdEncoding.EncodeToString(ciphertext)
	return encrypted, nil
}
