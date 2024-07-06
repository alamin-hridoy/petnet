package metrics

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strconv"
	"time"

	"github.com/spf13/viper"

	// influxdb
	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	"github.com/influxdata/influxdb-client-go/v2/api"
	"github.com/influxdata/influxdb-client-go/v2/api/write"
)

type Influxdb struct {
	w              api.WriteAPI
	close          func()
	reqMeasurement string
	defTag         map[string]string
	grpc           string
	http           string
}

type Config struct {
	Host         string `mapstructure:"host"`
	Port         int    `mapstructure:"port"`
	BatchSize    int    `mapstructure:"batchSize"`
	Bucket       string `mapstructure:"bucket"`
	Organization string `mapstructure:"organization"`
	Token        string `mapstructure:"token"`
	GRPCMeasure  string `mapstructure:"grpcmeasure"` // note: M is not capitalized
	HTTPMeasure  string `mapstructure:"httpmeasure"` // note: M is not capitalized
}

func NewInfluxDBClientFromConfig(config Config) (*Influxdb, error) {
	host := config.Host
	port := config.Port
	batchSize := config.BatchSize
	bucket := config.Bucket
	organization := config.Organization
	authToken := config.Token

	grpcMeasure := config.GRPCMeasure
	httpMeasure := config.HTTPMeasure

	switch "" {
	case host, organization, authToken:
		fmt.Println("skipping metrics initialization (missing config entries)")
		return nil, nil
	}

	influxURL, err := url.Parse("http://" + host + ":" + strconv.Itoa(port))
	if err != nil {
		return nil, err
	}

	idbcl := influxdb2.NewClientWithOptions(influxURL.String(), authToken,
		influxdb2.DefaultOptions().SetBatchSize(uint(batchSize)))

	influxCL, err := NewInfluxdb(idbcl, organization, bucket)
	if err != nil {
		return nil, err
	}
	influxCL.grpc = grpcMeasure
	influxCL.http = httpMeasure

	return influxCL, nil
}

// Configures an influxdb metrics clients using the default config entries.
// No-op when configuration entries are not available,
// to allow easy disabling on local and/or new deployments.
func NewInfluxDBClient(config *viper.Viper) (*Influxdb, error) {
	var allConfig struct {
		InfluxDB Config `mapstructure:"influxdb"`
	}

	if err := config.Unmarshal(&allConfig); err != nil {
		return nil, fmt.Errorf("cannot unmarshal influx config: %w", err)
	}

	return NewInfluxDBClientFromConfig(allConfig.InfluxDB)
}

// NewInfluxdb creates a client for reporting metrics to influxdb
func NewInfluxdb(cl influxdb2.Client, organization, bucket string) (*Influxdb, error) {
	if cl == nil {
		return nil, fmt.Errorf("missing influxdb client")
	}
	return &Influxdb{
		w:              cl.WriteAPI(organization, bucket),
		close:          cl.Close,
		reqMeasurement: "http_client_latency",
		defTag:         map[string]string{},
	}, nil
}

// ErrorsFunc will operate on all write errors until the client is closed.
func (r *Influxdb) ErrorsFunc(f func(error)) {
	if r == nil {
		return
	}
	for err := range r.w.Errors() {
		f(err)
	}
}

// Close the reporter.
func (r *Influxdb) Close() {
	if r == nil {
		return
	}
	r.w.Flush()
	r.close()
}

// DefaultTags adds default tags to all datapoints reported.
func (r *Influxdb) DefaultTags(t map[string]string) {
	if r == nil {
		return
	}
	r.defTag = t
}

// MetricsClient traces latencies for all requests sent via the client.
func (r *Influxdb) MetricsClient(cl *http.Client, measurement string) *http.Client {
	if r == nil {
		return cl
	}
	if measurement == "" {
		measurement = r.reqMeasurement
	}
	if measurement == "" {
		measurement = r.http
	}

	wr := new(http.Client)
	if cl != nil {
		*wr = *cl // copy existing configuration from the client
	}
	if wr.Transport == nil {
		wr.Transport = http.DefaultTransport
	}
	wr.Transport = &tripper{
		measurement: measurement,
		r:           r,
		rt:          wr.Transport,
	}
	return wr
}

// SetRequestMeasurement assigns a default measurement for all requests.
func (r *Influxdb) SetRequestMeasurement(measurement string) {
	if r == nil {
		return
	}
	if measurement == "" {
		measurement = "http_client_latency"
	}
	r.reqMeasurement = measurement
}

