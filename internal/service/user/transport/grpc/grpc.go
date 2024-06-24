package usergrpc

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/anhnmt/go-api-boilerplate/internal/core/entity"
	userentity "github.com/anhnmt/go-api-boilerplate/internal/service/user/entity"
	usercommand "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/command"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type grpcService struct {
	pb.UnimplementedUserServiceServer

	userCommand *usercommand.Command
}

func New(
	grpcSrv *grpc.Server,

	userCommand *usercommand.Command,
) pb.UserServiceServer {
	svc := &grpcService{
		userCommand: userCommand,
	}

	pb.RegisterUserServiceServer(grpcSrv, svc)

	return svc
}

func (s *grpcService) ListUsers(context.Context, *pb.ListUsersRequest) (*pb.ListUsersReply, error) {
	return &pb.ListUsersReply{
		Message: "Hello World",
	}, nil
}

func (s *grpcService) CreateUser(ctx context.Context, req *pb.CreateUserRequest) (*pb.CreateUserReply, error) {
	now := time.Now().UTC()

	createUser := &userentity.User{
		BaseEntity: entity.BaseEntity{
			ID:        uuid.NewString(),
			CreatedAt: now,
			UpdatedAt: now,
		},
		Name:  req.Name,
		Email: req.Email,
	}

	err := s.userCommand.Create(ctx, createUser)
	if err != nil {
		if errors.As(err, &gorm.ErrDuplicatedKey) {
			return nil, status.Error(codes.InvalidArgument, "user already exists")
		}

		return nil, status.Error(codes.Internal, "failed to create user")
	}

	return &pb.CreateUserReply{
		Message: "Created user: " + req.Name,
	}, nil
}
