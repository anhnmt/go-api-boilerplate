package authbusiness

import (
	"context"
	"fmt"
	"time"

	"github.com/anhnmt/go-fingerprint"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anhnmt/go-api-boilerplate/internal/common/jwtutils"
	"github.com/anhnmt/go-api-boilerplate/internal/core/entity"
	"github.com/anhnmt/go-api-boilerplate/internal/infrastructure/gormgen"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	sessionentity "github.com/anhnmt/go-api-boilerplate/internal/service/session/entity"
	sessioncommand "github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/command"
	userentity "github.com/anhnmt/go-api-boilerplate/internal/service/user/entity"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type Business struct {
	cfg            config.JWT
	userQuery      *userquery.Query
	sessionCommand *sessioncommand.Command
}

func New(
	cfg config.JWT,
	userQuery *userquery.Query,
	sessionCommand *sessioncommand.Command,
) *Business {
	return &Business{
		cfg:            cfg,
		userQuery:      userQuery,
		sessionCommand: sessionCommand,
	}
}

func (b *Business) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := b.userQuery.GetByEmail(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid password")
	}

	res := &pb.LoginResponse{
		TokenType: jwtutils.TokenType,
	}

	err = b.generateUserToken(ctx, user, res)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (b *Business) Info(ctx context.Context) (*pb.InfoResponse, error) {
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

	res := &pb.InfoResponse{
		Id:        cast.ToString(claims[jwtutils.Sub]),
		Email:     cast.ToString(claims[jwtutils.Email]),
		Name:      cast.ToString(claims[jwtutils.Name]),
		SessionId: cast.ToString(claims[jwtutils.Sid]),
	}

	return res, nil
}

func (b *Business) generateAccessToken(user *userentity.User, session *sessionentity.Session, now time.Time) (string, time.Time, error) {
	tokenTime, err := time.ParseDuration(b.cfg.TokenExpires)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("tokenExpires: %w", err)
	}

	tokenExpires := now.Add(tokenTime)
	accessToken, err := jwtutils.GenerateToken(jwt.MapClaims{
		jwtutils.Jti:   uuid.NewString(),
		jwtutils.Typ:   jwtutils.TokenType,
		jwtutils.Iat:   now,
		jwtutils.Exp:   tokenExpires.Unix(),
		jwtutils.Sid:   session.ID,
		jwtutils.Fgp:   session.Fingerprint,
		jwtutils.Sub:   user.ID,
		jwtutils.Name:  user.Name,
		jwtutils.Email: user.Email,
	}, []byte(b.cfg.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, tokenExpires, nil
}

func (b *Business) generateRefreshToken(userId string, session *sessionentity.Session, now time.Time) (string, time.Time, error) {
	refreshTime, err := time.ParseDuration(b.cfg.RefreshExpires)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("tokenExpires: %w", err)
	}

	refreshExpires := now.Add(refreshTime)
	refreshToken, err := jwtutils.GenerateToken(jwt.MapClaims{
		jwtutils.Jti: uuid.NewString(),
		jwtutils.Typ: jwtutils.RefreshType,
		jwtutils.Iat: now.Unix(),
		jwtutils.Exp: refreshExpires.Unix(),
		jwtutils.Sid: session.ID,
		jwtutils.Fgp: session.Fingerprint,
		jwtutils.Sub: userId,
	}, []byte(b.cfg.Secret))

	return refreshToken, refreshExpires, nil
}

func (b *Business) createUserSession(ctx context.Context, fg *fingerprint.Fingerprint, session *sessionentity.Session) error {
	if fg.IpAddress != nil {
		session.IpAddress = fg.IpAddress.Value
	}

	if fg.UserAgent != nil {
		session.UserAgent = fg.UserAgent.Raw

		if fg.UserAgent.Device != nil {
			session.Device = fg.UserAgent.Device.Name
			session.DeviceType = fg.UserAgent.Device.Type
		}

		if fg.UserAgent.OS != nil {
			session.OS = fg.UserAgent.OS.Name
		}

		if fg.UserAgent.Browser != nil {
			session.Browser = fg.UserAgent.Browser.Name
		}
	}

	err := b.sessionCommand.Create(ctx, session)
	if err != nil {
		return fmt.Errorf("failed create session: %w", err)
	}

	return nil
}

func (b *Business) generateUserToken(ctx context.Context, user *userentity.User, res *pb.LoginResponse) error {
	return b.sessionCommand.DB().Transaction(func(tx *gormgen.Query) error {
		now := time.Now().UTC()
		sessionId := uuid.NewString()
		fg := fingerprint.NewFingerprintContext(ctx)

		session := &sessionentity.Session{
			BaseEntity: entity.BaseEntity{
				ID:        sessionId,
				CreatedAt: now,
				UpdatedAt: now,
			},
			LastSeenAt:  &now,
			UserID:      user.ID,
			Fingerprint: fg.ID,
		}

		accessToken, tokenExpires, err := b.generateAccessToken(user, session, now)
		if err != nil {
			return err
		}

		refreshToken, refreshExpires, err := b.generateRefreshToken(user.ID, session, now)
		if err != nil {
			return err
		}

		session.ExpiresAt = &refreshExpires
		err = b.createUserSession(ctx, fg, session)
		if err != nil {
			return err
		}

		res.AccessToken = accessToken
		res.ExpiresAt = tokenExpires.Unix()
		res.RefreshToken = refreshToken
		res.RefreshExpiresAt = refreshExpires.Unix()

		return nil
	})
}
