package util

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"google.golang.org/grpc/codes"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
)

func TestErrorHandler_HTTPErrorHandler(t *testing.T) {
	mux := &runtime.ServeMux{}

	req := httptest.NewRequest("GET", "http://example.com", nil)
	marshaler := &runtime.JSONPb{
		UnmarshalOptions: protojson.UnmarshalOptions{DiscardUnknown: true},
	}

	eh := ErrorHandler{
		Log: logrus.NewEntry(logrus.StandardLogger()),
	}

	for _, tc := range []struct {
		name             string
		inputStatus      *status.Status
		expectedHttpCode int
		expectedBodyStr  string
	}{
		{
			name:             "no error",
			inputStatus:      status.New(codes.OK, "success"),
			expectedHttpCode: http.StatusOK,
			expectedBodyStr:  "",
		},
		{
			name:             "valid grpc code",
			inputStatus:      status.New(codes.InvalidArgument, "invalid args"),
			expectedHttpCode: http.StatusBadRequest,
			expectedBodyStr:  "invalid args",
		},
		{
			name:             "http code as grpc code",
			inputStatus:      status.New(422, "validation error"),
			expectedHttpCode: http.StatusUnprocessableEntity,
			expectedBodyStr:  "validation error",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			resRecorder := httptest.NewRecorder()

			err := tc.inputStatus.Err()

			eh.HTTPErrorHandler(context.Background(), mux, marshaler, resRecorder, req, err)

			if tc.expectedHttpCode != resRecorder.Code {
				t.Errorf("http code does not match.\nexpected: %d, got: %d", tc.expectedHttpCode, resRecorder.Code)
			}

			bodyStr := resRecorder.Body.String()
			if tc.expectedBodyStr != "" && !strings.Contains(bodyStr, tc.expectedBodyStr) {
				t.Errorf("body does not contain expected message.\nexpected: %s\ngot: %s\n", tc.expectedBodyStr, bodyStr)
			}
		})
	}
}
