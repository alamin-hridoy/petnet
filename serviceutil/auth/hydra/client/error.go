package client

import (
	"errors"
	"fmt"
)

type ErrorCode int

const (
	_ ErrorCode = iota
	HydraError
	InvalidOptions
	InvalidClientParam
)

type Error struct {
	Message string
	Err     error
	Code    ErrorCode
}

var emptyOptions = Error{
	Message: `host or host pattern cannot be empty`,
	Code:    InvalidOptions,
}

var invalidOptions = Error{
	Message: `host and hostPattern cannot be set at the same time`,
	Code:    InvalidOptions,
}

var hostAlreadySet = Error{
	Message: `namespace can't be used. Host is already set.`,
	Code:    InvalidOptions,
}

var missingClientID = Error{
	Message: `clientID is required`,
	Code:    InvalidClientParam,
}

var missingURI = Error{
	Message: `redirectURI is required`,
	Code:    InvalidClientParam,
}

// Error fulfils the error interface.
func (e Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s, err: %s", e.Message, e.Err)
	}
	return e.Message
}

// Cause of the underlying error.  Used for internal logging.
func (e Error) Cause() error {
	if e.Err == nil {
		return errors.New(e.Message)
	}
	return e.Err
}
