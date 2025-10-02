package observability

import (
	"context"
	"errors"
	"time"

	"github.com/ouz/goauthboilerplate/internal/config"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func SetupOTelSDK(ctx context.Context, logger *config.Logger) (shutdown func(context.Context) error, err error) {
	cfg := config.Get()

	if !cfg.Otel.MonitoringEnabled {
		logger.Info("Monitoring is disabled, skipping OpenTelemetry setup")
		return func(ctx context.Context) error { return nil }, nil
	}

	var shutdownFuncs []func(context.Context) error

	shutdown = func(ctx context.Context) error {
		var err error
		for _, fn := range shutdownFuncs {
			err = errors.Join(err, fn(ctx))
		}
		shutdownFuncs = nil
		return err
	}

	handleErr := func(inErr error) {
		err = errors.Join(inErr, shutdown(ctx))
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.Otel.ServiceName),
			semconv.ServiceVersion("1.0.0"),
			semconv.ServiceNamespace("production"),
		),
	)
	if err != nil {
		logger.Warn("Failed to create OpenTelemetry resource", "error", err)
		return func(ctx context.Context) error { return nil }, nil
	}

	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient())
	if err != nil {
		logger.Warn("Failed to create trace exporter, monitoring will be disabled", "error", err)
		return func(ctx context.Context) error { return nil }, nil
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	metricExporter, err := otlpmetrichttp.New(ctx)
	if err != nil {
		logger.Warn("Failed to create metric exporter, metrics will be disabled", "error", err)
		return func(ctx context.Context) error { return nil }, nil
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)
	if err != nil {
		handleErr(err)
		return
	}
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	logExporter, err := otlploghttp.New(ctx)
	if err != nil {
		logger.Warn("Failed to create log exporter, logs will be disabled", "error", err)
		return func(ctx context.Context) error { return nil }, nil
	}

	logProvider := log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(logExporter)),
		log.WithResource(res),
	)
	shutdownFuncs = append(shutdownFuncs, logProvider.Shutdown)

	global.SetLoggerProvider(logProvider)

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		logger.Warn("Failed to start runtime instrumentation", "error", err)
		// Don't fatal here, just continue without runtime metrics
	}

	return
}
