package otelb

import (
	"errors"

	"github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

type LogHook struct{}

func (h *LogHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

// Fire hook for logging tracing events
func (h *LogHook) Fire(e *logrus.Entry) error {
	s := trace.SpanFromContext(e.Context)
	// todo: fix so trace-id and span-id are logged, something is wrong with this
	// and they just zeroes
	// tid := s.SpanContext().TraceID()
	// sid := s.SpanContext().SpanID()
	// e.Data["trace-id"] = tid
	// e.Data["span-id"] = sid

	as := []attribute.KeyValue{
		semconv.CodeFilepathKey.String(e.Caller.File),
		semconv.CodeLineNumberKey.Int(e.Caller.Line),
	}
	switch e.Level {
	case logrus.ErrorLevel:
		s.RecordError(errors.New(e.Message), trace.WithAttributes(as...))
		s.SetStatus(codes.Error, "internal error")
	default:
		as := append(as, semconv.MessageIDKey.String(e.Message))
		s.AddEvent(e.Level.String(), trace.WithAttributes(as...))
	}
	return nil
}
