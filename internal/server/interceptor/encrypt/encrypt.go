package encryptinterceptor

import (
	"net/http"

	"github.com/bytedance/sonic"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"

	"github.com/anhnmt/go-api-boilerplate/internal/common"
)

var defaultGuardLists = []string{
	"/auth.v1.AuthService/Encrypt",
}

type EncryptInterceptor interface {
	Handler(http.Handler) http.Handler
}

type encryptInterceptor struct {
}

func New() EncryptInterceptor {
	return &encryptInterceptor{}
}

func (e *encryptInterceptor) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !lo.Contains(defaultGuardLists, r.URL.Path) {
			h.ServeHTTP(w, r)
			return
		}

		// handler decrypt data
		requestKey := r.Header.Get("X-Request-Key")
		if requestKey == "" {
			writeErrorResponse(w, codes.InvalidArgument, "request key not found")
			return
		}

		checksum := r.Header.Get("X-Checksum")
		if checksum == "" {
			writeErrorResponse(w, codes.InvalidArgument, "checksum not found")
			return
		}

		h.ServeHTTP(w, r)
	})
}

func writeErrorResponse(w http.ResponseWriter, c codes.Code, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(common.HTTPStatusFromCode(c))

	_ = sonic.ConfigDefault.NewEncoder(w).Encode(&ErrorResponse{
		Code:    common.StringFromCode(c),
		Message: msg,
	})
}
