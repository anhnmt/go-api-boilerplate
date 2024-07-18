package sessionservice

import (
	"go.uber.org/fx"

	sessioncommand "github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/command"
	sessionquery "github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/query"
)

// Module provided to fx
var Module = fx.Module("session_service", fx.Provide(
	sessioncommand.New,
	sessionquery.New,
))
