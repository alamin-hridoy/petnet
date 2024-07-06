package metrics

import (
	"context"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"

	// grpc
	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

// UnaryServerInterceptor returns a server interceptor for reporting request latency.
// default measurement is grpc_request_latency
func (r *Influxdb) UnaryServerInterceptor(measurement string, ignore []string) grpc.UnaryServerInterceptor {
	if r == nil { // no-op on nil
		return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
			return handler(ctx, req)
		}
	}

	if measurement == "" {
		measurement = r.grpc
	}
	if measurement == "" {
		measurement = "grpc_request_latency"
	}
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		mthd := strings.Trim(info.FullMethod, "/")
		for _, v := range ignore {
			if v == mthd {
				return handler(ctx, req)
			}
		}

		ctx = r.ServerSpan(ctx, measurement)
		SetTag(ctx, "grpc_service", path.Base(path.Dir(info.FullMethod)))
		SetTag(ctx, "grpc_method", path.Base(info.FullMethod))
		start := time.Now()
		resp, err := handler(ctx, req)
		elapsed := time.Since(start)
		r.WriteSpan(ctx, map[string]string{
			"error_code": status.Code(err).String(),
		}, map[string]interface{}{
			"latency": elapsed.Milliseconds(),
		})
		return resp, err
	}
}

// StreamServerInterceptor returns a server interceptor for reporting stream latency.
// default measurement is grpc_stream_latency
func (r *Influxdb) StreamServerInterceptor(measurement string) grpc.StreamServerInterceptor {
	if r == nil {
		return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
			return handler(srv, stream)
		}
	}

	if measurement == "" {
		measurement = r.grpc
	}
	if measurement == "" {
		measurement = "grpc_request_latency"
	}
	return func(srv interface{}, stream grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		start := time.Now()
		wrap := grpc_middleware.WrapServerStream(stream)
		ctx := r.ServerSpan(stream.Context(), measurement)
		SetTag(ctx, "grpc_service", path.Base(path.Dir(info.FullMethod)))
		SetTag(ctx, "grpc_method", path.Base(info.FullMethod))
		wrap.WrappedContext = ctx
		err := handler(srv, wrap)
		elapsed := time.Since(start)
		r.WriteSpan(ctx, map[string]string{
			"grpc_code": status.Code(err).String(),
		}, map[string]interface{}{
			"latency": elapsed.Milliseconds(),
		})
		return err
	}
}

// UnaryClientInterceptor returns a client interceptor for reporting request latency.
// default measurement is grpc_client_request_latency
func (r *Influxdb) UnaryClientInterceptor(measurement string, ignore []string) grpc.UnaryClientInterceptor {
	if r == nil { // no-op on nil
		return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
			return invoker(ctx, method, req, reply, cc, opts...)
		}
	}

	if measurement == "" {
		measurement = r.grpc
	}
	if measurement == "" {
		measurement = "grpc_client_request_latency"
	}
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) error {
		for _, v := range ignore {
			if v == method {
				return invoker(ctx, method, req, reply, cc, opts...)
			}
		}

		ctx = r.ServerSpan(ctx, measurement)
		SetTag(ctx, "target", path.Base(cc.Target()))
		SetTag(ctx, "grpc_method", path.Base(method))
		start := time.Now()
		err := invoker(ctx, method, req, reply, cc, opts...)
		elapsed := time.Since(start)
		r.WriteSpan(ctx, map[string]string{
			"error_code": status.Code(err).String(),
		}, map[string]interface{}{
			"latency": elapsed.Milliseconds(),
		})
		return err
	}
}

// StreamClientInterceptor returns a client interceptor for reporting stream latency.
// default measurement is grpc_client_request_latency
func (r *Influxdb) StreamClientInterceptor(measurement string) grpc.StreamClientInterceptor {
	if r == nil {
		return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
			return streamer(ctx, desc, cc, method, opts...)
		}
	}

	if measurement == "" {
		measurement = r.grpc
	}
	if measurement == "" {
		measurement = "grpc_client_request_latency"
	}
	return func(ctx context.Context, desc *grpc.StreamDesc, cc *grpc.ClientConn, method string, streamer grpc.Streamer, opts ...grpc.CallOption) (grpc.ClientStream, error) {
		start := time.Now()
		ctx = r.ServerSpan(ctx, measurement)
		SetTag(ctx, "target", path.Base(cc.Target()))
		SetTag(ctx, "grpc_method", path.Base(method))
		resp, err := streamer(ctx, desc, cc, method, opts...)
		elapsed := time.Since(start)
		r.WriteSpan(ctx, map[string]string{
			"grpc_code": status.Code(err).String(),
		}, map[string]interface{}{
			"latency": elapsed.Milliseconds(),
		})
		return resp, err
	}
}

// ResponseStats wraps responsewriter to record status and response size.
type ResponseStats struct {
	http.ResponseWriter
	statusCode int
	sizeBytes  int
}

func (w *ResponseStats) Status() int {
	return w.statusCode
}

func (w *ResponseStats) Size() int {
	return w.sizeBytes
}

// Write records the amount of the data written.
func (w *ResponseStats) Write(data []byte) (int, error) {
	if w.statusCode == 0 {
		w.statusCode = http.StatusOK
	}
	n, err := w.ResponseWriter.Write(data)
	w.sizeBytes += n
	return n, err
}

// WriteHeader records the status code.
func (w *ResponseStats) WriteHeader(code int) {
	w.ResponseWriter.WriteHeader(code)
	w.statusCode = code
}

type Handler = func(http.ResponseWriter, *http.Request)

// HTTPMiddleware is the full feature middleware.
// Combines HTTPReporter and HTTPSummary to give a drop-in metrics recorder.
func (r *Influxdb) HTTPMiddleware(measurement string) func(http.Handler) http.Handler {
	if r == nil { // noop
		return func(h http.Handler) http.Handler { return h }
	}
	rpt := r.HTTPReporter(measurement)
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(rpt(h.ServeHTTP))
	}
}

// HTTPReporter middleware initializes the context and reports the request fields and tags.
func (r *Influxdb) HTTPReporter(measurement string) func(Handler) Handler {
	if r == nil {
		return func(h Handler) Handler { return h }
	}
	if measurement == "" {
		measurement = r.http
	}
	if measurement == "" {
		measurement = "http_request_latency"
	}
	return func(h Handler) Handler {
		h = HTTPSummary(h)
		return func(w http.ResponseWriter, req *http.Request) {
			ctx := r.ServerSpan(req.Context(), measurement)
			h(w, req.WithContext(ctx))
			r.WriteSpan(ctx, nil, nil)
		}
	}
}

// HTTPSummary middleware adds http summary metrics to the request.
func HTTPSummary(h Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		SetField(ctx, "path", r.URL.Path)
		wr := &ResponseStats{ResponseWriter: w}
		start := time.Now()
		h(wr, r)
		elapsed := time.Since(start)
		SetField(ctx, "latency", elapsed.Milliseconds())
		SetField(ctx, "resp_size", wr.Size())
		SetTag(ctx, "status_code", strconv.Itoa(wr.Status()))
	}
}
