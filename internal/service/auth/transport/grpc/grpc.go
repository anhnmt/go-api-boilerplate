package authgrpc

import (
	"context"

	"github.com/golang-jwt/jwt/v5"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/internal/common/jwtutils"
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

func (s *grpcService) Info(_ context.Context, req *pb.InfoRequest) (*pb.InfoReply, error) {
	token, err := jwtutils.ParseToken(req.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte("mysecretkey"), nil
	})

	if err != nil {
		return nil, err
	}

	log.Info().Any("token", token).Msg("token")

	return nil, nil
}
