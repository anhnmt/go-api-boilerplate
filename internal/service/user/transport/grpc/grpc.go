package usergrpc

import (
	"context"
	"fmt"

	"github.com/bufbuild/protovalidate-go"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
	v, err := protovalidate.New(protovalidate.WithFailFast(true))
	if err != nil {
		fmt.Println("failed to initialize validator:", err)
	}

	if err = v.Validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &pb.CreateUserReply{}, nil
}
