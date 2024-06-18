package usergrpc

import (
	"context"

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
	return connect.NewResponse(&pb.ListUsersReply{
		Message: "Hello World",
	}), nil
}
