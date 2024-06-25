package userbusiness

import (
	"context"
	"errors"

	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/anhnmt/go-api-boilerplate/internal/common"
	userentity "github.com/anhnmt/go-api-boilerplate/internal/service/user/entity"
	usercommand "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/command"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type Business struct {
	userCommand *usercommand.Command
	userQuery   *userquery.Query
}

func New(
	userCommand *usercommand.Command,
	userQuery *userquery.Query,
) *Business {
	return &Business{
		userCommand: userCommand,
		userQuery:   userQuery,
	}
}

func (b *Business) ListUsers(ctx context.Context) ([]*pb.User, error) {
	data, err := b.userQuery.ListUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	return b.ConvertToPbUser(data), nil
}

func (b *Business) ConvertToPbUser(users []*userentity.User) []*pb.User {
	res := make([]*pb.User, len(users))

	for i, user := range users {
		res[i] = &pb.User{
			Id:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}
	}

	return res
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
		if errors.Is(common.GormTranslate(err), gorm.ErrDuplicatedKey) {
			return status.Error(codes.InvalidArgument, "user already exists")
		}

		return status.Error(codes.Internal, "failed to create user")
	}

	return nil
}
