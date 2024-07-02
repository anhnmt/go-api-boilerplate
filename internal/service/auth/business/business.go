package authbusiness

import (
	"context"
	"fmt"
	"time"

	"github.com/anhnmt/go-fingerprint"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anhnmt/go-api-boilerplate/internal/common/ctxutils"
	"github.com/anhnmt/go-api-boilerplate/internal/common/jwtutils"
	"github.com/anhnmt/go-api-boilerplate/internal/infrastructure/core/entity"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	authredis "github.com/anhnmt/go-api-boilerplate/internal/service/auth/repository/redis"
	sessionentity "github.com/anhnmt/go-api-boilerplate/internal/service/session/entity"
	sessioncommand "github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/command"
	sessionquery "github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/query"
	userentity "github.com/anhnmt/go-api-boilerplate/internal/service/user/entity"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

type Business struct {
	cfg            config.JWT
	userQuery      *userquery.Query
	sessionCommand *sessioncommand.Command
	sessionQuery   *sessionquery.Query
	authRedis      *authredis.Redis
}

func New(
	cfg config.JWT,
	userQuery *userquery.Query,
	sessionCommand *sessioncommand.Command,
	sessionQuery *sessionquery.Query,
	authRedis *authredis.Redis,
) *Business {
	return &Business{
		cfg:            cfg,
		userQuery:      userQuery,
		sessionCommand: sessionCommand,
		sessionQuery:   sessionQuery,
		authRedis:      authRedis,
	}
}

func (b *Business) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	user, err := b.userQuery.GetByEmailWithPassword(ctx, req.Email)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid email or password")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password))
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid password")
	}

	sessionId := uuid.NewString()
	accessToken, tokenExpires, refreshToken, refreshExpires, err := b.generateUserToken(ctx, user, sessionId)
	if err != nil {
		return nil, err
	}

	res := &pb.LoginResponse{
		TokenType:        jwtutils.TokenType,
		AccessToken:      accessToken,
		ExpiresAt:        tokenExpires.Unix(),
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpires.Unix(),
	}

	return res, err
}

func (b *Business) Info(ctx context.Context) (*pb.InfoResponse, error) {
	claims, err := ctxutils.ExtractCtxClaims(ctx)
	if err != nil {
		return nil, err
	}

	if claims[jwtutils.Typ] != jwtutils.TokenType {
		return nil, fmt.Errorf("invalid token type")
	}

	sessionId := cast.ToString(claims[jwtutils.Sid])
	tokenId := cast.ToString(claims[jwtutils.Jti])
	err = b.CheckBlacklist(ctx, sessionId, tokenId)
	if err != nil {
		return nil, err
	}

	email := cast.ToString(claims[jwtutils.Email])
	user, err := b.userQuery.GetByEmail(ctx, email)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed get user info")
	}

	res := &pb.InfoResponse{
		Id:        user.ID,
		Email:     email,
		Name:      user.Name,
		SessionId: sessionId,
		TokenId:   tokenId,
	}

	return res, nil
}

func (b *Business) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	claims, err := b.ExtractClaims(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if claims[jwtutils.Typ] != jwtutils.RefreshType {
		return nil, fmt.Errorf("invalid refresh token")
	}

	sessionId := cast.ToString(claims[jwtutils.Sid])
	tokenId := cast.ToString(claims[jwtutils.Jti])
	err = b.CheckBlacklist(ctx, sessionId, tokenId)
	if err != nil {
		return nil, err
	}

	userId := cast.ToString(claims[jwtutils.Sub])
	user, err := b.userQuery.GetByID(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user")
	}

	err = b.authRedis.SetTokenBlacklist(ctx, tokenId, sessionId)
	if err != nil {
		return nil, err
	}

	err = b.sessionCommand.UpdateLastSeenAt(ctx, sessionId, time.Now().UTC())
	if err != nil {
		return nil, err
	}

	accessToken, tokenExpires, refreshToken, refreshExpires, err := b.generateUserToken(ctx, user, sessionId)
	if err != nil {
		return nil, err
	}

	res := &pb.RefreshTokenResponse{
		TokenType:        jwtutils.TokenType,
		AccessToken:      accessToken,
		ExpiresAt:        tokenExpires.Unix(),
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpires.Unix(),
	}

	return res, nil
}

