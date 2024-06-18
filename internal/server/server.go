package server

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	_ "net/http/pprof"

	"connectrpc.com/vanguard"
	"github.com/rs/zerolog/log"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
	"golang.org/x/sync/errgroup"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
)

type Server interface {
	Start(context.Context, config.Server) error
}

type server struct {
	mux *http.ServeMux
}

func New(mux *http.ServeMux, services []*vanguard.Service) (Server, error) {
	opts := []vanguard.TranscoderOption{
		vanguard.WithDefaultServiceOptions(
			vanguard.WithTargetProtocols(
				vanguard.ProtocolConnect,
				vanguard.ProtocolGRPC,
				vanguard.ProtocolGRPCWeb,
			),
		),
	}

	transcoder, err := vanguard.NewTranscoder(services, opts...)
	if err != nil {
		return nil, err
	}

	mux.Handle("/", transcoder)

	return &server{
		mux: mux,
	}, nil
}

func (s *server) Start(ctx context.Context, cfg config.Server) error {
	g, _ := errgroup.WithContext(ctx)

	if cfg.Pprof.Enable {
		g.Go(func() error {
			addr := fmt.Sprintf(":%d", cfg.Pprof.Port)
			log.Info().Msgf("Starting pprof http://localhost%s", addr)

			return http.ListenAndServe(addr, http.DefaultServeMux)
		})
	}

	// Serve the http server on the http listener.
	g.Go(func() error {
		addr := fmt.Sprintf(":%d", cfg.Grpc.Port)
		log.Info().Msgf("Starting application http://localhost%s", addr)

		// create new http server
		srv := &http.Server{
			Addr: addr,
			// We use the h2c package in order to support HTTP/2 without TLS,
			// so we can handle gRPC requests, which requires HTTP/2, in
			// addition to Connect and gRPC-Web (which work with HTTP 1.1).
			Handler: h2c.NewHandler(
				s.mux,
				&http2.Server{},
			),
		}

		defer func() {
			_ = srv.Close()
		}()

		// run the server
		err := srv.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}

		return nil
	})

	return g.Wait()
}
