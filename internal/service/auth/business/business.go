package authbusiness

import (
	"context"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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

	tokenExpires := now.Add(time.Minute * 10)
	accessToken, err := generateAccessToken(user.ID, sessionId, now, tokenExpires)
	if err != nil {
		return nil, err
	}

	res := &pb.LoginReply{
		TokenType:   "Bearer",
		AccessToken: accessToken,
		ExpiresAt:   tokenExpires.Unix(),
	}

	return res, nil
}

func generateAccessToken(userId, sessionId string, tokenStart, tokenExpires time.Time) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"session_id": sessionId,
		"sub":        userId,
		"iat":        tokenStart.Unix(),
		"exp":        tokenExpires.Unix(),
	})

	return token.SignedString([]byte("secret"))
}
