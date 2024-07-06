package perahub

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
)

func TestHandleNonexError(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name      string
		errStruct PtnrErrStruct
		httpCode  int
		want      *Error
	}{
		{
			name:      "E1",
			errStruct: E1,
			httpCode:  400,
			want: &Error{
				Code:     "2001",
				GRPCCode: codes.InvalidArgument,
				Msg:      "M1",
				Type:     PartnerError,
			},
		},
		{
			name:      "E2",
			errStruct: E2,
			httpCode:  422,
			want: &Error{
				Code:     "422",
				GRPCCode: codes.InvalidArgument,
				Msg:      "The given data was invalid.",
				Type:     PartnerError,
				Errors: map[string][]string{
					"agent_code": {
						"The agent code field is required.",
					},
				},
			},
		},
		{
			name:      "E3",
			errStruct: E3,
			httpCode:  400,
			want: &Error{
				Code:     "2003",
				GRPCCode: codes.InvalidArgument,
				Msg:      "E3",
				Type:     PartnerError,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			v := PtnrErrStructs[test.errStruct]
			b, err := json.Marshal(&v.Error)
			if err != nil {
				t.Fatal(err)
			}

			err = handleNonexErr(ctx, b, "error-test", test.httpCode)
			if err == nil {
				t.Fatal("want error got nil")
			}

			tp, ok := err.(*Error)
			if !ok {
				t.Fatal("unknown error type")
			}

			if !cmp.Equal(test.want, tp) {
				t.Error(cmp.Diff(test.want, tp))
			}
		})
	}
}

func TestHandleBillsPaymentErr(t *testing.T) {
	t.Parallel()
	ctx := context.Background()

	tests := []struct {
		name      string
		errStruct BillPayErrStruct
		httpCode  int
		want      *Error
	}{
		{
			name:      "BP1",
			errStruct: BP1,
			httpCode:  400,
			want: &Error{
				Code:     "2001",
				GRPCCode: codes.InvalidArgument,
				Msg:      "M1",
				Type:     BillPayError,
			},
		},
		{
			name:      "BP2",
			errStruct: BP2,
			httpCode:  422,
			want: &Error{
				Code:     "422",
				GRPCCode: codes.InvalidArgument,
				Msg:      "validation_error",
				Type:     BillPayError,
				Errors: map[string][]string{
					"code": {
						"The code is required",
					},
				},
			},
		},
		{
			name:      "BP3",
			errStruct: BP3,
			httpCode:  400,
			want: &Error{
				Code:     "2003",
				GRPCCode: codes.InvalidArgument,
				Msg:      "BP3",
				Type:     BillPayError,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			v := BillPayErrStructs[test.errStruct]
			b, err := json.Marshal(&v.Error)
			if err != nil {
				t.Fatal(err)
			}

			err = handleBillsPaymentErr(ctx, b, "error-test", test.httpCode)
			if err == nil {
				t.Fatal("want error got nil")
			}

			tp, ok := err.(*Error)
			if !ok {
				t.Fatal("unknown error type")
			}

			if !cmp.Equal(test.want, tp) {
				t.Error(cmp.Diff(test.want, tp))
			}
		})
	}
}

func TestErrors(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name       string
		want       *Error
		nonexErr   bool
		drpErr     bool
		partnerErr bool
	}{
		{
			name: "Nonex Error",
			want: &Error{
				Code:       "200",
				GRPCCode:   codes.Internal,
				Msg:        codes.Internal.String(),
				UnknownErr: "<html>server-error<html>",
				Type:       NonexError,
			},
			nonexErr: true,
		},
		{
			name: "DRP Error",
			want: &Error{
				Code:       "400",
				GRPCCode:   codes.InvalidArgument,
				Msg:        codes.InvalidArgument.String(),
				UnknownErr: `{"non":{"standard":"error"}}`,
				Type:       PartnerError,
			},
			drpErr: true,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()

			s, m := newTestSvc(t, st)
			ctx := context.Background()

			m.nonexErr = test.nonexErr
			m.partnerErr = test.partnerErr
			m.drpErr = test.drpErr

			_, err := s.postNonex(ctx, s.nonexURL("ria/inquire"), RiaInquireRequest{})
			if err == nil {
				t.Fatal(err)
			}

			switch tp := err.(type) {
			case *Error:
				if !cmp.Equal(test.want, tp) {
					t.Error(cmp.Diff(test.want, tp))
				}
			default:
				t.Fatal(err)
			}
		})
	}
}
