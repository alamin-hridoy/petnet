package signup

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	upb "brank.as/rbac/gunk/v1/user"
)

// EmailConfirmation verifies the confirmation code sent to the users email
func (h *Handler) EmailConfirmation(ctx context.Context, req *upb.EmailConfirmationRequest) (*upb.EmailConfirmationResponse, error) {
	log := logging.FromContext(ctx).WithField("service", "signup.EmailConfirmation")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.Code, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	u, err := h.store.ConfirmEmail(ctx, req.Code)
	if err != nil {
		return nil, err
	}
	return &upb.EmailConfirmationResponse{
		Email:     u.Email,
		UserID:    u.ID,
		OrgID:     u.OrgID,
		Username:  u.Username,
		FirstName: u.FirstName,
		LastName:  u.LastName,
	}, nil
}
