package authbusiness

import (
	"context"
	"fmt"
	"time"

	"github.com/anhnmt/go-fingerprint"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cast"
	"go.uber.org/fx"
	"golang.org/x/crypto/bcrypt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/anhnmt/go-api-boilerplate/gen/pb"
	"github.com/anhnmt/go-api-boilerplate/internal/model"

	"github.com/anhnmt/go-api-boilerplate/internal/common/ctxutils"
	"github.com/anhnmt/go-api-boilerplate/internal/common/jwtutils"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	authredis "github.com/anhnmt/go-api-boilerplate/internal/service/auth/repository/redis"
	sessioncommand "github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/command"
	sessionquery "github.com/anhnmt/go-api-boilerplate/internal/service/session/repository/postgres/query"
	userquery "github.com/anhnmt/go-api-boilerplate/internal/service/user/repository/postgres/query"
)

type Business struct {
	config         config.JWT
	userQuery      *userquery.Query
	sessionCommand *sessioncommand.Command
	sessionQuery   *sessionquery.Query
	authRedis      *authredis.Redis
}

type Params struct {
	fx.In

	Config         config.JWT
	UserQuery      *userquery.Query
	SessionCommand *sessioncommand.Command
	SessionQuery   *sessionquery.Query
	AuthRedis      *authredis.Redis
}

func New(p Params) *Business {
	return &Business{
		config:         p.Config,
		userQuery:      p.UserQuery,
		sessionCommand: p.SessionCommand,
		sessionQuery:   p.SessionQuery,
		authRedis:      p.AuthRedis,
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
	fg := fingerprint.NewFingerprintContext(ctx)

	accessToken, tokenExpires, refreshToken, refreshExpires, err := b.generateUserToken(ctx, fg, user, sessionId)
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

	fg := fingerprint.NewFingerprintContext(ctx)
	sessionId := cast.ToString(claims[jwtutils.Sid])
	tokenId := cast.ToString(claims[jwtutils.Jti])

	err = b.CheckBlacklist(ctx, sessionId, tokenId)
	if err != nil {
		fingerId := cast.ToString(claims[jwtutils.Fgp])

		// detect leaked token here
		err = b.detectLeakedToken(ctx, fingerId, fg.ID, sessionId)
		if err != nil {
			return nil, err
		}

		return nil, err
	}

	userId := cast.ToString(claims[jwtutils.Sub])
	user, err := b.userQuery.GetByID(ctx, userId)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user")
	}

	err = b.updateLastSeenAtAndBlacklist(ctx, tokenId, sessionId)
	if err != nil {
		return nil, err
	}

	accessToken, tokenExpires, refreshToken, refreshExpires, err := b.generateUserToken(ctx, fg, user, sessionId)
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

	err = b.updateRevokedAndBlacklist(ctx, sessionId)
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

func (b *Business) RevokeAllSessions(ctx context.Context, req *pb.RevokeAllSessionsRequest) error {
	claims, err := b.ExtractClaims(req.RefreshToken)
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

	userId := cast.ToString(claims[jwtutils.Sub])

	var revokeAll string
	if !req.RevokeCurrent {
		revokeAll = sessionId
	}

	sessionIds, err := b.sessionQuery.FindByUserIdWithoutSessionId(ctx, userId, revokeAll)
	if err != nil {
		return err
	}

	if len(sessionIds) > 0 {
		err = b.authRedis.SetSessionsBlacklist(ctx, sessionIds)
		if err != nil {
			return err
		}

		err = b.sessionCommand.UpdateRevokedByUserIdWithoutSessionId(ctx, userId, revokeAll)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Business) Encrypt(_ context.Context, req *pb.EncryptRequest) (*pb.EncryptResponse, error) {
	log.Info().Str("key", req.Data).Msg("encrypting data")

	res := &pb.EncryptResponse{
		Data: "alo",
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
		return []byte(b.config.Secret), nil
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

func (b *Business) detectLeakedToken(ctx context.Context, claimsFingerId, fingerId, sessionId string) error {
	if claimsFingerId != fingerId {
		// revoke current session
		err := b.updateRevokedAndBlacklist(ctx, sessionId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Business) updateLastSeenAtAndBlacklist(ctx context.Context, tokenId, sessionId string) error {
	err := b.authRedis.SetTokenBlacklist(ctx, tokenId, sessionId)
	if err != nil {
		return err
	}

	err = b.sessionCommand.UpdateLastSeenAt(ctx, sessionId, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}

func (b *Business) updateRevokedAndBlacklist(ctx context.Context, sessionId string) error {
	err := b.sessionCommand.UpdateIsRevoked(ctx, sessionId, true, time.Now().UTC())
	if err != nil {
		return err
	}

	err = b.authRedis.SetSessionBlacklist(ctx, sessionId)
	if err != nil {
		return err
	}

	return nil
}

func (b *Business) generateAccessToken(user *model.User, tokenId, sessionId, fingerprint string, now time.Time) (string, time.Time, error) {
	tokenTime, err := time.ParseDuration(b.config.TokenExpires)
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
	}, []byte(b.config.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, tokenExpires, nil
}

func (b *Business) generateRefreshToken(userId, tokenId, sessionId, fingerprint string, now time.Time) (string, time.Time, error) {
	refreshTime, err := time.ParseDuration(b.config.RefreshExpires)
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
	}, []byte(b.config.Secret))

	return refreshToken, refreshExpires, nil
}

func (b *Business) createUserSession(ctx context.Context, fg *fingerprint.Fingerprint, userId, sessionId string, now, refreshExpires time.Time) error {
	session := &model.Session{
		BaseModel: model.BaseModel{
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

func (b *Business) generateUserToken(ctx context.Context, fg *fingerprint.Fingerprint, user *model.User, sessionId string) (accessToken string, tokenExpires time.Time, refreshToken string, refreshExpires time.Time, err error) {
	now := time.Now().UTC()
	tokenId := uuid.NewString()

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
