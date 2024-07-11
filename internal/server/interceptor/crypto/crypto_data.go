package cryptointerceptor

import (
	"bytes"
	"net/http"
)

type CryptoData struct {
	Data string `json:"data"`
}

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ResponseWriter wrap http response to response body
type ResponseWriter struct {
	http.ResponseWriter
	Buf *bytes.Buffer
}

func (rw *ResponseWriter) Write(p []byte) (int, error) {
	return rw.Buf.Write(p)
}

func (rw *ResponseWriter) Bytes() []byte {
	return rw.Buf.Bytes()
}

func (rw *ResponseWriter) Flush() {
	rw.Buf.Reset()
	if fw, ok := rw.ResponseWriter.(http.Flusher); ok {
		fw.Flush()
	}
}
