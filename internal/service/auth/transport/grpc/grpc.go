package authgrpc

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/gen/pb"
	authbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/auth/business"
)

type grpcService struct {
	pb.UnimplementedAuthServiceServer

	authBusiness *authbusiness.Business
}

type Params struct {
	fx.In

	GrpcSever    *grpc.Server
	AuthBusiness *authbusiness.Business
}

func New(p Params) pb.AuthServiceServer {
	svc := &grpcService{
		authBusiness: p.AuthBusiness,
	}

	pb.RegisterAuthServiceServer(p.GrpcSever, svc)
	return svc
}

func (s *grpcService) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	return s.authBusiness.Login(ctx, req)
}

func (s *grpcService) Info(ctx context.Context, _ *pb.InfoRequest) (*pb.InfoResponse, error) {
	return s.authBusiness.Info(ctx)
}

func (s *grpcService) RevokeToken(ctx context.Context, _ *pb.RevokeTokenRequest) (*pb.RevokeTokenResponse, error) {
	return nil, s.authBusiness.RevokeToken(ctx)
}

func (s *grpcService) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	return s.authBusiness.RefreshToken(ctx, req)
}

func (s *grpcService) ActiveSessions(ctx context.Context, req *pb.ActiveSessionsRequest) (*pb.ActiveSessionsResponse, error) {
	return s.authBusiness.ActiveSessions(ctx, req)
}

func (s *grpcService) RevokeAllSessions(ctx context.Context, req *pb.RevokeAllSessionsRequest) (*pb.RevokeAllSessionsResponse, error) {
	return nil, s.authBusiness.RevokeAllSessions(ctx, req)
}

func (s *grpcService) Encrypt(ctx context.Context, req *pb.EncryptRequest) (*pb.EncryptResponse, error) {
	return s.authBusiness.Encrypt(ctx, req)
}
