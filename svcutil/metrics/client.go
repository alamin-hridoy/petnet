package metrics

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"

	monitoring "cloud.google.com/go/monitoring/apiv3"
	metricpb "google.golang.org/genproto/googleapis/api/metric"
	monitoredrespb "google.golang.org/genproto/googleapis/api/monitoredres"
	monitoringpb "google.golang.org/genproto/googleapis/monitoring/v3"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/serviceutil/errors"
)

const (
	// the prefix that will be used when sending metrics to stackdriver monitoring,
	// all metrics will be grouped under this name (this is a stackdriver specific
	// namespace)
	metricPrefix = "custom.googleapis.com"
	// store the metrics as a 'global' resource in stackdriver monitoring
	metricResourceType = "global"
	// timeout for creating metrics in stackdriver monitoring
	metricTimeout = time.Second * 5
)

// MetricsClient can be used to send metrics to stackdriver monitoring
type MetricsClient struct {
	gcpProjectID string
	namespace    string

	// labels that will attached to every metric sent to stackdriver
	defaultMetricLabels map[string]string

	client *monitoring.MetricClient
	logger *logrus.Entry

	// startTime for cumulative metrics
	startTime time.Time
}

// NewMetricsClient will return a default metrics client that will
// log metrics to stdout
func NewMetricsClient(gcpProjectID, namespace string) *MetricsClient {
	// Set up the default logger
	logger := logrus.WithFields(logrus.Fields{
		"service": "metrics",
	})

	return &MetricsClient{
		namespace:           namespace,
		defaultMetricLabels: map[string]string{},
		logger:              logger,
		startTime:           time.Now().UTC(),
	}
}

// connectToService connects to Stackdriver monitoring
func (m *MetricsClient) connectToService() error {
	ctx := context.Background()
	client, err := monitoring.NewMetricClient(ctx)
	if err != nil {
		return errors.Wrap(err, "unable to create stack driver client")
	}
	m.client = client
	return nil
}

// projectResource returns a formatted gcp project identifier
func (m MetricsClient) projectResource() string {
	return "projects/" + m.gcpProjectID
}

// formatMetric returns the metric name prefixed with the metrics namespace
func (m MetricsClient) formatMetric(metricName string) string {
	return fmt.Sprintf("%s/brankas-%s/%s", metricPrefix, m.namespace, metricName)
}

// createSingleMetric currently only supports GAUGE and CUMULATIVE metrics. This will
// create a single data point and send it to stackdriver monitoring
func (m MetricsClient) createSingleMetric(metricType, metricName string, metricValue *monitoringpb.TypedValue, labels map[string]string) error {
	/*
		Interval: The time interval to which the data point applies.
		- GAUGE metrics, only the end time of the interval is used.
		- DELTA metrics, the start and end time should specify a
		non-zero interval, with subsequent points specifying contiguous
		and non-overlapping intervals (currently not supported).
		- CUMULATIVE metrics, the start and end time should
		specify a non-zero interval, with subsequent points specifying the
		same start time and increasing end times, until an event resets the
		cummuative value to zero and sets a new start time for the following
		points.
	*/

	var metricKind metricpb.MetricDescriptor_MetricKind
	metricInterval := &monitoringpb.TimeInterval{
		EndTime: &tspb.Timestamp{Seconds: time.Now().UTC().Unix()},
	}

	switch metricType {
	case "CUMULATIVE":
		metricKind = metricpb.MetricDescriptor_CUMULATIVE
		metricInterval.StartTime = &tspb.Timestamp{Seconds: m.startTime.Unix()}
	case "GAUGE":
		metricKind = metricpb.MetricDescriptor_GAUGE
	default:
		return fmt.Errorf("unknown metric type value received: %s", metricType)
	}

	// If we haven't configured a client yet, don't try and send the
	// metric to stackdriver. However, we do this after we check the
	// metricType incase that is invalid.
	if m.client == nil {
		return nil
	}

	// Add and overwrite the default metric labels with the passed
	// through labels
	metricLabels := m.defaultMetricLabels
	for k, v := range labels {
		metricLabels[k] = v
	}

	timeseries := monitoringpb.TimeSeries{
		MetricKind: metricKind,
		Metric: &metricpb.Metric{
			Type:   m.formatMetric(metricName),
			Labels: metricLabels,
		},
		Resource: &monitoredrespb.MonitoredResource{
			Labels: map[string]string{
				"project_id": m.gcpProjectID,
			},
			Type: metricResourceType,
		},
		Points: []*monitoringpb.Point{
			{
				Interval: metricInterval,
				Value:    metricValue,
			},
		},
	}

	createTimeSeriesRequest := &monitoringpb.CreateTimeSeriesRequest{
		Name:       m.projectResource(),
		TimeSeries: []*monitoringpb.TimeSeries{&timeseries},
	}

	ctx, cancel := context.WithTimeout(context.Background(), metricTimeout)
	defer cancel()
	if err := m.client.CreateTimeSeries(ctx, createTimeSeriesRequest); err != nil {
		return err
	}

	return nil
}

