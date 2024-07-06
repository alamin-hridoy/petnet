package metrics

import (
	"context"
	"log"
	"path"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/stats"
)

const (
	// Metric names for grpc metrics
	metricConnectionOpened = "grpc_connection_open"
	metricConnectionClosed = "grpc_connection_closed"
	metricRPCBegin         = "grpc_rpc_begin"
	metricRPCEnd           = "grpc_rpc_end"
	metricRPCRecvMessage   = "grpc_rpc_recv_message"
	metricRPCSentMessage   = "grpc_rpc_sent_message"
)

// Ensure that the MetricsGRPCHandler implements stats.Handler
var _ stats.Handler = (*MetricsGRPCHandler)(nil)

// MetricsGRPCHandler can be used as a grpc stats.Handler interface
// to record metrics
type MetricsGRPCHandler struct{}

// NewMetricsGRPCHandler returns a handler that satifies grpc stats.Handler
// and will record grpc metrics regarding connections and rpc requests
func NewMetricsGRPCHandler() *MetricsGRPCHandler {
	return &MetricsGRPCHandler{}
}

// HandleConn exists to satisfy gRPC stats.Handler and will send metrics
// for each open and closed connection
func (s *MetricsGRPCHandler) HandleConn(ctx context.Context, cs stats.ConnStats) {
	switch cs := cs.(type) {
	case *stats.ConnBegin:
		defaultClient.CountInt64(metricConnectionOpened, 1, nil)
	case *stats.ConnEnd:
		defaultClient.CountInt64(metricConnectionClosed, 1, nil)
	default:
		log.Printf("unknown stats.ConnStats type: %s\n", cs)
	}
}

// TagConn exists to satisfy gRPC stats.Handler
func (s *MetricsGRPCHandler) TagConn(ctx context.Context, cti *stats.ConnTagInfo) context.Context {
	// no-op
	return ctx
}

// HandleRPC implements per-RPC metrics
func (s *MetricsGRPCHandler) HandleRPC(ctx context.Context, rs stats.RPCStats) {
	switch rs := rs.(type) {
	case *stats.Begin:
		defaultClient.CountInt64(metricRPCBegin, 1, nil)
	case *stats.InPayload:
		defaultClient.CountInt64(metricRPCRecvMessage, 1, nil)
	case *stats.OutPayload:
		defaultClient.CountInt64(metricRPCSentMessage, 1, nil)
	case *stats.End:
		labels := map[string]string{
			"success": "true",
		}
		if rs.Error != nil {
			labels["success"] = "false"
		}
		defaultClient.CountInt64(metricRPCEnd, 1, labels)
	}
}

// TagRPC implements per-RPC context management.
func (s *MetricsGRPCHandler) TagRPC(ctx context.Context, rti *stats.RPCTagInfo) context.Context {
	// no-op
	return ctx
}

// UnaryServerInterceptor will create metrics for each grpc request about how long
// it takes to handle a request. This will create a metrics based on the grpc method.
func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		startTime := time.Now()
		resp, err := handler(ctx, req)

		// How long did the request take to the nearest millisecond.
		timeSince := time.Since(startTime) / time.Millisecond

		// eg: brankas.v2.account.AccountService
		service := path.Dir(info.FullMethod)[1:]
		// eg: RetrieveAccount
		method := path.Base(info.FullMethod)

		// Labels to add to the metric.
		labels := map[string]string{
			"success": "true",
			"service": service,
		}
		if err != nil {
			labels["success"] = "false"
		}

		// Gauge how long it took to handle the request.
		defaultClient.GaugeInt64("grpc_"+strings.ToLower(method), int64(timeSince), labels)

		return resp, err
	}
}
