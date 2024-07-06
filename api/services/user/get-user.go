package user

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	pnpb "brank.as/petnet/gunk/drp/v1/partner"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) GetUser(ctx context.Context, req *ppb.GetUserRequest) (*ppb.GetUserResponse, error) {
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

	creq, err := vp.GetUserValidate(ctx, req)
	if err != nil {
		return &ppb.GetUserResponse{
			Message: err.Error(),
		}, nil
	}

	res, err := s.user.GetUser(ctx, *creq)
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
	return &ppb.GetUserResponse{
		Code:    int32(res.Code),
		Message: res.Message,
		Result: &ppb.GUResult{
			User: &ppb.User{
				UserID:         int32(res.Result.Client.ClientID),
				UserNumber:     res.Result.Client.ClientNumber,
				FirstName:      res.Result.Client.FirstName,
				MiddleName:     res.Result.Client.MiddleName,
				LastName:       res.Result.Client.LastName,
				BirthDate:      res.Result.Client.BirthDate,
				MobileCountry:  int32(res.Result.Client.CPCountry.CountryID),
				PhoneCountry:   int32(res.Result.Client.CPCountry.CountryID),
				CountryAddress: int32(res.Result.Client.CPCountry.CountryID),
				SourceOfFund:   int32(res.Result.Client.CSOfFund.SourceOfFundID),
			},
		},
		RemcoID: int32(res.RemcoID),
	}, nil
}

func (*WISEVal) GetUserValidate(ctx context.Context, req *ppb.GetUserRequest) (*core.GetUserRequest, error) {
	return nil, fmt.Errorf("service not available for WISE")
}

func (*CEBVal) GetUserValidate(ctx context.Context, req *ppb.GetUserRequest) (*core.GetUserRequest, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.FirstName, validation.Required),
		validation.Field(&req.LastName, validation.Required),
		validation.Field(&req.BirthDate, validation.Required),
		validation.Field(&req.UserNumber, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	return &core.GetUserRequest{
		FirstName:    req.FirstName,
		LastName:     req.LastName,
		BirthDate:    req.BirthDate,
		ClientNumber: req.UserNumber,
	}, nil
}
