package service

import (
	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type Service interface {
	Register(grpcSrv *grpc.Server)
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

func (s *service) Register(grpcSrv *grpc.Server) {
	pb.RegisterUserServiceServer(grpcSrv, s.UserSvc)
	pb.RegisterAuthServiceServer(grpcSrv, s.AuthSvc)
}
