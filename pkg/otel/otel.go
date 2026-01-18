package otel

import (
	"context"
	"errors"
	stdlog "log"
	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
)

var Tracer trace.Tracer

type OtelConfig struct {
	ServiceName       string
	ServiceVersion    string
	ServiceNamespace  string
	ExporterEndpoint  string
	MonitoringEnabled bool
}

func SetupOTelSDK(ctx context.Context, cfg OtelConfig) (shutdown func(context.Context) error, err error) {
	if !cfg.MonitoringEnabled {
		stdlog.Println("[otel] monitoring disabled, skipping OpenTelemetry setup")
		Tracer = noop.NewTracerProvider().Tracer(cfg.ServiceName)
		return func(ctx context.Context) error { return nil }, nil
	}

	if cfg.ServiceVersion == "" {
		cfg.ServiceVersion = "1.0.0"
	}
	if cfg.ServiceNamespace == "" {
		cfg.ServiceNamespace = "production"
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

	// Create resource
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.ServiceVersion),
			semconv.ServiceNamespace(cfg.ServiceNamespace),
		),
	)
	if err != nil {
		return nil, err
	}

	// Set up propagator
	prop := propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(prop)

	// Parse endpoint to decide secure/insecure and strip scheme
	endpoint := cfg.ExporterEndpoint
	isInsecure := true
	if len(endpoint) > 8 && endpoint[:8] == "https://" {
		endpoint = endpoint[8:]
		isInsecure = false
	} else if len(endpoint) > 7 && endpoint[:7] == "http://" {
		endpoint = endpoint[7:]
	}

	// Set up trace exporter
	traceOpts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(endpoint),
	}
	if isInsecure {
		traceOpts = append(traceOpts, otlptracehttp.WithInsecure())
	}
	traceExporter, err := otlptrace.New(ctx, otlptracehttp.NewClient(traceOpts...))
	if err != nil {
		return nil, err
	}

	// Set up trace provider with 25% sampling rate
	tracerProvider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(traceExporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.TraceIDRatioBased(0.25)),
	)
	shutdownFuncs = append(shutdownFuncs, tracerProvider.Shutdown)
	otel.SetTracerProvider(tracerProvider)

	// Set up metric exporter
	metricOpts := []otlpmetrichttp.Option{
		otlpmetrichttp.WithEndpoint(endpoint),
	}
	if isInsecure {
		metricOpts = append(metricOpts, otlpmetrichttp.WithInsecure())
	}
	metricExporter, err := otlpmetrichttp.New(ctx, metricOpts...)
	if err != nil {
		handleErr(err)
		return
	}

	// Set up meter provider
	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)
	shutdownFuncs = append(shutdownFuncs, meterProvider.Shutdown)
	otel.SetMeterProvider(meterProvider)

	logOpts := []otlploghttp.Option{
		otlploghttp.WithEndpoint(endpoint),
	}
	if isInsecure {
		logOpts = append(logOpts, otlploghttp.WithInsecure())
	}
	logExporter, err := otlploghttp.New(ctx, logOpts...)
	if err != nil {
		stdlog.Printf("[otel] failed to create log exporter, logs will only go to stdout: %v", err)
	} else {
		logProvider := sdklog.NewLoggerProvider(
			sdklog.WithProcessor(sdklog.NewBatchProcessor(logExporter)),
			sdklog.WithResource(res),
		)
		shutdownFuncs = append(shutdownFuncs, logProvider.Shutdown)
		global.SetLoggerProvider(logProvider)
	}

	err = runtime.Start(runtime.WithMinimumReadMemStatsInterval(time.Second))
	if err != nil {
		stdlog.Fatal(err)
	}

	Tracer = otel.Tracer(cfg.ServiceName)
	stdlog.Printf("[otel] OpenTelemetry SDK initialized successfully for service: %s", cfg.ServiceName)
	return
}

func WrapHandlerWithTelemetry(handler http.Handler, serviceName string) http.Handler {
	skipPaths := map[string]bool{
		"/":        true,
		"/ready":   true,
		"/health":  true,
		"/ping":    true,
		"/metrics": true,
	}

	filter := func(req *http.Request) bool {
		return !skipPaths[req.URL.Path]
	}
	return otelhttp.NewHandler(handler, serviceName,
		otelhttp.WithFilter(filter),
		otelhttp.WithSpanOptions(trace.WithSpanKind(trace.SpanKindServer)),
	)
}

func NewHTTPClientWithTelemetry() *http.Client {
	return &http.Client{
		Transport: otelhttp.NewTransport(http.DefaultTransport),
	}
}

func GetTracer() trace.Tracer {
	return Tracer
}
