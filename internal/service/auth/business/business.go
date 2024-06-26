package authbusiness

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anhnmt/go-api-boilerplate/internal/common/jwtutils"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type Business struct {
	cfg       config.JWT
	userQuery *userquery.Query
}

func New(
	cfg config.JWT,
	userQuery *userquery.Query,
) *Business {
	return &Business{
		cfg:       cfg,
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
	secret := []byte(b.cfg.Secret)

	tokenTime, err := time.ParseDuration(b.cfg.TokenExpires)
	if err != nil {
		return nil, fmt.Errorf("tokenExpires: %w", err)
	}

	tokenExpires := now.Add(tokenTime)
	accessToken, err := jwtutils.GenerateToken(jwt.MapClaims{
		jwtutils.Jti:   uuid.NewString(),
		jwtutils.Iat:   now.Unix(),
		jwtutils.Exp:   tokenExpires.Unix(),
		jwtutils.Sid:   sessionId,
		jwtutils.Sub:   user.ID,
		jwtutils.Name:  user.Name,
		jwtutils.Email: user.Email,
		jwtutils.Typ:   jwtutils.TokenType,
	}, secret)
	if err != nil {
		return nil, err
	}

	refreshTime, err := time.ParseDuration(b.cfg.RefreshExpires)
	if err != nil {
		return nil, fmt.Errorf("tokenExpires: %w", err)
	}

	refreshExpires := now.Add(refreshTime)
	refreshToken, err := jwtutils.GenerateToken(jwt.MapClaims{
		jwtutils.Jti: uuid.NewString(),
		jwtutils.Iat: now.Unix(),
		jwtutils.Exp: refreshExpires.Unix(),
		jwtutils.Sid: sessionId,
		jwtutils.Sub: user.ID,
		jwtutils.Typ: jwtutils.RefreshType,
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

func (b *Business) Info(_ context.Context, req *pb.InfoRequest) error {
	token, err := jwtutils.ParseToken(req.Token, func(token *jwt.Token) (interface{}, error) {
		return []byte(b.cfg.Secret), nil
	})
	if err != nil {
		return err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return fmt.Errorf("invalid token")
	}

	log.Info().Any("claims", claims).Msg("token")

	return nil
}
