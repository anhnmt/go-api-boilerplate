package config

type JWT struct {
	Secret         string `mapstructure:"secret" defaultvalue:"this-is-super-secret-key"`
	TokenExpires   string `mapstructure:"token_expires" defaultvalue:"5m"`
	RefreshExpires string `mapstructure:"refresh_expires" defaultvalue:"24h"`
}
