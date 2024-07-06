package mw

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/csrf"
	"github.com/sirupsen/logrus"
	"github.com/urfave/negroni"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"

	"brank.as/petnet/serviceutil/logging"
	"brank.as/rbac/svcutil/otelb"
)

type User interface {
	Use(func(http.Handler) http.Handler)
}

func ChainHTTPMiddleware(usr User, log logrus.FieldLogger, mw ...func(http.Handler) http.Handler) {
	for _, f := range []func(http.Handler) http.Handler{
		Logger(log),
		Gzip,
		ContentType("text/html"),
	} {
		usr.Use(f)
	}
	for _, f := range mw {
		usr.Use(f)
	}
}

func Logger(log logrus.FieldLogger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			log, otl, ctx := otelb.Start(req.Context(), otelb.WithSpanName(req.URL.String()))
			defer otl.Span.End()

			lrw := negroni.NewResponseWriter(w)
			ctx = logging.WithLogger(ctx, log)
			h.ServeHTTP(lrw, req.WithContext(ctx))

			// todo: add more attributes
			otl.Span.SetAttributes(semconv.HTTPStatusCodeKey.String(strconv.Itoa(lrw.Status())))
			otl.Span.SetAttributes(semconv.HTTPMethodKey.String(req.Method))
			otl.Span.SetAttributes(semconv.HTTPTargetKey.String(req.URL.String()))
			otl.Span.SetAttributes(semconv.HTTPHostKey.String(strings.ToLower(req.Host)))
			otl.Span.SetAttributes(semconv.HTTPClientIPKey.String(GetIP(req)))
			otl.Span.SetAttributes(semconv.HTTPSchemeKey.String(req.URL.Scheme))
		})
	}
}

func GetIP(r *http.Request) string {
	forwarded := r.Header.Get("X-forwarded-for")
	if forwarded != "" {
		return forwarded
	}
	return r.RemoteAddr
}

func CSRF(secret []byte, opts ...csrf.Option) func(h http.Handler) http.Handler {
	opts = append([]csrf.Option{
		csrf.ErrorHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := logging.FromContext(r.Context())
			log.WithFields(logrus.Fields{
				"csrf_error": csrf.FailureReason(r).Error(),
				"token":      csrf.Token(r),
				"template":   csrf.TemplateField(r),
			}).Error("csrf error")
			fmt.Fprintln(w, csrf.FailureReason(r))
		})),
	}, opts...)
	return csrf.Protect(secret, opts...)
}

type contentWriter struct {
	http.ResponseWriter
	def string
	log *logrus.Entry
}

func (c contentWriter) Write(b []byte) (int, error) {
	if c.Header().Get("Content-Type") == "" {
		ct := http.DetectContentType(b)
		if ct == "" {
			ct = c.def
		}
		c.log.WithField("content-type", ct).Trace("set")
		c.Header().Set("Content-Type", ct)
	}
	c.log.WithField("content-type", c.Header().Get("Content-Type")).Trace("write")
	return c.ResponseWriter.Write(b)
}

func ContentType(typeDefault string) func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			h.ServeHTTP(contentWriter{
				ResponseWriter: w,
				def:            typeDefault, log: logging.FromContext(r.Context()),
			}, r)
		})
	}
}
