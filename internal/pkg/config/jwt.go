package config

type JWT struct {
	Secret         string `mapstructure:"secret" defaultvalue:"super-secret-9a443f5d-7bc8-4c90-bb50-2fd45050a74a"`
	TokenExpires   string `mapstructure:"token_expires" defaultvalue:"5m"`
	RefreshExpires string `mapstructure:"refresh_expires" defaultvalue:"24h"`
}
