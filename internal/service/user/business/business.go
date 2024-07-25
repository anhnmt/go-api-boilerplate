package userbusiness

import (
	"context"
	"errors"

	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"

	"github.com/anhnmt/go-api-boilerplate/gen/pb"
	"github.com/anhnmt/go-api-boilerplate/internal/model"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/util"
	usercommand "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/command"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
)

type Business struct {
	userCommand *usercommand.Command
	userQuery   *userquery.Query
}

type Params struct {
	fx.In

	UserCommand *usercommand.Command
	UserQuery   *userquery.Query
}

func New(p Params) *Business {
	return &Business{
		userCommand: p.UserCommand,
		userQuery:   p.UserQuery,
	}
}

func (b *Business) ListUsers(ctx context.Context) ([]*pb.User, error) {
	data, err := b.userQuery.ListUsers(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, "failed to list users")
	}

	return b.ConvertToPbUser(data), nil
}

func (b *Business) ConvertToPbUser(users []*model.User) []*pb.User {
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

	createUser := &model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(password),
	}

	err = b.userCommand.Create(ctx, createUser)
	if err != nil {
		if errors.Is(util.GormTranslate(err), gorm.ErrDuplicatedKey) {
			return status.Error(codes.InvalidArgument, "user already exists")
		}

		return status.Error(codes.Internal, "failed to create user")
	}

	return nil
}
