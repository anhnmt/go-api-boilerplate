package authbusiness

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anhnmt/go-api-boilerplate/internal/common/jwtutils"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	userentity "github.com/anhnmt/go-api-boilerplate/internal/service/user/entity"
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

	accessToken, tokenExpires, err := b.generateAccessToken(user, sessionId, now)
	if err != nil {
		return nil, err
	}

	refreshToken, refreshExpires, err := b.generateRefreshToken(user.ID, sessionId, now)
	if err != nil {
		return nil, err
	}

	res := &pb.LoginReply{
		TokenType:        jwtutils.TokenType,
		AccessToken:      accessToken,
		ExpiresAt:        tokenExpires,
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpires,
	}

	return res, nil
}

func (b *Business) Info(ctx context.Context) (*pb.InfoReply, error) {
	jwtToken, err := auth.AuthFromMD(ctx, jwtutils.TokenType)
	if err != nil {
		return nil, fmt.Errorf("failed get token")
	}

	token, err := jwtutils.ParseToken(jwtToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(b.cfg.Secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed parse token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims[jwtutils.Typ] != jwtutils.TokenType {
		return nil, fmt.Errorf("invalid token type")
	}

	res := &pb.InfoReply{
		Id:        cast.ToString(claims[jwtutils.Sub]),
		Email:     cast.ToString(claims[jwtutils.Email]),
		Name:      cast.ToString(claims[jwtutils.Name]),
		SessionId: cast.ToString(claims[jwtutils.Sid]),
	}

	return res, nil
}

func (b *Business) generateAccessToken(user *userentity.User, sessionId string, now time.Time) (string, int64, error) {
	tokenTime, err := time.ParseDuration(b.cfg.TokenExpires)
	if err != nil {
		return "", 0, fmt.Errorf("tokenExpires: %w", err)
	}

	tokenExpires := now.Add(tokenTime).Unix()
	accessToken, err := jwtutils.GenerateToken(jwt.MapClaims{
		jwtutils.Jti:   uuid.NewString(),
		jwtutils.Iat:   now.Unix(),
		jwtutils.Exp:   tokenExpires,
		jwtutils.Sid:   sessionId,
		jwtutils.Sub:   user.ID,
		jwtutils.Name:  user.Name,
		jwtutils.Email: user.Email,
		jwtutils.Typ:   jwtutils.TokenType,
	}, []byte(b.cfg.Secret))
	if err != nil {
		return "", 0, err
	}

	return accessToken, tokenExpires, nil
}

func (b *Business) generateRefreshToken(userId string, sessionId string, now time.Time) (string, int64, error) {
	refreshTime, err := time.ParseDuration(b.cfg.RefreshExpires)
	if err != nil {
		return "", 0, fmt.Errorf("tokenExpires: %w", err)
	}

	refreshExpires := now.Add(refreshTime).Unix()
	refreshToken, err := jwtutils.GenerateToken(jwt.MapClaims{
		jwtutils.Jti: uuid.NewString(),
		jwtutils.Iat: now.Unix(),
		jwtutils.Exp: refreshExpires,
		jwtutils.Sid: sessionId,
		jwtutils.Sub: userId,
		jwtutils.Typ: jwtutils.RefreshType,
	}, []byte(b.cfg.Secret))

	return refreshToken, refreshExpires, nil
}