type traceKey struct{}

type clientTrace struct {
	msmt      string
	start     time.Time
	conn      time.Duration
	writeHdr  time.Duration
	writeReq  time.Duration
	firstByte time.Duration
}

// since start of trace.
func (tr *clientTrace) since() time.Duration {
	if tr.start.IsZero() {
		tr.start = time.Now()
		return 0
	}
	return time.Since(tr.start)
}

// WithRequestMetricsMeasurement loads latency metrics collector in the context for use with http.Client
func (r *Influxdb) WithRequestMetricsMeasurement(ctx context.Context, measurement string) context.Context {
	if r == nil {
		return ctx
	}
	if ctx.Value(traceKey{}) != nil {
		return ctx
	}
	if measurement == "" {
		measurement = r.reqMeasurement
	}
	if measurement == "" {
		measurement = r.http
	}
	t := &clientTrace{
		msmt: measurement,
	}
	trace := &httptrace.ClientTrace{
		DNSStart:             func(httptrace.DNSStartInfo) { t.since() },
		GetConn:              func(string) { t.since() },
		GotConn:              func(httptrace.GotConnInfo) { t.conn = t.since() },
		WroteHeaders:         func() { t.writeHdr = t.since() },
		WroteRequest:         func(httptrace.WroteRequestInfo) { t.writeReq = t.since() },
		GotFirstResponseByte: func() { t.firstByte = t.since() },
	}
	ctx = WithTags(ctx, tagsSpan(ctx))
	ctx = WithFields(ctx, fieldsSpan(ctx))
	return httptrace.WithClientTrace(context.WithValue(ctx, traceKey{}, t), trace)
}

// WithRequestMetrics loads latency metrics collector in the context for use with http.Client
func (r *Influxdb) WithRequestMetrics(ctx context.Context) context.Context {
	if r == nil {
		return ctx
	}
	return r.WithRequestMetricsMeasurement(ctx, r.reqMeasurement)
}

// IsRequestRecorder tests if the context contains a recorder for http client requests.
func IsRequestRecorder(ctx context.Context) bool {
	t, ok := ctx.Value(traceKey{}).(*clientTrace)
	return ok && t != nil
}

// NewRequest creates an instrumented HTTP request for use with http.Client.Do call.
func (r *Influxdb) NewRequest(ctx context.Context, mthd string, u *url.URL, body io.Reader) (*http.Request, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	// initialize recorder
	ctx = r.WithRequestMetrics(ctx)
	// add standard tags
	SetTags(ctx, map[string]string{
		"host":   u.Host,
		"method": mthd,
	})
	return http.NewRequestWithContext(ctx, mthd, u.String(), body)
}

// FinalizeRequest finalizes and reports the metrics from the request.
// Merges all tags in current span created by `WithTags` and fields created by `WithFields`.
func (r *Influxdb) FinalizeRequest(ctx context.Context, tags map[string]string, fields map[string]interface{}) {
	if r == nil {
		return
	}
	tm := time.Now()
	t, ok := ctx.Value(traceKey{}).(*clientTrace)
	if !ok {
		return
	}
	tg := map[string]string{}
	// default first, to allow overrides
	for k, v := range r.defTag {
		tg[k] = v
	}
	for k, v := range tagsSpan(ctx) {
		tg[k] = v
	}
	for k, v := range tags {
		tg[k] = v
	}
	// set default measurements, allow overwrites.
	fld := map[string]interface{}{
		"conn":       t.conn.Milliseconds(),
		"header":     t.writeHdr.Milliseconds(),
		"write_req":  t.writeReq.Milliseconds(),
		"first_byte": t.firstByte.Milliseconds(),
		"total":      tm.Sub(t.start).Milliseconds(),
	}
	for k, v := range fieldsSpan(ctx) {
		fld[k] = v
	}
	for k, v := range fields {
		fld[k] = v
	}
	pt := write.NewPoint(t.msmt, tg, fld, tm)
	r.w.WritePoint(pt.SortFields().SortTags())
}

type tripper struct {
	measurement string
	r           *Influxdb
	rt          http.RoundTripper
}

// NewTransport to record requests metrics by intercepting the http round trip.
func (r *Influxdb) NewTransport(measurement string, rt http.RoundTripper) http.RoundTripper {
	if measurement == "" {
		measurement = r.reqMeasurement
	}
	if measurement == "" {
		measurement = r.http
	}
	if rt == nil || measurement == "" {
		rt = http.DefaultTransport
	}
	return &tripper{
		measurement: measurement,
		r:           r,
		rt:          rt,
	}
}

