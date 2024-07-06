package user

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/core"

	tspb "google.golang.org/protobuf/types/known/timestamppb"

	mpb "brank.as/rbac/gunk/v1/mfa"
	upb "brank.as/rbac/gunk/v1/user"
)

func (s *Handler) UpdateUser(ctx context.Context, req *upb.UpdateUserRequest) (*upb.UpdateUserResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "svc.user.updateuser")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
		validation.Field(&req.MFAType, validation.By(func(interface{}) error {
			return validation.Validate(mpb.MFA_name[mpb.MFA_value[req.MFAType.String()]],
				validation.Required)
		})),
		validation.Field(&req.Email, is.EmailFormat),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	mfa := req.MFAType.String()
	if req.MFAType == mpb.MFA_PASS {
		// MFA type PASS is only for password re-validation,
		// not valid for 2nd factor challenges.
		mfa = ""
	}

	m, err := s.usr.UpdateUser(ctx, core.User{
		ID:           req.UserID,
		FName:        req.FirstName,
		LName:        req.LastName,
		Email:        req.Email,
		PreferredMFA: mfa,
		EnableMFA:    req.LoginMFA == upb.EnableOpt_Enable,
		DisableMFA:   req.LoginMFA == upb.EnableOpt_Disable,
	})
	if err != nil {
		logging.WithError(err, log).Error("update user")
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to update record")
	}
	if m != nil {
		return &upb.UpdateUserResponse{
			MFAEventID: m.EventID,
			MFAType:    mpb.MFA(mpb.MFA_value[m.Type]),
		}, nil
	}

	return &upb.UpdateUserResponse{Updated: tspb.Now()}, nil
}
