package otelinterceptor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	grpcMetadata "google.golang.org/grpc/metadata"
)

// UnaryServerInterceptor returns a new unary server interceptor that extracts the real client IP from request headers.
// It checks if the request comes from a trusted peer, validates headers against trusted proxies list and trusted proxies count
// then it extracts the IP from the configured headers.
// The real IP is added to the request context.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, _ *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
		md := traceResponseFromContext(ctx)
		if md.Len() == 0 {
			return handler(ctx, req)
		}

		err := grpc.SetHeader(ctx, md)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}

// StreamServerInterceptor returns a new stream server interceptor that extracts the real client IP from request headers.
// It checks if the request comes from a trusted peer, validates headers against trusted proxies list and trusted proxies count
// then it extracts the IP from the configured headers.
// The real IP is added to the request context.
func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv any, stream grpc.ServerStream, _ *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := stream.Context()
		md := traceResponseFromContext(ctx)
		if md.Len() == 0 {
			return handler(srv, stream)
		}

		err := grpc.SetHeader(ctx, md)
		if err != nil {
			return err
		}

		return handler(srv, stream)
	}
}

// traceResponseFromContext extracts the trace from the context.
// https://uptrace.dev/opentelemetry/opentelemetry-traceparent.html#what-is-traceparent-header
// # {version}-{trace_id}-{span_id}-{trace_flags}
// The RayId header uses the version-trace_id-parent_id-trace_flags format where:
// version is always 00.
// trace_id is a hex-encoded trace id.
// span_id is a hex-encoded span id.
// trace_flags is a hex-encoded 8-bit field that contains tracing flags such as sampling, trace level, etc.
// RayId: 00-80e1afed08e019fc1110464cfa66635c-7a085853722dc6d2-01
func traceResponseFromContext(ctx context.Context) grpcMetadata.MD {
	span := trace.SpanContextFromContext(ctx)

	if !span.IsValid() {
		return nil
	}

	xRayId := fmt.Sprintf(
		"00-%s-%s-%s",
		span.TraceID().String(),
		span.SpanID().String(),
		span.TraceFlags().String(),
	)

	return grpcMetadata.Pairs("X-Ray-Id", xRayId)
}
