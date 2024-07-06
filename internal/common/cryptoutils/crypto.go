package cryptoutils

import (
	"bytes"
	"crypto/aes"
)

// PKCS7Padding padding data with PKCS7
func PKCS7Padding(ciphertext []byte) []byte {
	padding := aes.BlockSize - len(ciphertext)%aes.BlockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padText...)
}

// PKCS7UnPadding unPadding data with PKCS7
func PKCS7UnPadding(plantText []byte) []byte {
	length := len(plantText)
	unPadding := int(plantText[length-1])
	return plantText[:(length - unPadding)]
}
