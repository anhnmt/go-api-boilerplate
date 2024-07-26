package authbusiness

import (
	"context"
	"fmt"
	"time"

	"github.com/anhnmt/go-fingerprint"
	"github.com/casbin/casbin/v2"
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
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/util"

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
	rbac           *casbin.Enforcer
}

type Params struct {
	fx.In

	Config         config.JWT
	UserQuery      *userquery.Query
	SessionCommand *sessioncommand.Command
	SessionQuery   *sessionquery.Query
	AuthRedis      *authredis.Redis
	RBAC           *casbin.Enforcer
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

	sessionID := uuid.NewString()
	fg := fingerprint.NewFingerprintContext(ctx)

	accessToken, tokenExpires, refreshToken, refreshExpires, err := b.generateUserToken(ctx, fg, user, sessionID)
	if err != nil {
		return nil, err
	}

	res := &pb.LoginResponse{
		TokenType:        util.TokenType,
		AccessToken:      accessToken,
		ExpiresAt:        tokenExpires.Unix(),
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpires.Unix(),
	}

	return res, err
}

func (b *Business) Info(ctx context.Context) (*pb.InfoResponse, error) {
	claims, err := util.ExtractCtxClaims(ctx)
	if err != nil {
		return nil, err
	}

	sessionID := cast.ToString(claims[util.Sid])
	tokenID := cast.ToString(claims[util.Jti])
	err = b.CheckBlacklist(ctx, sessionID, tokenID)
	if err != nil {
		return nil, err
	}

	email := cast.ToString(claims[util.Email])
	user, err := b.userQuery.GetByEmail(ctx, email)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "failed get user info")
	}

	res := &pb.InfoResponse{
		Id:        user.ID,
		Email:     email,
		Name:      user.Name,
		SessionId: sessionID,
		TokenId:   tokenID,
		Role:      user.Role,
	}

	return res, nil
}

func (b *Business) RefreshToken(ctx context.Context, req *pb.RefreshTokenRequest) (*pb.RefreshTokenResponse, error) {
	claims, err := b.ExtractClaims(req.RefreshToken)
	if err != nil {
		return nil, err
	}

	if claims[util.Typ] != util.RefreshType {
		return nil, fmt.Errorf("invalid refresh token")
	}

	fg := fingerprint.NewFingerprintContext(ctx)
	sessionID := cast.ToString(claims[util.Sid])
	tokenID := cast.ToString(claims[util.Jti])

	err = b.CheckBlacklist(ctx, sessionID, tokenID)
	if err != nil {
		fingerID := cast.ToString(claims[util.Fgp])

		// detect leaked token here
		err2 := b.detectLeakedToken(ctx, fingerID, fg.ID, sessionID)
		if err2 != nil {
			return nil, err2
		}

		return nil, err
	}

	userID := cast.ToString(claims[util.Sub])
	user, err := b.userQuery.GetByID(ctx, userID)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid user")
	}

	err = b.updateLastSeenAtAndBlacklist(ctx, tokenID, sessionID)
	if err != nil {
		return nil, err
	}

	accessToken, tokenExpires, refreshToken, refreshExpires, err := b.generateUserToken(ctx, fg, user, sessionID)
	if err != nil {
		return nil, err
	}

	res := &pb.RefreshTokenResponse{
		TokenType:        util.TokenType,
		AccessToken:      accessToken,
		ExpiresAt:        tokenExpires.Unix(),
		RefreshToken:     refreshToken,
		RefreshExpiresAt: refreshExpires.Unix(),
	}

	return res, nil
}

