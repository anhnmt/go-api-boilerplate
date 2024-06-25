package usergrpc

import (
	"context"

	"google.golang.org/grpc"

	userbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/user/business"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type grpcService struct {
	pb.UnimplementedUserServiceServer

	userBusiness *userbusiness.Business
}

func New(
	grpcSrv *grpc.Server,

	userBusiness *userbusiness.Business,
) pb.UserServiceServer {
	svc := &grpcService{
		userBusiness: userBusiness,
	}

	pb.RegisterUserServiceServer(grpcSrv, svc)

	return svc
}

func (s *grpcService) ListUsers(ctx context.Context, _ *pb.ListUsersRequest) (*pb.ListUsersReply, error) {
	res, err := s.userBusiness.ListUsers(ctx)
	if err != nil {
		return nil, err
	}

	return &pb.ListUsersReply{
		Data: res,
	}, nil
}

func (s *grpcService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	err := s.userBusiness.CreateUser(ctx, req)
	if err != nil {
		return nil, err
	}

	return &pb.CreateUserReply{
		Message: "Created user: " + req.Name,
	}, nil
}
