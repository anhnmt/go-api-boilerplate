package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"time"

	"connectrpc.com/vanguard"
	"connectrpc.com/vanguard/vanguardgrpc"
	"github.com/rs/cors"
	"github.com/rs/zerolog/log"
	"go.uber.org/fx"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"google.golang.org/grpc"
	"google.golang.org/grpc/encoding"
	"google.golang.org/protobuf/encoding/protojson"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
	cryptointerceptor "github.com/anhnmt/go-api-boilerplate/internal/server/interceptor/crypto"
)

var opts = []vanguard.TranscoderOption{
	vanguard.WithDefaultServiceOptions(
		vanguard.WithTargetProtocols(
			vanguard.ProtocolGRPC,
			vanguard.ProtocolGRPCWeb,
		),
	),
}

func init() {
	encoding.RegisterCodec(vanguardgrpc.NewCodec(&vanguard.JSONCodec{
		// These fields can be used to customize the serialization and
		// de-serialization behavior. The options presented below are
		// highly recommended.
		MarshalOptions: protojson.MarshalOptions{
			EmitUnpopulated: false,
			UseProtoNames:   true,
		},
		UnmarshalOptions: protojson.UnmarshalOptions{
			DiscardUnknown: true,
		},
	}))
}

type Params struct {
	fx.In

	Config            config.Server
	GrpcServer        *grpc.Server
	CryptoInterceptor cryptointerceptor.CryptoInterceptor
}

type Server struct {
	config config.Server
	mux    http.Handler
	srv    *http.Server
}

func New(lc fx.Lifecycle, p Params) (*Server, error) {
	var mux http.Handler
	var err error
	mux, err = vanguardgrpc.NewTranscoder(p.GrpcServer, opts...)
	if err != nil {
		return nil, err
	}

	// Add encrypt interceptor
	mux = p.CryptoInterceptor.Handler(mux)

	// Add CORS support
	mux = withCORS(mux)

	srv := &Server{
		mux:    mux,
		config: p.Config,
	}

	lc.Append(fx.StartStopHook(
		srv.Start,
		srv.Stop,
	))

	return srv, nil
}

func (s *Server) Start() error {
	errChan := make(chan error)

	if s.config.Pprof.Enable {
		go func(errChan chan error) {
			addr := fmt.Sprintf(":%d", s.config.Pprof.Port)
			log.Info().Msgf("Starting pprof http://localhost%s", addr)

			errChan <- http.ListenAndServe(addr, http.DefaultServeMux)
		}(errChan)
	}

	// create new http Server
	addr := fmt.Sprintf(":%d", s.config.Grpc.Port)
	s.srv = &http.Server{
		Addr: addr,
		// We use the h2c package in order to support HTTP/2 without TLS,
		// so we can handle gRPC requests, which requires HTTP/2, in
		// addition to Connect and gRPC-Web (which work with HTTP 1.1).
		Handler: h2c.NewHandler(
			s.mux,
			&http2.Server{},
		),
	}

	// Serve the http Server on the http listener.
	go func(errChan chan error) {
		log.Info().Msgf("Starting application http://localhost%s", addr)

		// run the Server
		err := s.srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			errChan <- err
		}
	}(errChan)

	select {
	case err := <-errChan:
		if err != nil {
			return err
		}
	case <-time.After(100 * time.Millisecond):
		return nil
	}

	return nil
}

func (s *Server) Stop(ctx context.Context) error {
	if s.srv != nil {
		return s.srv.Shutdown(ctx)
	}

	return nil
}

// withCORS adds CORS support to a gRPC HTTP handler.
func withCORS(h http.Handler) http.Handler {
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"}, // replace with your domain
		AllowedMethods: []string{
			http.MethodPost, // for all protocols
		},
		AllowedHeaders: []string{
			"Content-Type",                // for all protocols
			"Connect-Protocol-Version",    // for Connect
			"Connect-Timeout-Ms",          // for Connect
			"Grpc-Timeout",                // for gRPC-web
			"X-Grpc-Web",                  // for gRPC-web
			"X-User-Agent",                // for all protocols
			cryptointerceptor.XRequestKey, // for encrypt interceptor
		},
		ExposedHeaders: []string{
			"Grpc-Status",             // for gRPC-web
			"Grpc-Message",            // for gRPC-web
			"Grpc-Status-Details-Bin", // for gRPC-web
		},
		MaxAge: 7200, // 2 hours in seconds
	})
	return c.Handler(h)
}
