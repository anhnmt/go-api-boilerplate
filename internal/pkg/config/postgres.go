package config

import (
	"fmt"
	"net/url"
)

const dbType = "postgres"

var _ ParseDSN = (*PostgresBase)(nil)

type ParseDSN interface {
	ParseDSN() url.URL
}

type Postgres struct {
	Migrate bool         `mapstructure:"migrate" defaultvalue:"true"`
	Debug   bool         `mapstructure:"debug" defaultvalue:"false"`
	Writer  PostgresBase `mapstructure:"writer"`
	Reader  PostgresBase `mapstructure:"reader"`
}

type PostgresBase struct {
	Enable          bool   `mapstructure:"enable"`
	Host            string `mapstructure:"host" defaultvalue:"localhost"`
	Port            int    `mapstructure:"port" defaultvalue:"5432"`
	User            string `mapstructure:"user" defaultvalue:"postgres"`
	Password        string `mapstructure:"password" defaultvalue:"postgres"`
	Database        string `mapstructure:"database" defaultvalue:"postgres"`
	SSLMode         string `mapstructure:"ssl_mode" defaultvalue:"disable"`
	ApplicationName string `mapstructure:"application_name"`

	MaxIdleConns int `mapstructure:"max_idle_conns" defaultvalue:"5"`
	MaxOpenConns int `mapstructure:"max_open_conns" defaultvalue:"10"`

	MaxConnIdleTime string `mapstructure:"max_conn_idle_time" defaultvalue:"5m"`
	MaxConnLifetime string `mapstructure:"max_conn_lifetime" defaultvalue:"15m"`
}

func (p PostgresBase) ParseDSN() url.URL {
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
