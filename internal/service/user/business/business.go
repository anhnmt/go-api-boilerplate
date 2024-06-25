package userbusiness

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	userentity "github.com/anhnmt/go-api-boilerplate/internal/service/user/entity"
	usercommand "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/command"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type Business struct {
	userCommand *usercommand.Command
}

func New(
	userCommand *usercommand.Command,
) *Business {
	return &Business{
		userCommand: userCommand,
	}
}

func (b *Business) CreateUser(ctx context.Context, req *pb.CreateUserRequest) error {
	password, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return status.Error(codes.Internal, "failed to hash password")
	}

	createUser := &userentity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(password),
	}

	err = b.userCommand.Create(ctx, createUser)
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return status.Error(codes.InvalidArgument, "user already exists")
		}

		return status.Error(codes.Internal, "failed to create user")
	}

	return nil
}
