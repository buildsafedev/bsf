package instrumentation

import (
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/trace"
)

// Instrumentor will have all the configured clients to enable instrumentation for the service
// This is expected to be used by all the handlers in the service
type Instrumentor struct {
	Tracer  trace.Tracer
	Metrics *metricInstruments
}

const (
	// TODO: configure such that version is set at build time
	version = "v0.0.1"
	svcName = "buildsafe/v1"
)

// SVCResource returns the resource for the service
func SVCResource() *resource.Resource {
	return resource.NewSchemaless(
		attribute.String("service.name", svcName),
	)
}

// SetupInstrumentation sets up the instrumentation for the service
func SetupInstrumentation(meterProvider metric.MeterProvider, traceProvider trace.TracerProvider) Instrumentor {
	meter := meterProvider.Meter(svcName, metric.WithInstrumentationVersion(version))
	metrics := newMetrics(meter)

	tracer := traceProvider.Tracer(svcName, trace.WithInstrumentationVersion(version))

	return Instrumentor{
		Tracer:  tracer,
		Metrics: metrics,
	}

}
