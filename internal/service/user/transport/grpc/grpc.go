package usergrpc

import (
	"context"

	"google.golang.org/grpc"

	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type grpcService struct {
	pb.UnimplementedUserServiceServer
}

func New(
	grpcSrv *grpc.Server,
) pb.UserServiceServer {
	svc := &grpcService{}

	pb.RegisterUserServiceServer(grpcSrv, svc)

	return svc
}

func (s *grpcService) ListUsers(context.Context, *pb.ListUsersRequest) (*pb.ListUsersReply, error) {
	return &pb.ListUsersReply{
		Message: "Hello World",
	}, nil
}

func (s *grpcService) CreateUser(_ context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	return &pb.CreateUserReply{
		Message: "Created user: " + req.Name,
	}, nil
}
