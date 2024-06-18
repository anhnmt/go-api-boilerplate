package usergrpc

import (
	"context"
	"errors"

	"connectrpc.com/connect"
	"connectrpc.com/vanguard"

	userbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/user/business"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
	"github.com/anhnmt/go-api-boilerplate/proto/pb/pbconnect"
)

type grpcService struct {
	pbconnect.UnimplementedUserServiceHandler

	business userbusiness.Business
}

func New(
	services *[]*vanguard.Service,
	// business userbusiness.Business,
) pbconnect.UserServiceHandler {
	svc := &grpcService{
		// business: business,
	}

	*services = append(*services, vanguard.NewService(
		pbconnect.NewUserServiceHandler(svc),
	))

	return svc
}

func (s *grpcService) ListUsers(context.Context, *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersReply], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.v1.UserService.ListUsers is not implemented"))
}
