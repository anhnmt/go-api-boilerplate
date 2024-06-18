package usergrpc

import (
	"context"
	"errors"
	"net/http"

	"connectrpc.com/connect"

	userbusiness "github.com/anhnmt/go-api-boilerplate/internal/service/user/business"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
	"github.com/anhnmt/go-api-boilerplate/proto/pb/pbconnect"
)

type grpcService struct {
	pbconnect.UnimplementedUserServiceHandler

	business userbusiness.Business
}

func New(
	mux *http.ServeMux,
	// business userbusiness.Business,
) pbconnect.UserServiceHandler {
	svc := &grpcService{
		// business: business,
	}

	mux.Handle(pbconnect.NewUserServiceHandler(svc))
	return svc
}

func (s *grpcService) ListUsers(context.Context, *connect.Request[pb.ListUsersRequest]) (*connect.Response[pb.ListUsersReply], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("user.v1.UserService.ListUsers is not implemented"))
}
