package service

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	spb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) ValidateServiceAccess(ctx context.Context, req *spb.ValidateServiceAccessRequest) (*spb.ValidateServiceAccessResponse, error) {
	log := logging.FromContext(ctx)
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, required, is.UUID),
		validation.Field(&req.Partner, validation.When(!req.IsAnyPartnerEnabled, required)),
		validation.Field(&req.SvcName, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	res, err := s.st.ValidateSvcRequest(ctx, storage.ValidateSvcRequestFilter{
		OrgID:               req.GetOrgID(),
		Partner:             req.GetPartner(),
		SvcName:             req.GetSvcName(),
		IsAnyPartnerEnabled: req.GetIsAnyPartnerEnabled(),
	})
	if err != nil {
		logging.WithError(err, log).Error("Validate Service Access")
		return &spb.ValidateServiceAccessResponse{
			Enabled: false,
		}, nil
	}
	return &spb.ValidateServiceAccessResponse{
		Enabled: res.Enabled,
	}, nil
}
