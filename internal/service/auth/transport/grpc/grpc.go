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

func (s *grpcService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.authBusiness.Login(ctx, req)
}

func (s *grpcService) Info(ctx context.Context, _ *pb.InfoRequest) (*pb.InfoResponse, error) {
	return s.authBusiness.Info(ctx)
}

func (s *grpcService) RevokeToken(context.Context, *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RevokeToken not implemented")
}

func (s *grpcService) RefreshToken(context.Context, *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RefreshToken not implemented")
}
