package main

import (
	"context"
	"fmt"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	_ "go.uber.org/automaxprocs"
	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/cmd/api-server/config"
	"github.com/anhnmt/go-api-boilerplate/gen/gormgen"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/base"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/logger"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/postgres"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/redis"
	"github.com/anhnmt/go-api-boilerplate/internal/pkg/util"
	"github.com/anhnmt/go-api-boilerplate/internal/server"
	"github.com/anhnmt/go-api-boilerplate/internal/server/grpc"
	"github.com/anhnmt/go-api-boilerplate/internal/server/interceptor"
	"github.com/anhnmt/go-api-boilerplate/internal/service"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	app := fx.New(
		fx.WithLogger(logger.NewFxLogger),
		fx.Provide(
			util.ProvideCtx(ctx),
			gormgen.Use,
		),
		fx.Invoke(
			base.AutoMigrate,
		),
		config.Module,
		logger.Module,
		postgres.Module,
		redis.Module,
		interceptor.Module,
		grpc.Module,
		service.Module,
		server.Module,
	)

	// Initialize tracing and handle the tracer provider shutdown
	stopTracing, err := initTracing(ctx)
	if err != nil {
		panic(fmt.Errorf("failed to initialize tracing: %w", err))
	}

	defer stopTracing()

	if err := app.Start(ctx); err != nil {
		panic(fmt.Errorf("failed to start application: %w", err))
	}

	<-app.Done()

	if err := app.Stop(ctx); err != nil {
		panic(fmt.Errorf("failed to stop application: %w", err))
	}

	log.Info().Msg("Gracefully shutting down")
}

// Initialize OpenTelemetry tracing and return a function to stop the tracer provider
func initTracing(ctx context.Context) (func() error, error) {
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %v", err)
	}

	// grpcExporter, err := otlptracegrpc.New(
	//     ctx,
	//     // otlptracegrpc.WithEndpoint("localhost:4317"),
	//     otlptracegrpc.WithEndpointURL("http://localhost:4318/v1/traces"),
	//     otlptracegrpc.WithInsecure(),
	// )

	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %v", err)
	}

	// Create a simple span processor that writes to the exporter
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		// sdktrace.WithBatcher(grpcExporter),
	)
	otel.SetTracerProvider(tp)

	// Set the global propagator to use W3C Trace Context
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	// Return a function to stop the tracer provider
	return func() error {
		if err := tp.Shutdown(context.Background()); err != nil {
			return fmt.Errorf("failed to shut down tracer provider: %v", err)
		}

		return nil
	}, nil
}