type recorderError struct {
	err error
	r   *Influxdb
	ctx context.Context
}

func (re *recorderError) Unwrap() error {
	return re.err
}

func (re *recorderError) Error() string {
	return re.err.Error()
}

func IsRecorderError(err error) bool {
	r := &recorderError{}
	return errors.Is(err, r)
}

func ErrorContext(err error) context.Context {
	r := &recorderError{}
	if errors.As(err, &r) {
		return r.ctx
	}
	return nil
}

func CloseError(err error) { FinalizeError(err, nil, nil) }

func FinalizeError(err error, tag map[string]string, fld map[string]interface{}) {
	e := &recorderError{}
	if errors.As(err, &e) {
		if fld == nil {
			fld = map[string]interface{}{}
		}
		if fld["error"] == nil {
			fld["error"] = err.Error()
		}
		e.r.FinalizeRequest(e.ctx, tag, fld)
	}
}

// RoundTrip fulfils the http.RoundTripper interface
func (tr *tripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if tr.r == nil {
		return tr.rt.RoundTrip(req)
	}
	if !IsRequestRecorder(req.Context()) {
		ctx := tr.r.WithRequestMetricsMeasurement(req.Context(), tr.measurement)
		RequestTags(ctx, req)
		req = req.WithContext(ctx)
	}
	req.Context().Value(traceKey{}).(*clientTrace).start = time.Now()
	resp, err := tr.rt.RoundTrip(req)
	if err != nil {
		return nil, &recorderError{err: err, ctx: req.Context(), r: tr.r}
	}
	ctx := resp.Request.Context()
	SetTag(ctx, "http_code", strconv.Itoa(resp.StatusCode))
	r := &Response{r: resp.Body, ctx: ctx, tr: tr}
	resp.Body = r
	resp.Request = resp.Request.WithContext(context.WithValue(ctx, responseKey{}, r))
	return resp, nil
}

var _ io.ReadCloser = (*Response)(nil)

type responseKey struct{}

type Response struct {
	len int
	r   io.ReadCloser
	ctx context.Context
	tr  *tripper
}

func (c *Response) Read(b []byte) (int, error) {
	n, err := c.r.Read(b)
	c.len += n
	return n, err
}

func (c *Response) Close() error {
	c.tr.r.FinalizeRequest(c.ctx, nil, nil)
	if c.r == nil {
		return nil
	}
	return c.r.Close()
}

func (c *Response) Tag(t map[string]string) {
	if c != nil {
		SetTags(c.ctx, t)
	}
}

func (c *Response) Field(f map[string]interface{}) {
	if c != nil {
		SetFields(c.ctx, f)
	}
}

type tagsKey struct{}

func ResponseTags(resp *http.Response, tags map[string]string) {
	responseRecorder(resp).Tag(tags)
}

func ResponseTag(resp *http.Response, key, value string) {
	responseRecorder(resp).Tag(map[string]string{key: value})
}

func ResponseFields(resp *http.Response, fld map[string]interface{}) {
	responseRecorder(resp).Field(fld)
}

func ResponseField(resp *http.Response, fld string, value interface{}) {
	responseRecorder(resp).Field(map[string]interface{}{fld: value})
}

func IsResponseRecorder(resp *http.Response) bool {
	return responseRecorder(resp) != nil
}

func responseRecorder(resp *http.Response) *Response {
	if resp == nil || resp.Request == nil {
		return nil
	}
	r, _ := resp.Request.Context().Value(responseKey{}).(*Response)
	return r
}

// WithTags creates a new span for metrics tags.  Does not inherit tags from the current span.
// To inherit a snapshot of tags from the current span, combine: `WithTags(ctx, Tags(ctx))`
// Any `SetTag` requests in the new span will be independent of the current/previous span.
func WithTags(ctx context.Context, tags map[string]string) context.Context {
	t := make(map[string]string)
	for k, v := range tags {
		t[k] = v
	}
	return withTags(ctx, t)
}

func withTags(ctx context.Context, tags map[string]string) context.Context {
	return context.WithValue(ctx, tagsKey{}, tags)
}

func SetTag(ctx context.Context, tag, value string) {
	t := Tags(ctx)
	if t == nil {
		return
	}
	t[tag] = value
}

