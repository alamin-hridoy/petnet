package signup

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"

	upb "brank.as/rbac/gunk/v1/user"
)

// ResendConfirmEmail resends the confirmation email to the user
func (h *Handler) ResendConfirmEmail(ctx context.Context, req *upb.ResendConfirmEmailRequest) (*upb.ResendConfirmEmailResponse, error) {
	log := logging.FromContext(ctx).WithField("service", "signup.ResendConfirmEmail")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.UserID, validation.Required, is.UUIDv4),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	u, code, err := h.store.GetConfirmationCode(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	if err := h.mailer.ConfirmEmail(u.Email, code); err != nil {
		logging.WithError(err, log).Error("sending confirm email")
		return nil, err
	}
	return &upb.ResendConfirmEmailResponse{}, nil
}
