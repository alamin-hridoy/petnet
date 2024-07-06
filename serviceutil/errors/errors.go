// Package errors is a thin layer on top of github.com/pkg/errors, which also
// maintains the GRPC error interface when wrapping errors. That enables code to
// wrap GRPC errors without the status information, such as the status code,
// being lost.
package errors

import (
	"fmt"
	"strings"

	"github.com/pkg/errors"
	"google.golang.org/grpc/status"
)

// List is a list of errors, which is easy to use, while also implementing the
// error interface properly.
type List []error

func (l List) Error() string {
	switch len(l) {
	case 0:
		return ""
	case 1:
		return l[0].Error()
	default:
		var b strings.Builder
		fmt.Fprintf(&b, "%d errors:\n", len(l))
		for _, err := range l {
			fmt.Fprintln(&b, err)
		}
		return b.String()
	}
}

// New is an alias for github.com/pkg/errors.New.
func New(message string) error {
	return errors.New(message)
}

// Cause is an alias for github.com/pkg/errors.Cause.
func Cause(err error) error {
	return errors.Cause(err)
}

type Causer interface {
	Cause() error
}

type StackTracer interface {
	StackTrace() errors.StackTrace
}

type wrapError interface {
	error
	// It is important to also keep the other wrapping methods from
	// pkg/errors. Otherwise, stack traces and other data may also be
	// hidden.
	fmt.Formatter
	Causer
	StackTracer
}

type GRPCStatuser interface {
	GRPCStatus() *status.Status
}

type exposedGRPCStatus struct {
	wrapError
	GRPCStatuser
}

// ExposedGRPCStatus returns an err containing a grpc status, if its cause
// contained one.
func ExposedGRPCStatus(err error) error {
	gerr, ok := errors.Cause(err).(GRPCStatuser)
	if !ok {
		return err
	}
	return exposedGRPCStatus{err.(wrapError), gerr}
}

// Wrap is like github.com/pkg/errors.Wrap, but it keeps the grpc status.
func Wrap(err error, message string) error {
	return ExposedGRPCStatus(errors.Wrap(err, message))
}

// Wrapf is like github.com/pkg/errors.Wrapf, but it keeps the grpc status.
func Wrapf(err error, format string, args ...interface{}) error {
	return ExposedGRPCStatus(errors.Wrapf(err, format, args...))
}
