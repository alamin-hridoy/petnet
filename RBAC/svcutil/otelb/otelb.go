// TODO: probably good to move this into brank.as/rbac/serviceutil so it can be shared
//
// To instrument add this to the beginning of a function
// log, otl, ctx := otelb.Start(ctx)
// defer otl.Span.End()
package otelb

import (
	"context"
	"log"

	"go.opentelemetry.io/otel"
	otlpgrpc "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"brank.as/rbac/serviceutil/logging"
)

// Initializes an OTLP exporter, and configures the corresponding trace providers.
// usually used in main function
func InitOTELProvider(ctx context.Context, svcName string, u string) func() {
	if u == "" {
		logging.FromContext(ctx).Debug("tracing disabled")
		return func() {}
	}
	logging.FromContext(ctx).Debug("dialing trace collector")
	exp, err := otlpgrpc.New(ctx,
		otlpgrpc.WithInsecure(),
		otlpgrpc.WithEndpoint(u),
		otlpgrpc.WithDialOption(),
	)
	if err != nil {
		log.Fatalln("dialing OTLP agent")
	}

	// TODO: add more resource attributes
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String(svcName),
		),
	)
	if err != nil {
		log.Fatalln("adding resource attributes")
	}

	bsp := sdktrace.NewBatchSpanProcessor(exp)
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sdktrace.AlwaysSample()),
		sdktrace.WithResource(res),
		sdktrace.WithSpanProcessor(bsp),
	)

	// set global propagator to tracecontext (the default is no-op).
	otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTracerProvider(tp)

	return func() {
		tp.Shutdown(ctx)
		if err != nil {
			log.Fatalln("stopping provider")
		}
		exp.Shutdown(ctx)
		if err != nil {
			log.Fatalln("stopping exporter")
		}
	}
}
