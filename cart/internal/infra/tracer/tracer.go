package tracer

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.37.0"
	t "go.opentelemetry.io/otel/trace"
)

type TracerManager struct {
	Tracer         t.Tracer
	TracerProvider *trace.TracerProvider
}

func NewTracer(ctx context.Context, endpoint string) (*TracerManager, error) {
	exporter, err := otlptracehttp.New(ctx, otlptracehttp.WithEndpointURL(endpoint))
	if err != nil {
		return nil, err
	}
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName("cart"),
			semconv.DeploymentEnvironmentName("development"),
			semconv.URLFull("jaeger"),
		),
	)
	if err != nil {
		return nil, err
	}
	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(exporter),
		trace.WithResource(res),
	)

	otel.SetTracerProvider(tracerProvider)

	tracer := otel.GetTracerProvider().Tracer("cart")
	return &TracerManager{
		Tracer:         tracer,
		TracerProvider: tracerProvider,
	}, nil
}
