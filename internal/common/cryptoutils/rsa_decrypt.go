package cryptoutils

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

func DecryptRSA(data string, key []byte) ([]byte, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return nil, fmt.Errorf("failed to parse PEM block containing the private key")
	}

	rsaPrivateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	pkey, ok := rsaPrivateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("key is not a valid RSA private key")
	}

	ciphertext, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return nil, err
	}

	rawData, err := rsa.DecryptPKCS1v15(rand.Reader, pkey, ciphertext)
	if err != nil {
		return nil, err
	}

	return rawData, nil
}

func DecryptRSAString(data string, key []byte) (string, error) {
	rawData, err := DecryptRSA(data, key)
	if err != nil {
		return "", err
	}

	return string(rawData), nil
}
