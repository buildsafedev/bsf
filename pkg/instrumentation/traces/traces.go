package traces

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"

	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc/credentials"
)

// localExporter returns a console exporter.
func localExporter(w io.Writer) (trace.SpanExporter, error) {
	return stdouttrace.New(
		stdouttrace.WithWriter(w),

		// Use human readable output.
		stdouttrace.WithPrettyPrint(),

		// Do not print timestamps for the demo.
		stdouttrace.WithoutTimestamps(),
	)
}

// SetupTracing sets up the tracing provider
func SetupTracing(r *resource.Resource) *trace.TracerProvider {
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if collectorURL == "" {
		return localTraceProvider(r)
	}

	return grpcTracer(collectorURL, r)
}

func localTraceProvider(r *resource.Resource) *trace.TracerProvider {
	// Create file if it doesn't exist
	f, err := os.Create("trace.json")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	exp, err := localExporter(f)
	if err != nil {
		log.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exp),
		trace.WithResource(r),
	)

	return tp
}

func grpcTracer(collectorURL string, r *resource.Resource) *trace.TracerProvider {
	secureOption := otlptracegrpc.WithTLSCredentials(credentials.NewClientTLSFromCert(nil, ""))
	if len(os.Getenv("OTEL_INSECURE_MODE")) > 0 {
		secureOption = otlptracegrpc.WithInsecure()
	}

	exporter, err := otlptrace.New(
		context.Background(),
		otlptracegrpc.NewClient(
			secureOption,
			otlptracegrpc.WithEndpoint(collectorURL),
		),
	)

	if err != nil {
		fmt.Printf("failed to create trace exporter: %v", err)
	}

	tp := trace.NewTracerProvider(
		trace.WithSampler(trace.AlwaysSample()),
		trace.WithBatcher(exporter),
		trace.WithResource(r))

	return tp
}
