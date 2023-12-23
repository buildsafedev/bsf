package instrumentation

import (
	"fmt"
	"sync"

	"go.opentelemetry.io/otel/metric"
)

// Chose metric type based on the use case:
// https://prometheus.io/docs/concepts/metric_types/

const (
	metricKeyFormat   = "rpc.%s"
	requestCounter    = "request_counter"
	messageKey        = "message"
	unitDimensionless = "1"
	unitBytes         = "By"
	unitMilliseconds  = "ms"
)

type metricInstruments struct {
	initOnce sync.Once
	initErr  error

	// Define custom metrics for the service here.
	RegisteredUsers metric.Int64Counter
}

func (i *metricInstruments) init(meter metric.Meter) *metricInstruments {
	i.initOnce.Do(func() {
		i.RegisteredUsers, i.initErr = meter.Int64Counter(
			formatkeys(requestCounter),
			metric.WithUnit(unitDimensionless),
		)
		if i.initErr != nil {
			return
		}

	})
	return i
}

func newMetrics(meter metric.Meter) *metricInstruments {
	i := metricInstruments{}
	return i.init(meter)
}

func formatkeys(metricName string) string {
	return fmt.Sprintf(metricKeyFormat, metricName)
}
