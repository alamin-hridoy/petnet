package metrics

import "github.com/sirupsen/logrus"

var defaultClient *MetricsClient

func init() {
	defaultClient = NewMetricsClient("local", "local")
}

// SetGCPProjectID will set the metrics client to connect to stackdriver
// monitoring and send metrics there
func Configure(projectID, namespace string) error {
	return defaultClient.Configure(projectID, namespace)
}

// SetLogger sets the logger
func SetLogger(logger *logrus.Entry) {
	defaultClient.SetLogger(logger)
}

// AddDefaultLabels will append to the default labels
func AddDefaultLabels(labels map[string]string) {
	defaultClient.AddDefaultLabels(labels)
}

// CountInt64 provides a wrapper around `MetricsClient.Count` with the 'metricValue' type as int64
func CountInt64(metricName string, metricValue int64, labels map[string]string) {
	defaultClient.CountInt64(metricName, metricValue, labels)
}

// CountFloat64 provides a wrapper around `MetricsClient.Count` with the 'metricValue' type as float64
func CountFloat64(metricName string, metricValue float64, labels map[string]string) {
	defaultClient.CountFloat64(metricName, metricValue, labels)
}

// GaugeInt64 provides a wrapper around `MetricsClient.Gauge` with the 'metricValue' type as int64
func GaugeInt64(metricName string, metricValue int64, labels map[string]string) {
	defaultClient.GaugeInt64(metricName, metricValue, labels)
}

// GaugeFloat64 provides a wrapper around `MetricsClient.Gauge` with the 'metricValue' type as float64
func GaugeFloat64(metricName string, metricValue float64, labels map[string]string) {
	defaultClient.GaugeFloat64(metricName, metricValue, labels)
}

// GaugedBool provides a wrapper around `MetricsClient.Gauge` with the 'metricValue' type as bool
func GaugeBool(metricName string, metricValue bool, labels map[string]string) {
	defaultClient.GaugeBool(metricName, metricValue, labels)
}

// GaugeString provides a wrapper around `MetricsClient.Gauge` with the 'metricValue' type as string
func GaugeString(metricName string, metricValue string, labels map[string]string) {
	defaultClient.GaugeString(metricName, metricValue, labels)
}
