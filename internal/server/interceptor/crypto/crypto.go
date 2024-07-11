package cryptointerceptor

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"

	"github.com/bytedance/sonic"
	"github.com/samber/lo"
	"google.golang.org/grpc/codes"

	"github.com/anhnmt/go-api-boilerplate/internal/common"
	"github.com/anhnmt/go-api-boilerplate/internal/common/cryptoutils"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
)

const XRequestKey = "X-Request-Key"

var defaultGuardLists = []string{
	"/auth.v1.AuthService/Encrypt",
}

type CryptoInterceptor interface {
	Handler(http.Handler) http.Handler
}

type cryptoInterceptor struct {
	cfg config.Crypto
}

func New(cfg config.Crypto) CryptoInterceptor {
	return &cryptoInterceptor{
		cfg: cfg,
	}
}

func (c *cryptoInterceptor) Handler(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !lo.Contains(defaultGuardLists, r.URL.Path) {
			h.ServeHTTP(w, r)
			return
		}

		// handler decrypt data
		requestKey := r.Header.Get(XRequestKey)
		if requestKey == "" {
			writeErrorResponse(w, codes.InvalidArgument, "request key not found")
			return
		}

		// xChecksum := r.Header.Get("X-Checksum")
		// if xChecksum == "" {
		//     writeErrorResponse(w, codes.InvalidArgument, "checksum not found")
		//     return
		// }

		rawRequestKey, err := cryptoutils.DecryptRSAString(requestKey, c.cfg.PrivateKeyBytes())
		if err != nil {
			writeErrorResponse(w, codes.Internal, "fail to decrypt request key")
			return
		}

		// Decrypt request body
		var payload CryptoData
		if err = sonic.ConfigDefault.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeErrorResponse(w, codes.Internal, "fail to marshal request payload")
			return
		}

		decryptedBody, err := cryptoutils.DecryptAES(payload.Data, rawRequestKey)
		if err != nil {
			writeErrorResponse(w, codes.Internal, "fail to decrypt request body")
			return
		}

		// rewrite request body with decrypted data
		buf := bytes.NewBuffer(decryptedBody)
		r.Body = io.NopCloser(buf)
		r.ContentLength = int64(buf.Len())

		// Wrap response write by ResponseWriter
		responseWriter := httptest.NewRecorder()
		// Serve request through server mux
		h.ServeHTTP(responseWriter, r)

		// Copy responseWriter.buf to w (default response writer)
		cipherText, err := cryptoutils.EncryptAES(responseWriter.Body.Bytes(), rawRequestKey)
		if err != nil {
			writeErrorResponse(w, codes.Internal, "fail to encrypt response body")
			return
		}

		w.Header().Set("Content-Type", "application/json")
		err = sonic.ConfigDefault.NewEncoder(w).Encode(&CryptoData{
			Data: cipherText,
		})
		if err != nil {
			writeErrorResponse(w, codes.Internal, "fail to write response body")
			return
		}
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
