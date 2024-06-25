package authgrpc

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/auth/business"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type grpcService struct {
	pb.UnimplementedAuthServiceServer

	authBusiness *authbusiness.Business
}

func New(
	grpcSrv *grpc.Server,
	authBusiness *authbusiness.Business,
) pb.AuthServiceServer {
	svc := &grpcService{
		authBusiness: authBusiness,
	}

	pb.RegisterAuthServiceServer(grpcSrv, svc)

	return svc
}

func (s *grpcService) Login(context.Context, *pb.LoginRequest) (*pb.LoginReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
