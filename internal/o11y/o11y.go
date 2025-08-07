// Package o11y provides observability (tracing, metrics, logging) setup for the application.
// Package o11y implementa instrumentação e tracing com OpenTelemetry.
package o11y

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
)

func Tracer() trace.Tracer {
	return otel.Tracer("chi-app")
}

func InitTracer(ctx context.Context) func(context.Context) error {
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithInsecure(),
		otlptracehttp.WithEndpoint("otel-collector:4318"),
	)
	if err != nil {
		log.Fatal().Err(err).Msg("Erro ao inicializar OTLP exporter")
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("chi-app"),
		)),
	)

	otel.SetTracerProvider(tp)
	return tp.Shutdown
}
