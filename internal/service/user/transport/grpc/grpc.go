package usergrpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"

	"github.com/anhnmt/go-api-boilerplate/proto/pb"
	"github.com/anhnmt/go-api-boilerplate/proto/pb/pbconnect"
)

type grpcService struct {
	pbconnect.UnimplementedUserServiceHandler

	business Business
}

func New(business Business) pbconnect.UserServiceHandler {
	return &grpcService{
		business: business,
	}
}

func (s *grpcService) ListUsers(context.Context, *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersReply], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.v1.UserService.ListUsers is not implemented"))
}