func (b *Business) RevokeToken(ctx context.Context) error {
	claims, err := ctxutils.ExtractCtxClaims(ctx)
	if err != nil {
		return err
	}

	if claims[jwtutils.Typ] != jwtutils.RefreshType {
		return fmt.Errorf("invalid refresh token")
	}

	sessionId := cast.ToString(claims[jwtutils.Sid])
	tokenId := cast.ToString(claims[jwtutils.Jti])
	err = b.CheckBlacklist(ctx, sessionId, tokenId)
	if err != nil {
		return err
	}

	err = b.sessionCommand.UpdateIsRevoked(ctx, sessionId, true, time.Now().UTC())
	if err != nil {
		return err
	}

	err = b.authRedis.SetSessionBlacklist(ctx, sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (b *Business) ActiveSessions(ctx context.Context, _ *pb.ActiveSessionsRequest) (*pb.ActiveSessionsResponse, error) {
	claims, err := ctxutils.ExtractCtxClaims(ctx)
	if err != nil {
		return nil, err
	}

	if claims[jwtutils.Typ] != jwtutils.TokenType {
		return nil, fmt.Errorf("invalid token")
	}

	sessionId := cast.ToString(claims[jwtutils.Sid])
	tokenId := cast.ToString(claims[jwtutils.Jti])
	err = b.CheckBlacklist(ctx, sessionId, tokenId)
	if err != nil {
		return nil, err
	}

	userId := cast.ToString(claims[jwtutils.Sub])
	sessions, err := b.sessionQuery.FindByUserIdAndSessionId(ctx, userId, sessionId)
	if err != nil {
		return nil, err
	}

	res := &pb.ActiveSessionsResponse{
		Data: sessions,
	}

	return res, nil
}

func (b *Business) CheckBlacklist(ctx context.Context, sessionId, tokenId string) error {
	if err := b.authRedis.CheckSessionBlacklist(ctx, sessionId); err != nil {
		return err
	}

	if err := b.authRedis.CheckTokenBlacklist(ctx, tokenId); err != nil {
		return err
	}

	return nil
}

func (b *Business) ExtractClaims(rawToken string) (jwt.MapClaims, error) {
	token, err := jwtutils.ParseToken(rawToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(b.cfg.Secret), nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (b *Business) generateAccessToken(user *userentity.User, tokenId, sessionId, fingerprint string, now time.Time) (string, time.Time, error) {
	tokenTime, err := time.ParseDuration(b.cfg.TokenExpires)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("tokenExpires: %w", err)
	}

	tokenExpires := now.Add(tokenTime)
	accessToken, err := jwtutils.GenerateToken(jwt.MapClaims{
		jwtutils.Jti:   tokenId,
		jwtutils.Typ:   jwtutils.TokenType,
		jwtutils.Iat:   now,
		jwtutils.Exp:   tokenExpires.Unix(),
		jwtutils.Sid:   sessionId,
		jwtutils.Fgp:   fingerprint,
		jwtutils.Sub:   user.ID,
		jwtutils.Name:  user.Name,
		jwtutils.Email: user.Email,
	}, []byte(b.cfg.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, tokenExpires, nil
}

func (b *Business) generateRefreshToken(userId, tokenId, sessionId, fingerprint string, now time.Time) (string, time.Time, error) {
	refreshTime, err := time.ParseDuration(b.cfg.RefreshExpires)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("tokenExpires: %w", err)
	}

	refreshExpires := now.Add(refreshTime)
	refreshToken, err := jwtutils.GenerateToken(jwt.MapClaims{
		jwtutils.Jti: tokenId,
		jwtutils.Typ: jwtutils.RefreshType,
		jwtutils.Iat: now.Unix(),
		jwtutils.Exp: refreshExpires.Unix(),
		jwtutils.Sid: sessionId,
		jwtutils.Fgp: fingerprint,
		jwtutils.Sub: userId,
	}, []byte(b.cfg.Secret))

	return refreshToken, refreshExpires, nil
}

func (b *Business) createUserSession(ctx context.Context, fg *fingerprint.Fingerprint, userId, sessionId string, now, refreshExpires time.Time) error {
	session := &sessionentity.Session{
		BaseEntity: entity.BaseEntity{
			ID:        sessionId,
			CreatedAt: now,
			UpdatedAt: now,
		},
		LastSeenAt:  &now,
		UserID:      userId,
		Fingerprint: fg.ID,
		ExpiresAt:   &refreshExpires,
	}

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

	err := b.sessionCommand.CreateOnConflict(ctx, session)
	if err != nil {
		return fmt.Errorf("failed create session: %w", err)
	}

	return nil
}

func (b *Business) generateUserToken(ctx context.Context, user *userentity.User, sessionId string) (accessToken string, tokenExpires time.Time, refreshToken string, refreshExpires time.Time, err error) {
	now := time.Now().UTC()
	tokenId := uuid.NewString()
	fg := fingerprint.NewFingerprintContext(ctx)

	accessToken, tokenExpires, err = b.generateAccessToken(user, tokenId, sessionId, fg.ID, now)
	if err != nil {
		return
	}

	refreshToken, refreshExpires, err = b.generateRefreshToken(user.ID, tokenId, sessionId, fg.ID, now)
	if err != nil {
		return
	}

	err = b.createUserSession(ctx, fg, user.ID, sessionId, now, refreshExpires)
	return
}
