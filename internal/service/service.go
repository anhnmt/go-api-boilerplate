package service

import (
	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/gen/pb"
)

type Service interface {
	Register(grpcSrv grpc.ServiceRegistrar)
}

type service struct {
	UserSvc pb.UserServiceServer
	AuthSvc pb.AuthServiceServer
}

func New(
	userSvc pb.UserServiceServer,
	authSvc pb.AuthServiceServer,
) Service {
	return &service{
		UserSvc: userSvc,
		AuthSvc: authSvc,
	}
}

func (s *service) Register(grpcSrv grpc.ServiceRegistrar) {
	pb.RegisterUserServiceServer(grpcSrv, s.UserSvc)
	pb.RegisterAuthServiceServer(grpcSrv, s.AuthSvc)
}
