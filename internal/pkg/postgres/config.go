package postgres

import (
	"fmt"
	"net/url"
)

const dbType = "postgres"

type Config struct {
	Migrate bool `validate:"boolean" mapstructure:"migrate" defaultvalue:"true"`
	Debug   bool `validate:"boolean" mapstructure:"debug" defaultvalue:"false"`
	Otel    bool `validate:"boolean" mapstructure:"otel" defaultvalue:"false"`
	Writer  Base `validate:"required" mapstructure:"writer"`
	Reader  Base `mapstructure:"reader"`
}

type ParseDSN interface {
	ParseDSN() url.URL
}

var _ ParseDSN = (*Base)(nil)

type Base struct {
	Enable          bool   `validate:"boolean" mapstructure:"enable"`
	Host            string `validate:"required" mapstructure:"host" defaultvalue:"localhost"`
	Port            int    `validate:"required" mapstructure:"port" defaultvalue:"5432"`
	User            string `validate:"required" mapstructure:"user" defaultvalue:"postgres"`
	Password        string `validate:"required" mapstructure:"password" defaultvalue:"postgres"`
	Database        string `validate:"required" mapstructure:"database" defaultvalue:"postgres"`
	SSLMode         string `mapstructure:"ssl_mode" defaultvalue:"disable"`
	ApplicationName string `mapstructure:"application_name"`

	MaxIdleConns int `mapstructure:"max_idle_conns" defaultvalue:"5"`
	MaxOpenConns int `mapstructure:"max_open_conns" defaultvalue:"10"`

	MaxConnIdleTime string `mapstructure:"max_conn_idle_time" defaultvalue:"5m"`
	MaxConnLifetime string `mapstructure:"max_conn_lifetime" defaultvalue:"15m"`
}

func (p Base) ParseDSN() url.URL {
	dsn := url.URL{
		Scheme: dbType,
		User:   url.UserPassword(p.User, p.Password),
		Host:   fmt.Sprintf("%s:%d", p.Host, p.Port),
		Path:   p.Database,
	}

	q := dsn.Query()
	q.Add("sslmode", p.SSLMode)

	if p.ApplicationName != "" {
		q.Add("application_name", p.ApplicationName)
	}

	return dsn
}
