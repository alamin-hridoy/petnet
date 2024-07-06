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

func (s *Svc) RegisterUser(ctx context.Context, req *ppb.RegisterUserRequest) (*ppb.RegisterUserResponse, error) {
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

	err := vp.RegisterUserValidate(ctx, req)
	if err != nil {
		return nil, err
	}

	res, err := s.user.RegisterUser(ctx, core.RegisterUserReq{
		Email:     req.GetEmail(),
		FirstName: req.GetFirstName(),
		LastName:  req.GetLastName(),
		BirthDate: req.GetBirthDate(),
		CpCtryID:  req.GetMobileCountryID(),
		ContactNo: req.GetContactNo(),
		TpCtryID:  req.GetPhoneCountryID(),
		TpArCode:  req.GetPhoneArCode(),
		CrtyAdID:  req.GetCountryAddID(),
		PAdd:      req.GetProvinceAdd(),
		CAdd:      req.GetCurrentAdd(),
		UserID:    req.GetUserID(),
		SOFID:     req.GetSourceOfID(),
		Tin:       req.GetTin(),
		TpNo:      req.GetPhoneNo(),
		AgentCode: req.GetAgentCode(),
	})
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
	return &ppb.RegisterUserResponse{
		Code:    int32(res.Code),
		Message: res.Message,
		Result: &ppb.RUResult{
			ResultStatus: res.Result.ResultStatus,
			MessageID:    int32(res.Result.MessageID),
			LogID:        int32(res.Result.LogID),
			UserID:       int32(res.Result.ClientID),
			UserNumber:   res.Result.ClientNo,
		},
		RemcoID: int32(res.RemcoID),
	}, nil
}

func (*WISEVal) RegisterUserValidate(ctx context.Context, req *ppb.RegisterUserRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Email, validation.Required, is.Email),
	); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}

func (*CEBVal) RegisterUserValidate(ctx context.Context, req *ppb.RegisterUserRequest) error {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.FirstName, validation.Required),
		validation.Field(&req.LastName, validation.Required),
		validation.Field(&req.BirthDate, validation.Required),
		validation.Field(&req.MobileCountryID, validation.Required),
		validation.Field(&req.ContactNo, validation.Required),
		validation.Field(&req.PhoneCountryID, validation.Required),
		validation.Field(&req.PhoneArCode),
		validation.Field(&req.CountryAddID, validation.Required),
		validation.Field(&req.ProvinceAdd),
		validation.Field(&req.CurrentAdd, validation.Required),
		validation.Field(&req.UserID, validation.Required),
		validation.Field(&req.SourceOfID, validation.Required),
		validation.Field(&req.Tin),
		validation.Field(&req.PhoneNo),
		validation.Field(&req.AgentCode, validation.Required),
	); err != nil {
		return status.Error(codes.InvalidArgument, err.Error())
	}
	return nil
}
