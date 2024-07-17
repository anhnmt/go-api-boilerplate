package config

type JWT struct {
	Secret         string `validate:"required" mapstructure:"secret" defaultvalue:"this-is-super-secret-key"`
	TokenExpires   string `validate:"required" mapstructure:"token_expires" defaultvalue:"5m"`
	RefreshExpires string `validate:"required" mapstructure:"refresh_expires" defaultvalue:"24h"`
}
