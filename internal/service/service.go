package service

import (
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

func initServices(
	_ pb.UserServiceServer,
	_ pb.AuthServiceServer,
) error {
	return nil
}
