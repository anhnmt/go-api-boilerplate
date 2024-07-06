// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package grpc

import (
	"github.com/anhnmt/go-api-boilerplate/internal/infrastructure/gormgen"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	"github.com/anhnmt/go-api-boilerplate/internal/server/interceptor/auth"
	"github.com/anhnmt/go-api-boilerplate/internal/service"
	"github.com/anhnmt/go-api-boilerplate/internal/service/auth/business"
	"github.com/anhnmt/go-api-boilerplate/internal/service/auth/repository/redis"
	"github.com/anhnmt/go-api-boilerplate/internal/service/auth/transport/grpc"
	"github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/command"
	"github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/query"
	"github.com/anhnmt/go-api-boilerplate/internal/service/user/business"
	"github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/command"
	"github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
	"github.com/anhnmt/go-api-boilerplate/internal/service/user/transport/grpc"
	"github.com/redis/go-redis/v9"
	"google.golang.org/grpc"
)

// Injectors from wire.go:

func New(gormQuery *gormgen.Query, rdb redis.UniversalClient, cfgGrpc config.Grpc, cfgJWT config.JWT) (*grpc.Server, error) {
	query := userquery.New(gormQuery)
	command := sessioncommand.New(gormQuery)
	sessionqueryQuery := sessionquery.New(gormQuery)
	authredisRedis := authredis.New(rdb)
	business := authbusiness.New(cfgJWT, query, command, sessionqueryQuery, authredisRedis)
	authInterceptor := authinterceptor.New(business)
	usercommandCommand := usercommand.New(gormQuery)
	userbusinessBusiness := userbusiness.New(usercommandCommand, query)
	userServiceServer := usergrpc.New(userbusinessBusiness)
	authServiceServer := authgrpc.New(business)
	serviceService := service.New(userServiceServer, authServiceServer)
	server := initServer(cfgGrpc, authInterceptor, serviceService)
	return server, nil
}
