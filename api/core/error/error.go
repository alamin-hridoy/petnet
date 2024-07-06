package error

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/integration/perahub"
)

const (
	MsgDatabaseError             = "Database error"
	MsgPerahubDatabaseError      = "Perahub Database error"
	MsgPartnerDoesntExist        = "Partner doesn't exist"
	MsgInvalidInput              = "Invalid input for parameters"
	MsgControlNumberNotFound     = "Control number not found"
	MsgIdentifierAlreadyExists   = "Identifier already exists"
	MsgConnectionError           = "Connection error"
	MsgTransactionAlreadyClaimed = "Transaction Already Claimed"
	MsgServiceNotAvailableWise   = "Service not available for Wise"
	MsgDRPInternalError          = "DRP internal error"
	MsgPerahubInternalError      = "Perahub internal error"
	MsgNeedToConfirm             = "Need to confirm first"
	MsgUnknownError              = "Unknown Error"
	MsgNotFound                  = "Not found"
)

// Error is error returned fro core
type Error struct {
	Code    codes.Code
	Message string
}

// Error ...
func (r *Error) Error() string {
	if r == nil {
		return ""
	}

	return r.Message
}

// NewCoreError ...
func NewCoreError(code codes.Code, msg string) *Error {
	return &Error{
		Code:    code,
		Message: msg,
	}
}

// ToCoreError ...
func ToCoreError(err error) *Error {
	switch err.(type) {
	case *Error:
		cErr, ok := err.(*Error)
		if !ok || cErr == nil {
			return NewCoreError(codes.Internal, MsgDRPInternalError)
		}

		return cErr

	case *perahub.Error:
		pErr, ok := err.(*perahub.Error)
		if !ok || pErr == nil {
			return NewCoreError(codes.Internal, MsgPerahubInternalError)
		}

		msg := pErr.Msg
		if pErr.UnknownErr != "" {
			msg = pErr.UnknownErr
		}

		if pErr.GRPCCode == codes.Internal {
			msg = MsgPerahubInternalError
		}

		return NewCoreError(pErr.GRPCCode, msg)

	case interface{ GRPCStatus() *status.Status }:
		grpcErr, ok := err.(interface{ GRPCStatus() *status.Status })
		if !ok || grpcErr == nil {
			return NewCoreError(codes.Internal, MsgPerahubInternalError)
		}

		return NewCoreError(grpcErr.GRPCStatus().Code(), grpcErr.GRPCStatus().Message())

	default:
		return NewCoreError(codes.Internal, MsgPerahubInternalError)
	}
}
