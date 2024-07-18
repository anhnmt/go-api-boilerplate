package usergrpc

import (
	"context"

	"go.uber.org/fx"
	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/gen/pb"
	userbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/user/business"
)

type grpcService struct {
	pb.UnimplementedUserServiceServer

	userBusiness *userbusiness.Business
}

type Params struct {
	fx.In

	GrpcServer   *grpc.Server
	UserBusiness *userbusiness.Business
}

func New(p Params) pb.UserServiceServer {
	svc := &grpcService{
		userBusiness: p.UserBusiness,
	}

	pb.RegisterUserServiceServer(p.GrpcServer, svc)
	return svc
}

func (s *grpcService) ListUsers(ctx context.Context, _ *pb.ListUsersRequest) (*pb.ListUsersResponse, error) {
	res, err := s.userBusiness.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListUsersResponse{
		Data: res,
	}, nil
}

func (s *grpcService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	err := s.userBusiness.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserResponse{
		Message: "Created user: " + req.Name,
	}, nil
}