func (b *Business) RevokeToken(ctx context.Context, req *pb.RevokeTokenRequest) error {
	claims, err := b.ExtractClaims(req.RefreshToken)
	if err != nil {
		return err
	}

	if claims[util.Typ] != util.RefreshType {
		return fmt.Errorf("invalid refresh token")
	}

	sessionID := cast.ToString(claims[util.Sid])
	tokenID := cast.ToString(claims[util.Jti])
	err = b.CheckBlacklist(ctx, sessionID, tokenID)
	if err != nil {
		return err
	}

	err = b.updateRevokedAndBlacklist(ctx, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (b *Business) ActiveSessions(ctx context.Context, req *pb.ActiveSessionsRequest) (*pb.ActiveSessionsResponse, error) {
	claims, err := util.ExtractCtxClaims(ctx)
	if err != nil {
		return nil, err
	}

	sessionID := cast.ToString(claims[util.Sid])
	tokenID := cast.ToString(claims[util.Jti])
	err = b.CheckBlacklist(ctx, sessionID, tokenID)
	if err != nil {
		return nil, err
	}

	page := int(req.GetPage())
	limit := int(req.GetLimit())
	if limit == 0 {
		limit = 10
	}

	userID := cast.ToString(claims[util.Sub])
	total, err := b.sessionQuery.CountByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	offset := util.GetOffset(page, limit)
	pageCount := util.TotalPage(total, limit)
	sessions, err := b.sessionQuery.FindByUserIDAndSessionID(ctx, userID, sessionID, limit, offset)
	if err != nil {
		return nil, err
	}

	res := &pb.ActiveSessionsResponse{
		Data:      sessions,
		Count:     int32(len(sessions)),
		PageCount: int32(pageCount),
		Total:     int32(total),
	}

	return res, nil
}

func (b *Business) RevokeAllSessions(ctx context.Context, req *pb.RevokeAllSessionsRequest) error {
	claims, err := b.ExtractClaims(req.RefreshToken)
	if err != nil {
		return err
	}

	if claims[util.Typ] != util.RefreshType {
		return fmt.Errorf("invalid refresh token")
	}

	sessionID := cast.ToString(claims[util.Sid])
	tokenID := cast.ToString(claims[util.Jti])
	err = b.CheckBlacklist(ctx, sessionID, tokenID)
	if err != nil {
		return err
	}

	userID := cast.ToString(claims[util.Sub])

	var revokeAll string
	if !req.RevokeCurrent {
		revokeAll = sessionID
	}

	sessionIDs, err := b.sessionQuery.FindByUserIDWithoutSessionID(ctx, userID, revokeAll)
	if err != nil {
		return err
	}

	if len(sessionIDs) > 0 {
		err = b.authRedis.SetSessionsBlacklist(ctx, sessionIDs)
		if err != nil {
			return err
		}

		err = b.sessionCommand.UpdateRevokedByUserIDWithoutSessionID(ctx, userID, revokeAll)
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

func (b *Business) CheckBlacklist(ctx context.Context, sessionID, tokenID string) error {
	if err := b.authRedis.CheckSessionBlacklist(ctx, sessionID); err != nil {
		return err
	}

	if err := b.authRedis.CheckTokenBlacklist(ctx, tokenID); err != nil {
		return err
	}

	return nil
}

func (b *Business) ExtractClaims(rawToken string) (jwt.MapClaims, error) {
	token, err := util.ParseToken(rawToken, func(token *jwt.Token) (interface{}, error) {
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

func (b *Business) detectLeakedToken(ctx context.Context, claimsFingerID, fingerID, sessionID string) error {
	if claimsFingerID != fingerID {
		// revoke current session
		err := b.updateRevokedAndBlacklist(ctx, sessionID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Business) updateLastSeenAtAndBlacklist(ctx context.Context, tokenID, sessionID string) error {
	err := b.authRedis.SetTokenBlacklist(ctx, tokenID, sessionID)
	if err != nil {
		return err
	}

	err = b.sessionCommand.UpdateLastSeenAt(ctx, sessionID, time.Now().UTC())
	if err != nil {
		return err
	}

	return nil
}

func (b *Business) updateRevokedAndBlacklist(ctx context.Context, sessionID string) error {
	err := b.sessionCommand.UpdateIsRevoked(ctx, sessionID, true, time.Now().UTC())
	if err != nil {
		return err
	}

	err = b.authRedis.SetSessionBlacklist(ctx, sessionID)
	if err != nil {
		return err
	}

	return nil
}

func (b *Business) generateAccessToken(user *model.User, tokenID, sessionID, fingerprint string, now time.Time) (string, time.Time, error) {
	tokenTime, err := time.ParseDuration(b.config.TokenExpires)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("tokenExpires: %w", err)
	}

	tokenExpires := now.Add(tokenTime)
	accessToken, err := util.GenerateToken(jwt.MapClaims{
		util.Jti:   tokenID,
		util.Typ:   util.TokenType,
		util.Iat:   now,
		util.Exp:   tokenExpires.Unix(),
		util.Sid:   sessionID,
		util.Fgp:   fingerprint,
		util.Sub:   user.ID,
		util.Name:  user.Name,
		util.Email: user.Email,
		util.Role:  user.Role,
	}, []byte(b.config.Secret))
	if err != nil {
		return "", time.Time{}, err
	}

	return accessToken, tokenExpires, nil
}

func (b *Business) generateRefreshToken(userID, tokenID, sessionID, fingerprint string, now time.Time) (string, time.Time, error) {
	refreshTime, err := time.ParseDuration(b.config.RefreshExpires)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("tokenExpires: %w", err)
	}

	refreshExpires := now.Add(refreshTime)
	refreshToken, err := util.GenerateToken(jwt.MapClaims{
		util.Jti: tokenID,
		util.Typ: util.RefreshType,
		util.Iat: now.Unix(),
		util.Exp: refreshExpires.Unix(),
		util.Sid: sessionID,
		util.Fgp: fingerprint,
		util.Sub: userID,
	}, []byte(b.config.Secret))

	return refreshToken, refreshExpires, nil
}

func (b *Business) createUserSession(ctx context.Context, fg *fingerprint.Fingerprint, userID, sessionID string, now, refreshExpires time.Time) error {
	session := &model.Session{
		BaseModel: model.BaseModel{
			ID:        sessionID,
			CreatedAt: now,
			UpdatedAt: now,
		},
		LastSeenAt:  &now,
		UserID:      userID,
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

func (b *Business) generateUserToken(ctx context.Context, fg *fingerprint.Fingerprint, user *model.User, sessionID string) (accessToken string, tokenExpires time.Time, refreshToken string, refreshExpires time.Time, err error) {
	now := time.Now().UTC()
	tokenID := uuid.NewString()

	accessToken, tokenExpires, err = b.generateAccessToken(user, tokenID, sessionID, fg.ID, now)
	if err != nil {
		return
	}

	refreshToken, refreshExpires, err = b.generateRefreshToken(user.ID, tokenID, sessionID, fg.ID, now)
	if err != nil {
		return
	}

	err = b.createUserSession(ctx, fg, user.ID, sessionID, now, refreshExpires)
	return
}
