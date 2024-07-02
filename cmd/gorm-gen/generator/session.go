package generator

import (
	"github.com/anhnmt/go-api-boilerplate/proto/pb"
)

// Dynamic SQL
type Session interface {
	// select id, fingerprint, user_agent, os, device_type, device, ip_address, created_at as login_time, last_seen_at as last_seen
	// {{if sessionId != ""}}
	// , CASE
	//	WHEN id = @sessionId THEN true
	//	ELSE false
	// END as is_current
	// {{end}}
	// from sessions
	// where user_id = @userId
	// and is_revoked = false
	// and expires_at >= NOW() - INTERVAL '24 hours'
	// order by
	// {{if sessionId != ""}}
	// is_current DESC,
	// {{end}}
	// last_seen_at DESC, updated_at DESC, expires_at DESC;
	FindByUserIdAndSessionId(userId, sessionId string) ([]*pb.ActiveSessions, error)
}
