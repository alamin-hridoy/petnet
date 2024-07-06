package session

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	authpb "brank.as/rbac/gunk/v1/authenticate"
)

// Error creates an error with details.
func Error(c codes.Code, msg string, details *authpb.SessionError) error {
	p := status.New(c, msg)
	if e, err := p.WithDetails(details); err == nil {
		return e.Err()
	}
	return p.Err()
}

func FromError(err error) *authpb.SessionError {
	s, ok := status.FromError(err)
	if !ok {
		return nil
	}
	if len(s.Proto().GetDetails()) == 0 {
		return nil
	}
	det := &authpb.SessionError{}
	if err := s.Proto().Details[0].UnmarshalTo(det); err != nil {
		return nil
	}
	return det
}
