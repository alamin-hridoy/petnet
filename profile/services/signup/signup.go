package signup

import (
	"context"
	"errors"
	"unicode"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	upb "brank.as/petnet/gunk/dsa/v1/user"
	core "brank.as/petnet/profile/core/signup"
)

var (
	errMissingSpecial  = errors.New("Password needs at least one special character")
	errMissingUpper    = errors.New("Password needs at least one uppercase character")
	errMissingNumber   = errors.New("Password needs at least one number")
	errIncorrectLength = errors.New("Password should be between 8 and 64 characters long")
)

func (s *Svc) Signup(ctx context.Context, req *upb.SignupRequest) (*upb.SignupResponse, error) {
	if err := validate(req); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.core.Signup(ctx, core.SignupReq{
		Username:   req.GetUsername(),
		FirstName:  req.GetFirstName(),
		LastName:   req.GetLastName(),
		Email:      req.GetEmail(),
		Password:   req.GetPassword(),
		InviteCode: req.GetInviteCode(),
		OrgID:      req.GetOrgID(),
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to signup")
	}
	return &upb.SignupResponse{
		UserID: res.UserID,
		OrgID:  res.OrgID,
	}, nil
}

func validate(req *upb.SignupRequest) error {
	if req.GetInviteCode() == "" {
		validation.Validate(&req.Username, validation.Required, validation.Length(1, 70))
		validation.Validate(&req.Email, validation.Required, is.Email, validation.Required, validation.Length(3, 254))
	}
	if err := validation.ValidateStruct(req,
		validation.Field(&req.FirstName, validation.Required, validation.Length(1, 70)),
		validation.Field(&req.LastName, validation.Required, validation.Length(1, 70)),
	); err != nil {
		return err
	}

	pwd := req.GetPassword()
	if len(pwd) < 8 || len(pwd) > 64 {
		return errIncorrectLength
	}

	var hasNumber, hasUpper, hasSpecial bool
	for _, c := range pwd {
		switch {
		case unicode.IsNumber(c):
			hasNumber = true
		case unicode.IsUpper(c):
			hasUpper = true
		case unicode.IsPunct(c) || unicode.IsSymbol(c):
			hasSpecial = true
		}
	}
	if !hasNumber {
		return errMissingNumber
	}
	if !hasUpper {
		return errMissingUpper
	}
	if !hasSpecial {
		return errMissingSpecial
	}
	return nil
}
