package util

import (
	"fmt"
	"net/http"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
)

var notFoundMsg = []string{
	"Control number does not exist",
	"Transaction does not exist",
	"Transaction not found",
	"Not Found",
	"Transfer not found",
	"No order found matching supplied PIN and BeneficiaryAmount",
	"Invalid Transaction, please contact headoffice",
	"Sorry but the U Cash Padala Control Number was invalid",
	"web service during transaction lookup",
	"NO TRANSACTION FOUND",
}

var alreadyClaimMsg = []string{
	"Transaction Already Claimed",
	"is already Paid",
	"Order has already been marked as paid to beneficiary",
	"Claimed already",
	"Paid",
	"U CASH PADALA CONTROL NUMBER IS NOT VALID",
	"Transfer already PAID by different partner",
}

// HandleServiceErr ...
// TODO(vitthal): Move this in Httphandler middleware after moving all specific errors to partners
func HandleServiceErr(err error) error {
	if err == nil {
		return nil
	}

	httpStatus, message := toHttpStatusMessage(err)
	// TODO(vitthal): add details core to error
	if httpStatus == http.StatusUnprocessableEntity || strings.Contains(message, "The given data was invalid") {
		// TODO(vitthal): move this under perahub errors
		message = fmt.Sprintf(coreerror.MsgInvalidInput+" :(%s)", message)
	}

	if strings.Contains(message, "pq: ") {
		message = coreerror.MsgDatabaseError
		httpStatus = http.StatusInternalServerError
	}

	// TODO: Move to specific partner
	if strings.Contains(message, "service not available for Wise") {
		message = coreerror.MsgServiceNotAvailableWise
		httpStatus = http.StatusInternalServerError
	}

	// TODO: check whether its internal error or should be service not available?
	if strings.Contains(message, "json: cannot unmarshal object into Go struct field") {
		message = coreerror.MsgConnectionError
		httpStatus = http.StatusServiceUnavailable
	}

	// TODO(vitthal): move bellow error checking in core partners
	for _, v := range alreadyClaimMsg {
		if strings.Contains(message, v) {
			message = coreerror.MsgTransactionAlreadyClaimed
			httpStatus = http.StatusConflict
			break
		}
	}

	for _, v := range notFoundMsg {
		if strings.Contains(message, v) {
			message = coreerror.MsgControlNumberNotFound
			httpStatus = http.StatusNotFound
			break
		}
	}

	return status.Error(codes.Code(httpStatus), message)
}

func getHttpStatusCodeFromGrpcCode(code codes.Code) int {
	switch code {
	case codes.Unauthenticated:
		return http.StatusUnauthorized
	case codes.DataLoss:
		return http.StatusInternalServerError
	case codes.Unavailable:
		return http.StatusServiceUnavailable
	case codes.Internal:
		return http.StatusInternalServerError
	case codes.Unimplemented:
		return http.StatusNotImplemented
	case codes.OutOfRange:
		return http.StatusServiceUnavailable
	case codes.Aborted:
		return http.StatusNotAcceptable
	case codes.FailedPrecondition:
		return http.StatusPreconditionFailed
	case codes.ResourceExhausted:
		return http.StatusTooManyRequests
	case codes.PermissionDenied:
		return http.StatusForbidden
	case codes.AlreadyExists:
		return http.StatusConflict
	case codes.NotFound:
		return http.StatusNotFound
	case codes.DeadlineExceeded:
		return http.StatusRequestTimeout
	case codes.InvalidArgument:
		return http.StatusBadRequest
	case codes.Canceled:
		return http.StatusRequestTimeout
	case codes.Unknown:
		return http.StatusInternalServerError
	default:
		return http.StatusInternalServerError
	}
}

func toHttpStatusMessage(err error) (int, string) {
	var (
		grpcCode codes.Code
		message  string
	)

	switch err.(type) {
	case *coreerror.Error:
		if cErr, ok := err.(*coreerror.Error); ok && cErr != nil {
			grpcCode = cErr.Code
			message = cErr.Message
		}

	case interface{ GRPCStatus() *status.Status }:
		if grpcErr, ok := err.(interface{ GRPCStatus() *status.Status }); ok && grpcErr != nil {
			st := grpcErr.GRPCStatus()
			if st != nil {
				grpcCode = st.Code()
				message = st.Message()
			}
		}

	case *perahub.Error:
		if pErr, ok := err.(*perahub.Error); ok && pErr != nil {
			msg := pErr.Msg
			if pErr.UnknownErr != "" {
				msg = pErr.UnknownErr
			}

			grpcCode = pErr.GRPCCode
			message = msg
		}

	default:
		grpcCode = codes.Unknown
		message = err.Error()
	}

	// In case code and message not set. Though in very rare case
	if grpcCode <= 0 {
		grpcCode = codes.Unknown
	}

	if message == "" {
		message = err.Error()
	}

	return getHttpStatusCodeFromGrpcCode(grpcCode), message
}
