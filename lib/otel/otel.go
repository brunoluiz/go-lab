package otel

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

// SetupOTelSDK bootstraps the OpenTelemetry pipeline.
// If it does not return an error, make sure to call shutdown for proper cleanup.
func SetupOTelSDK(ctx context.Context) (func(context.Context) error, error) {
	var shutdownFuncs []func(context.Context) error
	var err error

	shutdown := func(ctx context.Context) error {
		var shutdownErr error
		for _, fn := range shutdownFuncs {
			shutdownErr = errors.Join(shutdownErr, fn(ctx))
		}
		shutdownFuncs = nil
		return shutdownErr
	}

	// Set up propagator.
	prop := newPropagator()
	otel.SetTextMapPropagator(prop)

	// Set up trace provider.
	tracerProvider, err := newTracerProvider(ctx)
	if err != nil {
		return shutdown, errors.Join(err, shutdown(ctx))
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up meter provider.
	meterProvider, err := newMeterProvider(ctx)
	if err != nil {
		return shutdown, errors.Join(err, shutdown(ctx))
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	return shutdown, err
}

func newPropagator() propagation.TextMapPropagator {
	return propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
}

func newTracerProvider(ctx context.Context) (*sdktrace.TracerProvider, error) {
	exporter := os.Getenv("OTEL_TRACES_EXPORTER")
	if exporter == "" {
		exporter = "otlp"
	}

	var exp sdktrace.SpanExporter
	var err error

	switch exporter {
	case "otlp":
		exp, err = otlptracehttp.New(ctx)
	default:
		return nil, fmt.Errorf("unsupported traces exporter: %s", exporter)
	}

	if err != nil {
		return nil, err
	}

	res := resource.Default()

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	return tp, nil
}

func newMeterProvider(ctx context.Context) (*metric.MeterProvider, error) {
	exporter := os.Getenv("OTEL_METRICS_EXPORTER")
	if exporter == "" {
		exporter = "otlp"
	}

	var metricExporter metric.Reader

	switch exporter {
	case "console":
		exp, err := stdoutmetric.New()
		if err != nil {
			return nil, err
		}
		metricExporter = metric.NewPeriodicReader(exp, metric.WithInterval(3*time.Second))
	case "otlp":
		exp, err := otlpmetrichttp.New(ctx)
		if err != nil {
			return nil, err
		}
		metricExporter = metric.NewPeriodicReader(exp, metric.WithInterval(3*time.Second))
	default:
		return nil, fmt.Errorf("unsupported metrics exporter: %s", exporter)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metricExporter),
	)
	return meterProvider, nil
}
