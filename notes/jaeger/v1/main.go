package main

import (
	"context"
	"io"
	"log"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	trace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

const name = "some-service"

// span1
//
//	-> span2
//	-> span3
//			-> some3Child
func some1(ctx context.Context) {
	ctx, span := otel.Tracer(name).Start(ctx, "some1")
	defer span.End()
	some2(ctx)
}

func some2(ctx context.Context) {
	_, span := otel.Tracer(name).Start(ctx, "some2")
	defer span.End()
	some3(ctx)
}

func some3(ctx context.Context) {
	ctx, span := otel.Tracer(name).Start(ctx, "some3")
	defer span.End()
	span.SetAttributes(attribute.String("request.n", "abc"))
	span.SetStatus(codes.Error, "invalid request")
	func(ctx context.Context) {
		_, span := otel.Tracer(name).Start(ctx, "some3Child")
		span.SetStatus(codes.Ok, "over")
		span.End()
	}(ctx)
}

func main() {
	jaegerExporter, err := JaegerExporter()
	if err != nil {
		log.Fatalf("new exporter error: %v", err)
	}
	tp := trace.NewTracerProvider(
		trace.WithBatcher(jaegerExporter),
		trace.WithResource(newResource()),
	)

	otel.SetTracerProvider(tp)
	// ----------------
	ctx := context.Background()
	some1(ctx)
	time.Sleep(time.Second * 6)
}

func newExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),
		// Use human-readable output.
		stdouttrace.WithPrettyPrint(),
		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

func newResource() *resource.Resource {
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(name),
			semconv.ServiceVersion("v0.1.0"),
			attribute.String("environment", "demo"),
		),
	)
	return r
}

func JaegerExporter() (*jaeger.Exporter, error) {
	endpoint := jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces"))
	exporter, err := jaeger.New(endpoint)
	if err != nil {
		return nil, err
	}
	return exporter, nil
}
