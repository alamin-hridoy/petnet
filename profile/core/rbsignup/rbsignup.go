package rbsignup

import (
	"context"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	spb "brank.as/rbac/gunk/v1/user"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SignupReq struct {
	Username   string
	FirstName  string
	LastName   string
	Email      string
	Password   string
	InviteCode string
	OrgID      string
}

type SignupResp struct {
	UserID string
	OrgID  string
}

type Svc struct {
	spb.SignupClient
}

func New(scl spb.SignupClient) *Svc {
	return &Svc{
		SignupClient: scl,
	}
}

// ResetPassword after confirmation.
func (s *Svc) ResetPassword(ctx context.Context, req *spb.ResetPasswordRequest,
	opts ...grpc.CallOption,
) (*spb.ResetPasswordResponse, error) {
	log := logging.FromContext(ctx)
	rpr, err := s.SignupClient.ResetPassword(ctx, req, opts...)
	if err != nil {
		logging.WithError(err, log).Error("change password with code")
		if err == storage.NotFound {
			return nil, status.Error(codes.InvalidArgument, "invalid code")
		}
		return nil, err
	}
	return rpr, nil
}
