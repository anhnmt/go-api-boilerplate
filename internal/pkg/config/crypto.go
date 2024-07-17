package config

import (
	"encoding/base64"
)

type Crypto struct {
	PublicKey  string `validate:"required" mapstructure:"public_key"`
	PrivateKey string `validate:"required" mapstructure:"private_key"`
}

func (c *Crypto) PublicKeyBytes() []byte {
	data, _ := base64.StdEncoding.DecodeString(c.PublicKey)
	return data
}

func (c *Crypto) PrivateKeyBytes() []byte {
	data, _ := base64.StdEncoding.DecodeString(c.PrivateKey)
	return data
}
