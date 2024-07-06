package user

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	pnpb "brank.as/petnet/gunk/drp/v1/partner"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) GetProfile(ctx context.Context, req *ppb.GetProfileRequest) (*ppb.GetProfileResponse, error) {
	log := logging.FromContext(ctx)
	pn := req.GetRemitPartner()

	if !static.PartnerExists(pn, "PH") {
		log.Error("partner doesn't exist")
		return nil, status.Error(codes.NotFound, "partner doesn't exist")
	}

	vp, ok := s.validators[pn]
	if !ok {
		return nil, status.Error(codes.NotFound, fmt.Sprint("missing service for ", pn))
	}

	creq, err := vp.GetProfileValidate(ctx, req)
	if err != nil {
		return nil, err
	}

	res, err := s.user.GetProfile(ctx, *creq)
	if err != nil {
		switch t := err.(type) {
		case *perahub.Error:
			if t.Type == perahub.PartnerError {
				return nil, perahub.GRPCError(t.GRPCCode, "partner error", &pnpb.Error{
					Code:    t.Code,
					Message: t.Msg,
				})
			}
			return nil, status.Errorf(codes.Internal, "internal error occurred")
		}
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Errorf(codes.Internal, "internal error occurred")
	}
	return &ppb.GetProfileResponse{
		Profile: &ppb.Profile{
			ID:         res.ID,
			Type:       res.Type,
			FirstName:  res.FirstName,
			LastName:   res.LastName,
			BirthDate:  res.BirthDate,
			Phone:      res.Phone,
			Occupation: res.Occupation,
			Address: &ppb.Address{
				Address1:   res.Address.Address1,
				City:       res.Address.City,
				State:      res.Address.State,
				PostalCode: res.Address.PostalCode,
				Country:    res.Address.Country,
			},
		},
	}, nil
}

func (*WISEVal) GetProfileValidate(ctx context.Context, req *ppb.GetProfileRequest) (*core.GetProfileReq, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.Required, is.Email),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &core.GetProfileReq{
		Email: req.Email,
	}, nil
}

func (*CEBVal) GetProfileValidate(ctx context.Context, req *ppb.GetProfileRequest) (*core.GetProfileReq, error) {
	return nil, fmt.Errorf("service not available for Cebuana")
}
