package logger

type Config struct {
	ServiceName string `validate:"required" mapstructure:"service_name"`
	Format      string `validate:"required,oneof=text json" mapstructure:"format" defaultvalue:"text"`
	Level       string `validate:"required,oneof=debug info warn error panic fatal" mapstructure:"level" defaultvalue:"info"`

	// Config file
	File       string `mapstructure:"file"`
	MaxSize    int    `validate:"required_with=File,number" mapstructure:"max_size" defaultvalue:"100"` // MB
	MaxBackups int    `validate:"required_with=File,number" mapstructure:"max_backups" defaultvalue:"5"`
	MaxAge     int    `validate:"required_with=File,number" mapstructure:"max_age" defaultvalue:"28"` // days
}
