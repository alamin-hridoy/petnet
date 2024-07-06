package middleware

import (
	"context"
	"strings"
	"time"

	"github.com/grpc-ecosystem/go-grpc-prometheus/packages/grpcstatus"
	prom "github.com/prometheus/client_golang/prometheus"
	"google.golang.org/grpc"
)

type ServerMetrics struct {
	labels                 []string
	serverHandledCounter   *prom.CounterVec
	serverHandledHistogram *prom.HistogramVec
}

// NewServerMetrics returns a ServerMetric which exposes the grpc service metrics for prometheus.
// SeverMetricLabels should contain the name for the custom labels that we want to attach to all the
// metrics.
func NewServerMetrics(labelExtractor LabelExtractor) *ServerMetrics {
	labels := append([]string{"grpc_service", "grpc_method", "grpc_status"}, labelExtractor.LabelNames()...)
	return &ServerMetrics{
		labels: labels,
		serverHandledCounter: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "grpc_server_handled_total",
				Help: "Total number of RPCs completed on the server, regardless of success or failure.",
			}, labels,
		),
		serverHandledHistogram: prom.NewHistogramVec(
			prom.HistogramOpts{
				Name:    "grpc_server_handling_seconds",
				Help:    "Histogram of response latency (seconds) of gRPC that had been application-level handled by the server.",
				Buckets: prom.DefBuckets,
			}, labels,
		),
	}
}

func (m *ServerMetrics) Describe(ch chan<- *prom.Desc) {
	m.serverHandledCounter.Describe(ch)
	m.serverHandledHistogram.Describe(ch)
}

func (m *ServerMetrics) Collect(ch chan<- prom.Metric) {
	m.serverHandledCounter.Collect(ch)
	m.serverHandledHistogram.Collect(ch)
}

// LabelExtractor must extract the needed labels for each one of the metrics and return
// an array of labels in the SAME ORDER than the ServerMetricLabels used for creating a NewServerMetrics()
type LabelExtractor interface {
	LabelNames() []string
	Labels(context.Context) map[string]string
}

// DefaultLabelExtractor is a dummy LabelExtractor which returns the empty
// list when processing the context to get the CustomLabels
type DefaultLabelExtractor struct{}

// LabelNames returns the names of the extra labels per metric
func (d *DefaultLabelExtractor) LabelNames() []string {
	return []string{}
}

// Labels returns the empty list
func (d *DefaultLabelExtractor) Labels(ctx context.Context) map[string]string {
	return map[string]string{}
}

// Method used for spliting the service/method names of a grpc service
func splitMethodName(fullMethodName string) (string, string) {
	fullMethodName = strings.TrimPrefix(fullMethodName, "/") // remove leading slash
	if i := strings.Index(fullMethodName, "/"); i >= 0 {
		return fullMethodName[:i], fullMethodName[i+1:]
	}
	return "unknown", "unknown"
}

func (m *ServerMetrics) metricLabels(ctx context.Context, labelExtractor LabelExtractor, info *grpc.UnaryServerInfo) map[string]string {
	service, method := splitMethodName(info.FullMethod)

	// Populate basic labels
	labels := map[string]string{
		"grpc_service": service,
		"grpc_method":  method,
	}

	// Populate custom labels
	for k, v := range labelExtractor.Labels(ctx) {
		labels[k] = v
	}

	// Populate non-initialized custom labels with default value
	for _, labelName := range m.labels {
		if _, ok := labels[labelName]; !ok {
			labels[labelName] = "default"
		}
	}
	return labels
}

// UnaryServerInterceptor is a gRPC server-side interceptor that provides Prometheus monitoring for Unary RPCs.
func (m *ServerMetrics) UnaryServerInterceptor(labelExtractor LabelExtractor) func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		metricLabels := m.metricLabels(ctx, labelExtractor, info)
		monitor := newServerReporter(m, metricLabels)
		resp, err := handler(ctx, req)
		st, _ := grpcstatus.FromError(err)
		monitor.labels["grpc_status"] = st.Code().String()
		monitor.Handled()
		return resp, err
	}
}

type serverReporter struct {
	metrics   *ServerMetrics
	labels    map[string]string
	startTime time.Time
}

func newServerReporter(m *ServerMetrics, labels map[string]string) *serverReporter {
	r := &serverReporter{
		metrics:   m,
		labels:    labels,
		startTime: time.Now().UTC(),
	}
	return r
}

func (r *serverReporter) Handled() {
	var orderedLabels []string
	for _, labelName := range r.metrics.labels {
		orderedLabels = append(orderedLabels, r.labels[labelName])
	}

	r.metrics.serverHandledCounter.WithLabelValues(orderedLabels...).Inc()
	r.metrics.serverHandledHistogram.WithLabelValues(orderedLabels...).Observe(time.Since(r.startTime).Seconds())
}
