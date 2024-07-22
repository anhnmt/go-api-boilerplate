package config

type App struct {
	Name    string `validate:"required" mapstructure:"name" defaultvalue:"api-server"`
	Version string `mapstructure:"version"`
}
