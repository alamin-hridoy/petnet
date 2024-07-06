package otelb

import (
	"context"
	"runtime"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"

	"brank.as/rbac/serviceutil/logging"
)

type Option func(*config)

type config struct {
	spanName   string
	tracerName string
}

// WithSpanName option sets the span name, if empty defaults to function name
// of the the caller to Start is used
func WithSpanName(sn string) Option {
	return func(c *config) { c.spanName = sn }
}

// WithTracerName option sets the tracer name, if empty a default value is set by
// opentelemetry
func WithTracerName(sn string) Option {
	return func(c *config) { c.tracerName = sn }
}

type Otel struct {
	Span trace.Span
}

// Start a span which is used to trace a function
func Start(ctx context.Context, opts ...Option) (*logrus.Entry, *Otel, context.Context) {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}

	// todo: leaving tracename as empty by default for now, it will be set to its default
	// we might want to set this to package name or something, not really sure
	// what the usecase is for this yet
	var t trace.Tracer
	if c.tracerName != "" {
		t = otel.Tracer(c.tracerName)
	} else {
		t = otel.Tracer("")
	}

	var s trace.Span
	if c.spanName != "" {
		ctx, s = t.Start(ctx, c.spanName)
	} else {
		pc, _, _, ok := runtime.Caller(1)
		details := runtime.FuncForPC(pc)
		// todo: might be a better way of setting the span name to the function name
		// maybe we just want the function name without the full path
		if ok && details != nil {
			ctx, s = t.Start(ctx, details.Name())
		} else {
			ctx, s = t.Start(ctx, "error")
		}
	}
	log := logging.FromContext(ctx).WithContext(ctx)
	log.Logger.AddHook(&LogHook{})
	otl := &Otel{Span: s}
	return log, otl, ctx
}
