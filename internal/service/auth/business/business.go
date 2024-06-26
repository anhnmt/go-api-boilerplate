package authbusiness

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anhnmt/go-api-boilerplate/internal/common/jwtutils"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type Business struct {
	userQuery *userquery.Query
}

func New(
	userQuery *userquery.Query,
) *Business {
	return &Business{
		userQuery: userQuery,
	}
}

func (b *Business) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginReply, error) {
	user, err := b.userQuery.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid password")
	}

	sessionId := uuid.NewString()
	now := time.Now().UTC()
	secret := []byte("mysecretkey")

	tokenExpires := now.Add(time.Minute * 10)
	accessToken, err := jwtutils.GenerateToken(jwt.MapClaims{
		"jti":   uuid.NewString(),
		"iat":   now.Unix(),
		"exp":   tokenExpires.Unix(),
		"sid":   sessionId,
		"sub":   user.ID,
		"name":  user.Name,
		"email": user.Email,
		"typ":   jwtutils.TokenType,
	}, secret)
	if err != nil {
		return nil, err
	}

	refreshExpires := now.Add(time.Hour * 24)
	refreshToken, err := jwtutils.GenerateToken(jwt.MapClaims{
		"jti": uuid.NewString(),
		"iat": now.Unix(),
		"exp": refreshExpires.Unix(),
		"sid": sessionId,
		"sub": user.ID,
		"typ": jwtutils.RefreshType,
	}, secret)
	if err != nil {
		return nil, err
	}

	res := &pb.LoginReply{
		TokenType:        jwtutils.TokenType,
		AccessToken:      accessToken,
		ExpiresAt:        tokenExpires.Unix(),
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpires.Unix(),
	}

	return res, nil
}
