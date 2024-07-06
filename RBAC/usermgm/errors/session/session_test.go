package session_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc/codes"
	"google.golang.org/protobuf/testing/protocmp"

	"brank.as/rbac/usermgm/errors/session"

	authpb "brank.as/rbac/gunk/v1/authenticate"
)

func TestError(t *testing.T) {
	want := &authpb.SessionError{
		Message:           "Incorrect username or password.",
		TrackingAttempts:  true,
		RemainingAttempts: 3,
		ErrorDetails: map[string]string{
			"username": "Invalid email or password. Please try again.",
			"email":    "user@email.com",
		},
	}
	e := session.Error(codes.AlreadyExists, "test error", want)
	if e == nil {
		t.Error("nil error value")
	}

	det := session.FromError(e)
	if det == nil {
		t.Fatal("missing details")
	}
	if !cmp.Equal(want, det, protocmp.Transform()) {
		t.Error(cmp.Diff(want, det, protocmp.Transform()))
	}
}
