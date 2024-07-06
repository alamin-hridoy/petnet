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

// ForgotPassword is for sending a password reset mail containing
func (h *Handler) ForgotPassword(ctx context.Context, req *upb.ForgotPasswordRequest) (*upb.ForgotPasswordResponse, error) {
	log := logging.FromContext(ctx).WithField("service", "signup.ForgotPasswordSend")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.Required, is.Email),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := h.store.ResetPasswordInit(ctx, req.Email); err != nil {
		return nil, err
	}
	return &upb.ForgotPasswordResponse{}, nil
}

// ResetPassword is for resetting the user's password with the sent token
func (h *Handler) ResetPassword(ctx context.Context, req *upb.ResetPasswordRequest) (*upb.ResetPasswordResponse, error) {
	log := logging.FromContext(ctx).WithField("service", "signup.EmailConfirmation")

	if err := validation.ValidateStruct(req,
		validation.Field(&req.Password, validation.Required, validation.Length(1, 64)),
		validation.Field(&req.Code, validation.Required),
	); err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	if err := h.store.ResetPassword(ctx, req.Code, req.Password); err != nil {
		return nil, err
	}

	return &upb.ResetPasswordResponse{}, nil
}
