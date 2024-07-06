package revenue_commission

import (
	"context"
	"encoding/json"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"

	revcom_int "brank.as/petnet/api/integration/revenue-commission"
	"brank.as/petnet/api/util"
	revcom "brank.as/petnet/gunk/drp/v1/revenue-commission"
	"brank.as/petnet/serviceutil/logging"
)

// CreateDSACommissionTier Creates DSA Commission Tier record.
func (s *Svc) CreateDSACommissionTier(ctx context.Context, req *revcom.CreateDSACommissionTierRequest) (*revcom.DSACommissionTier, error) {
	log := logging.FromContext(ctx)
	err := validation.ValidateStruct(req,
		validation.Field(&req.TierNo, validation.Required),
		validation.Field(&req.Maximum, validation.Required, is.Float),
		validation.Field(&req.Minimum, validation.Required, is.Float),
		validation.Field(&req.UpdatedBy, validation.Required),
	)
	if err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.dsaCommissionStore.CreateCommissionTier(ctx, &revcom_int.SaveCommissionTierRequest{
		TierNo:    req.TierNo,
		Minimum:   json.Number(req.Minimum),
		Maximum:   json.Number(req.Maximum),
		UpdatedBy: req.UpdatedBy,
	})

	return toDSACommissionTierResponse(log, res, err)
}

// UpdateDSACommissionTier Update DSA Commission Tier record by ID.
func (s *Svc) UpdateDSACommissionTier(ctx context.Context, req *revcom.UpdateDSACommissionTierRequest) (*revcom.DSACommissionTier, error) {
	log := logging.FromContext(ctx)
	err := validation.ValidateStruct(req,
		validation.Field(&req.TierID, validation.Required),
		validation.Field(&req.TierNo, validation.Required),
		validation.Field(&req.Maximum, validation.Required, is.Float),
		validation.Field(&req.Minimum, validation.Required, is.Float),
		validation.Field(&req.UpdatedBy, validation.Required),
	)
	if err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.dsaCommissionStore.UpdateCommissionTier(ctx, &revcom_int.SaveCommissionTierRequest{
		TierID:    req.TierID,
		TierNo:    req.TierNo,
		Minimum:   json.Number(req.Minimum),
		Maximum:   json.Number(req.Maximum),
		UpdatedBy: req.UpdatedBy,
	})

	return toDSACommissionTierResponse(log, res, err)
}

// GetDSACommissionTierByID Gets DSA Commission Tier record by ID.
func (s *Svc) GetDSACommissionTierByID(ctx context.Context, req *revcom.GetDSACommissionTierByIDRequest) (*revcom.DSACommissionTier, error) {
	if req == nil || req.TierID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "tierID is required")
	}

	log := logging.FromContext(ctx)
	res, err := s.dsaCommissionStore.GetCommissionTierByID(ctx, req.TierID)

	return toDSACommissionTierResponse(log, res, err)
}

// ListDSACommissionTier List all DSA Commission Tier records.
func (s *Svc) ListDSACommissionTier(ctx context.Context, r *emptypb.Empty) (*revcom.ListDSACommissionTierResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.dsaCommissionStore.ListCommissionTier(ctx)
	if err != nil {
		logging.WithError(err, log).Error("dsa commission store error")
		return nil, util.HandleServiceErr(err)
	}

	list := make([]*revcom.DSACommissionTier, 0, len(res))
	for _, c := range res {
		list = append(list, toDSACommissionTier(&c))
	}

	return &revcom.ListDSACommissionTierResponse{
		CommissionTierList: list,
	}, nil
}

// DeleteDSACommissionTier Deletes DSA Commission Tier record by ID.
func (s *Svc) DeleteDSACommissionTier(ctx context.Context, req *revcom.DeleteDSACommissionTierRequest) (*revcom.DSACommissionTier, error) {
	if req == nil || req.TierID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "tierID is required")
	}

	log := logging.FromContext(ctx)
	res, err := s.dsaCommissionStore.DeleteCommissionTier(ctx, req.TierID)

	return toDSACommissionTierResponse(log, res, err)
}

func toDSACommissionTierResponse(log *logrus.Entry, c *revcom_int.CommissionTier, err error) (*revcom.DSACommissionTier, error) {
	if err != nil {
		logging.WithError(err, log).Error("dsa commission store error")
		return nil, util.HandleServiceErr(err)
	}

	return toDSACommissionTier(c), nil
}

func toDSACommissionTier(c *revcom_int.CommissionTier) *revcom.DSACommissionTier {
	tierID, _ := c.TierID.Int64()
	var createdAt, updatedAt *timestamppb.Timestamp
	if c.CreatedAt != nil {
		createdAt = timestamppb.New(*c.CreatedAt)
	}

	if c.UpdatedAt != nil {
		updatedAt = timestamppb.New(*c.UpdatedAt)
	}

	return &revcom.DSACommissionTier{
		TierID:    uint32(tierID),
		TierNo:    c.TierNo,
		Minimum:   c.Minimum.String(),
		Maximum:   c.Maximum.String(),
		UpdatedBy: c.UpdatedBy,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}
}