// Count is a cumulative metric that represents a single numerical value that
// only ever goes up. A counter is typically used to count requests served,
// tasks completed, errors occurred, etc. Counters should not be used to
// expose current counts of items whose number can also go down,
// e.g. the number of currently running goroutines. Use gauges for this use case
func (m MetricsClient) Count(metricName string, metricValue *monitoringpb.TypedValue, labels map[string]string) {
	if err := m.createSingleMetric("CUMULATIVE", metricName, metricValue, labels); err != nil {
		m.logger.WithError(err).WithFields(logrus.Fields{
			"metric_type": "CUMULATIVE",
			"metric_name": metricName,
		}).Errorf("error creating single count metric")
	}
}

// CountInt64 provides a wrapper around `MetricsClient.Count` with the 'metricValue' type as int64
func (m MetricsClient) CountInt64(metricName string, metricValue int64, labels map[string]string) {
	go m.Count(metricName, monitoringTypedValueInt64(metricValue), labels)
}

// CountFloat64 provides a wrapper around `MetricsClient.Count` with the 'metricValue' type as float64
func (m MetricsClient) CountFloat64(metricName string, metricValue float64, labels map[string]string) {
	go m.Count(metricName, monitoringTypedValueFloat64(metricValue), labels)
}

// Gauge is a metric that represents a single numerical value that can arbitrarily
// go up and down.Gauges are typically used for measured values like temperatures
// or current memory usage, but also "counts" that can go up and down, like the
// number of running goroutines.
func (m MetricsClient) Gauge(metricName string, metricValue *monitoringpb.TypedValue, labels map[string]string) {
	if err := m.createSingleMetric("GAUGE", metricName, metricValue, labels); err != nil {
		m.logger.WithError(err).WithFields(logrus.Fields{
			"metric_type": "GAUGE",
			"metric_name": metricName,
		}).Errorf("error creating single gauge metric")
	}
}

// GaugeInt64 provides a wrapper around `MetricsClient.Gauge` with the 'metricValue' type as int64
func (m MetricsClient) GaugeInt64(metricName string, metricValue int64, labels map[string]string) {
	go m.Gauge(metricName, monitoringTypedValueInt64(metricValue), labels)
}

// GaugeFloat64 provides a wrapper around `MetricsClient.Gauge` with the 'metricValue' type as float64
func (m MetricsClient) GaugeFloat64(metricName string, metricValue float64, labels map[string]string) {
	go m.Gauge(metricName, monitoringTypedValueFloat64(metricValue), labels)
}

// GaugedBool provides a wrapper around `MetricsClient.Gauge` with the 'metricValue' type as bool
func (m MetricsClient) GaugeBool(metricName string, metricValue bool, labels map[string]string) {
	go m.Gauge(metricName, monitoringTypedValueBool(metricValue), labels)
}

// GaugeString provides a wrapper around `MetricsClient.Gauge` with the 'metricValue' type as string
func (m MetricsClient) GaugeString(metricName, metricValue string, labels map[string]string) {
	go m.Gauge(metricName, monitoringTypedValueString(metricValue), labels)
}

// SetLogger sets the logger
func (m *MetricsClient) SetLogger(logger *logrus.Entry) {
	m.logger = logger.WithFields(logrus.Fields{
		"service":           "metrics",
		"metrics_namespace": m.namespace,
	})
}

// SetGCPProjectID will set the metrics client to connect to stackdriver
// monitoring and send metrics there
func (m *MetricsClient) Configure(projectID, namespace string) error {
	m.gcpProjectID = projectID
	m.namespace = namespace
	return m.connectToService()
}

// AddDefaultLabels will append to the default labels
func (m *MetricsClient) AddDefaultLabels(labels map[string]string) {
	for k, v := range labels {
		m.defaultMetricLabels[k] = v
	}
}

func monitoringTypedValueInt64(value int64) *monitoringpb.TypedValue {
	return &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_Int64Value{Int64Value: value}}
}

func monitoringTypedValueFloat64(value float64) *monitoringpb.TypedValue {
	return &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_DoubleValue{DoubleValue: value}}
}

func monitoringTypedValueBool(value bool) *monitoringpb.TypedValue {
	return &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_BoolValue{BoolValue: value}}
}

func monitoringTypedValueString(value string) *monitoringpb.TypedValue {
	return &monitoringpb.TypedValue{Value: &monitoringpb.TypedValue_StringValue{StringValue: value}}
}