func SetTags(ctx context.Context, tags map[string]string) {
	t := Tags(ctx)
	if t == nil {
		return
	}
	for k, v := range tags {
		t[k] = v
	}
}

func Tags(ctx context.Context) map[string]string {
	t := tagsSpan(ctx)
	if t != nil {
		return t
	}
	span := getSpan(ctx)
	if span != nil {
		return span.tags
	}
	return nil
}

func tagsSpan(ctx context.Context) map[string]string {
	t, ok := ctx.Value(tagsKey{}).(map[string]string)
	if ok && t != nil {
		return t
	}
	return nil
}

func RequestTags(ctx context.Context, req *http.Request) {
	t := Tags(ctx)
	if t == nil {
		return
	}
	t["host"] = req.URL.Host
	t["method"] = req.Method
}

type fieldsKey struct{}

// WithFields creates a new span for metrics fields.  Does not inherit fields from the current span.
// To inherit a snapshot of fields from the current span, combine: `WithFields(ctx, Fields(ctx))`
// Any `SetField` requests in the new span will be independent of the current/previous span.
func WithFields(ctx context.Context, fields map[string]interface{}) context.Context {
	fld := make(map[string]interface{})
	for k, v := range fields {
		fld[k] = v
	}
	return withFields(ctx, fld)
}

func withFields(ctx context.Context, fields map[string]interface{}) context.Context {
	if fields == nil {
		return context.WithValue(ctx, fieldsKey{}, nil)
	}
	return context.WithValue(ctx, fieldsKey{}, fields)
}

func SetField(ctx context.Context, key string, value interface{}) {
	t := Fields(ctx)
	if t == nil {
		return
	}
	t[key] = value
}

func SetFields(ctx context.Context, fields map[string]interface{}) {
	t := Fields(ctx)
	if t == nil {
		return
	}
	for k, v := range fields {
		t[k] = v
	}
}

func Fields(ctx context.Context) map[string]interface{} {
	fld := fieldsSpan(ctx)
	if fld != nil {
		return fld
	}
	span := getSpan(ctx)
	if span == nil {
		return nil
	}
	return span.fields
}

func fieldsSpan(ctx context.Context) map[string]interface{} {
	fld, ok := ctx.Value(fieldsKey{}).(map[string]interface{})
	if ok && fld != nil {
		return fld
	}
	return nil
}

type (
	spanKey struct{}
	span    struct {
		measurement string
		tags        map[string]string
		fields      map[string]interface{}
	}
)

// ServerSpan creates a metrics recorder for a server handling a request.
// Span will contain any previous tags created by `WithTags` and fields created by `WithFields`.
func (r *Influxdb) ServerSpan(ctx context.Context, measurement string) context.Context {
	if r == nil {
		return ctx
	}
	if measurement == "" {
		measurement = r.reqMeasurement
	}
	f, t := fieldsSpan(ctx), tagsSpan(ctx)
	if f == nil {
		// no fields or fields are from span.
		f = map[string]interface{}{}
	} else {
		// mask existing entry to avoid accidental duplication
		ctx = withFields(ctx, nil)
	}
	if t == nil {
		// no tags or tags are from span.
		t = map[string]string{}
	} else {
		// mask existing entry to avoid accidental duplication
		ctx = withTags(ctx, nil)
	}
	return context.WithValue(ctx, spanKey{}, span{measurement: measurement, tags: t, fields: f})
}

// WriteSpan publishes the recorded metrics to influxdb.
func (r *Influxdb) WriteSpan(ctx context.Context, tags map[string]string, fields map[string]interface{}) {
	if r == nil {
		return
	}
	span := getSpan(ctx)
	if span == nil {
		return
	}
	tg := map[string]string{}
	// default first, to allow overrides
	for k, v := range r.defTag {
		tg[k] = v
	}
	for k, v := range span.tags {
		tg[k] = v
	}
	// Combine all tag sources
	for k, v := range tagsSpan(ctx) {
		tg[k] = v
	}
	for k, v := range tags {
		tg[k] = v
	}
	// Combine all field sources
	for k, v := range fieldsSpan(ctx) {
		span.fields[k] = v
	}
	for k, v := range fields {
		span.fields[k] = v
	}
	r.w.WritePoint(write.NewPoint(span.measurement, tg, span.fields, time.Now()).
		SortFields().SortTags())
}

// getSpan extracts the span from the context.
func getSpan(ctx context.Context) *span {
	s, ok := ctx.Value(spanKey{}).(span)
	if !ok {
		return nil
	}
	return &s
}
