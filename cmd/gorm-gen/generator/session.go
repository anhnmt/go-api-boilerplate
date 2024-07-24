package generator

import (
	"github.com/anhnmt/go-api-boilerplate/gen/pb"
)

// Dynamic SQL
type Session interface {
	// select id,
	// fingerprint,
	// user_agent,
	// os,
	// device_type,
	// device,
	// ip_address,
	// created_at as login_time,
	// last_seen_at as last_seen
	// {{if sessionID != ""}}
	// , CASE
	//	WHEN id = @sessionID THEN true
	//	ELSE false
	// END as is_current
	// {{end}}
	// from sessions
	// where user_id = @userID
	// and is_revoked = false
	// and expires_at >= NOW() - INTERVAL '24 hours'
	// order by
	// {{if sessionID != ""}}
	// is_current DESC,
	// {{end}}
	// last_seen_at DESC, updated_at DESC, expires_at DESC
	// LIMIT {{if limit == 0}} 10 {{else}} @limit {{end}}
	// OFFSET @offset;
	FindByUserIDAndSessionID(userID, sessionID string, limit, offset int) ([]*pb.ActiveSessions, error)

	// select count(1) AS total
	// from sessions
	// where user_id = @userID
	// and is_revoked = false
	// and expires_at >= NOW() - INTERVAL '24 hours'
	CountByUserID(userID string) (int, error)

	// update sessions
	// set is_revoked = true
	// where user_id = @userID
	// {{if sessionID != ""}}
	// and id <> @sessionID
	// {{end}}
	// and is_revoked = false
	// and expires_at >= NOW() - INTERVAL '24 hours'
	UpdateRevokedByUserIDWithoutSessionID(userID, sessionID string) error

	// select id
	// from sessions
	// where user_id = @userID
	// {{if sessionID != ""}}
	// and id <> @sessionID
	// {{end}}
	// and is_revoked = false
	// and expires_at >= NOW() - INTERVAL '24 hours'
	FindByUserIDWithoutSessionID(userID, sessionID string) ([]string, error)
}
