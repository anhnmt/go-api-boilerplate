package cryptoutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func DecryptRSA(data string, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private privateKey")
	}

	rsaPrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	rawData, err := rsa.DecryptPKCS1v15(rand.Reader, rsaPrivateKey, ciphertext)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func DecryptRSAString(data string, privateKey []byte) (string, error) {
	rawData, err := DecryptRSA(data, privateKey)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}
