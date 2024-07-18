package service

import (
	"go.uber.org/fx"

	authservice "github.com/anhnmt/go-api-boilerplate/internal/service/auth"
	sessionservice "github.com/anhnmt/go-api-boilerplate/internal/service/session"
	userservice "github.com/anhnmt/go-api-boilerplate/internal/service/user"
)

// Module provided to fx
var Module = fx.Module("service",
	authservice.Module,
	sessionservice.Module,
	userservice.Module,
)
