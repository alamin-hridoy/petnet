package signup

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/storage"

	upb "brank.as/rbac/gunk/v1/user"
)

// Signup signs up a user using the passed in details
func (h *Handler) Signup(ctx context.Context, req *upb.SignupRequest) (*upb.SignupResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.signup.Signup")

	rqr := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Username,
			validation.When(req.GetInviteCode() == "", rqr, validation.Length(1, 70))),
		validation.Field(&req.Email,
			validation.When(req.GetInviteCode() == "", is.Email, rqr, validation.Length(3, 254))),
		validation.Field(&req.FirstName, rqr, validation.Length(1, 70)),
		validation.Field(&req.LastName, rqr, validation.Length(1, 70)),
		validation.Field(&req.Password, rqr, validation.Length(8, 64)),
	); err != nil {
		log.WithError(err).Error("invalid request")
		return nil, status.Errorf(codes.InvalidArgument, err.Error())
	}

	dbUser := storage.User{
		OrgID:      req.OrgID,
		Username:   req.Username,
		FirstName:  req.FirstName,
		LastName:   req.LastName,
		Email:      req.Email,
		InviteCode: req.InviteCode,
	}
	u, code, err := h.store.CreateUser(ctx, dbUser, storage.Credential{
		Username: req.Username,
		Password: req.Password,
	})
	if err != nil {
		return nil, err
	}

	if code != "" {
		if err := h.mailer.ConfirmEmail(u.Email, code); err != nil {
			log.WithError(err).Error("sending confirm email")
		}
	}
	return &upb.SignupResponse{
		UserID: u.ID,
		OrgID:  u.OrgID,
	}, nil
}
