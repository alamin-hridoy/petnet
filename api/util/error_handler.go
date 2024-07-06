package util

import (
	"context"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
)

type ErrorHandler struct {
	Log *logrus.Entry
}

func NewErrorHandler(log *logrus.Entry) *ErrorHandler {
	return &ErrorHandler{
		Log: log,
	}
}

// HTTPErrorHandler is for handling if custom/http error codes are set other than grpc codes
func (eh *ErrorHandler) HTTPErrorHandler(ctx context.Context, mux *runtime.ServeMux, marshaler runtime.Marshaler, w http.ResponseWriter, r *http.Request, err error) {
	s := status.Convert(err)
	logStr := "error handler: status code less than 300 received"
	if int(s.Code()) >= 300 {
		// If custom codes are set for http response,
		// use HTTPStatusError to setup custom http code by DefaultHTTPErrorHandler
		err = &runtime.HTTPStatusError{
			HTTPStatus: int(s.Code()),
			Err:        err,
		}

		logStr = "error handler: status code greater than 300 received"
	}

	eh.Log.Debug(logStr)

	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}
