package metrics

import (
	"log"

	"go.opentelemetry.io/otel/exporters/prometheus"
	metricsdk "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
)

// SetupMetrics sets up the metrics provider
func SetupMetrics(r *resource.Resource) (metricsdk.Reader, *metricsdk.MeterProvider) {
	exporter, err := prometheus.New()
	if err != nil {
		log.Fatal(err)
	}
	provider := metricsdk.NewMeterProvider(metricsdk.WithReader(exporter), metricsdk.WithResource(r))
	return exporter, provider
}
