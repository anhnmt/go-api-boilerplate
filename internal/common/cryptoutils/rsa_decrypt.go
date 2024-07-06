package cryptoutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func DecryptRSA(data, key string) (string, error) {
	block, _ := pem.Decode([]byte(key))
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the private key")
	}

	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return "", err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	rawData, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, ciphertext)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}
