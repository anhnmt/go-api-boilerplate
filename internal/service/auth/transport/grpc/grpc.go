package authgrpc

import (
	"context"

	"google.golang.org/grpc"

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

func (s *grpcService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	return s.authBusiness.Login(ctx, req)
}

func (s *grpcService) Info(ctx context.Context, _ *pb.InfoRequest) (*pb.InfoReply, error) {
	return s.authBusiness.Info(ctx)
}
