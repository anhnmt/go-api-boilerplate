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
func EncryptRSA(data, key []byte) (string, error) {
	block, _ := pem.Decode(key)
	if block == nil {
		return "", fmt.Errorf("failed to parse PEM block containing the public key")
	}

	rsaPublicKey, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return "", err
	}

	pkey, ok := rsaPublicKey.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("key is not a valid RSA public key")
	}

	ciphertext, err := rsa.EncryptPKCS1v15(rand.Reader, pkey, data)
	if err != nil {
		return "", err
	}

	encrypted := base64.StdEncoding.EncodeToString(ciphertext)
	return encrypted, nil
}
