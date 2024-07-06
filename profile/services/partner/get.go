package partner

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	ppb "brank.as/petnet/gunk/dsa/v2/partner"
)

func (s *Svc) GetPartners(ctx context.Context, req *ppb.GetPartnersRequest) (*ppb.GetPartnersResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pnr, err := s.core.GetPartners(ctx, req.GetOrgID())
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to get partner record")
	}
	return &ppb.GetPartnersResponse{Partners: pnr}, nil
}

func (s *Svc) GetPartner(ctx context.Context, req *ppb.GetPartnersRequest) (*ppb.GetPartnerResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, validation.Required, is.UUIDv4),
		validation.Field(&req.Type, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	pnr, err := s.core.GetPartner(ctx, req.GetOrgID(), req.GetType())
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to get partner record")
	}
	return &ppb.GetPartnerResponse{Partner: &ppb.Partner{
		ID:        pnr.ID,
		OrgID:     pnr.OrgID,
		Type:      pnr.Type,
		Partner:   pnr.Partner,
		Created:   timestamppb.New(pnr.Created),
		Updated:   timestamppb.New(pnr.Updated),
		UpdatedBy: pnr.UpdatedBy,
		Status:    pnr.Status,
	}}, nil
}
