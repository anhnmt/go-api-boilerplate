package otel

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
	"go.uber.org/fx"

	"github.com/anhnmt/go-api-boilerplate/internal/pkg/config"
)

type Params struct {
	fx.In

	Ctx       context.Context
	AppConfig config.App
	Config    Config
}

func New(lc fx.Lifecycle, p Params) error {
	if p.Config.Endpoint == "" {
		return nil
	}

	// Ensure default SDK resources and the required service name are set.
	r, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(p.AppConfig.Name),
		),
	)
	if err != nil {
		return fmt.Errorf("failed to create resource: %v", err)
	}

	opts := []sdktrace.TracerProviderOption{
		sdktrace.WithResource(r),
	}

	exporter, err := newExporter(p)
	if err != nil {
		return fmt.Errorf("failed to create exporter: %v", err)
	}

	opts = append(opts, exporter)

	// Create a simple span processor that writes to the exporter
	tp := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(tp)

	// Set the global propagator to use W3C Trace Context
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	lc.Append(fx.StopHook(tp.Shutdown))
	return nil
}

func newExporter(p Params) (sdktrace.TracerProviderOption, error) {
	if p.Config.Type == "grpc" {
		exporter, err := otlptracegrpc.New(
			p.Ctx,
			otlptracegrpc.WithEndpoint(p.Config.Endpoint),
			otlptracegrpc.WithInsecure(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create exporter: %v", err)
		}

		return sdktrace.WithBatcher(exporter), nil
	}

	if p.Config.Type == "http" {
		exporter, err := otlptracehttp.New(
			p.Ctx,
			otlptracehttp.WithEndpoint(p.Config.Endpoint),
			otlptracehttp.WithInsecure(),
		)
		if err != nil {
			return nil, fmt.Errorf("failed to create exporter: %v", err)
		}

		return sdktrace.WithBatcher(exporter), nil
	}

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		return nil, fmt.Errorf("failed to create exporter: %v", err)
	}

	return sdktrace.WithBatcher(exporter), nil
}