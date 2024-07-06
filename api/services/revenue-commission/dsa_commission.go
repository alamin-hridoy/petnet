package revenue_commission

import (
	"context"
	"strconv"
	"time"

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

// CreateDSACommission Creates DSA Commission record.
func (s *Svc) CreateDSACommission(ctx context.Context, req *revcom.CreateDSACommissionRequest) (*revcom.DSACommission, error) {
	log := logging.FromContext(ctx)
	err := validation.ValidateStruct(req,
		validation.Field(&req.DsaCode, validation.Required),
		validation.Field(&req.CommissionType, validation.Required),
		validation.Field(&req.CommissionCurrency, validation.Required, is.Digit),
		validation.Field(&req.CommissionAmount, validation.Required, is.Float),
		validation.Field(&req.UpdatedBy, validation.Required),
		validation.Field(&req.TrxType, validation.Required),
		validation.Field(&req.RemitType, validation.Required),
	)
	if err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	commAmt, _ := strconv.ParseFloat(req.CommissionAmount, 64)
	effectiveDate := ""
	if req.EffectiveDate != nil && req.EffectiveDate.IsValid() {
		effectiveDate = req.EffectiveDate.AsTime().Format("2006-01-02")
	}

	res, err := s.dsaCommissionStore.CreateDSACommission(ctx, &revcom_int.SaveDSACommissionRequest{
		DsaCode:            req.DsaCode,
		CommissionType:     toIntegrationCommissionType(req.CommissionType),
		TierID:             req.TierID,
		CommissionAmount:   commAmt,
		CommissionCurrency: req.CommissionCurrency,
		UpdatedBy:          req.UpdatedBy,
		EffectiveDate:      effectiveDate,
		TrxType:            req.TrxType,
		RemitType:          req.RemitType,
	})

	return toDSACommissionResponse(log, res, err)
}

// UpdateDSACommission Update DSA Commission record by ID.
func (s *Svc) UpdateDSACommission(ctx context.Context, req *revcom.UpdateDSACommissionRequest) (*revcom.DSACommission, error) {
	log := logging.FromContext(ctx)
	err := validation.ValidateStruct(req,
		validation.Field(&req.CommID, validation.Required),
		validation.Field(&req.DsaCode, validation.Required),
		validation.Field(&req.CommissionType, validation.Required),
		validation.Field(&req.CommissionCurrency, validation.Required, is.Digit),
		validation.Field(&req.CommissionAmount, validation.Required, is.Float),
		validation.Field(&req.UpdatedBy, validation.Required),
		validation.Field(&req.TrxType, validation.Required),
		validation.Field(&req.RemitType, validation.Required),
	)
	if err != nil {
		logging.WithError(err, log).Error("invalid request")
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	commAmt, _ := strconv.ParseFloat(req.CommissionAmount, 64)
	effectiveDate := ""
	if req.EffectiveDate != nil && req.EffectiveDate.IsValid() {
		effectiveDate = req.EffectiveDate.AsTime().Format("2006-01-02")
	}

	res, err := s.dsaCommissionStore.UpdateDSACommission(ctx, &revcom_int.SaveDSACommissionRequest{
		ID:                 req.CommID,
		DsaCode:            req.DsaCode,
		CommissionType:     toIntegrationCommissionType(req.CommissionType),
		TierID:             req.TierID,
		CommissionAmount:   commAmt,
		CommissionCurrency: req.CommissionCurrency,
		UpdatedBy:          req.UpdatedBy,
		EffectiveDate:      effectiveDate,
		TrxType:            req.TrxType,
		RemitType:          req.RemitType,
	})

	return toDSACommissionResponse(log, res, err)
}

// GetDSACommissionByID Gets DSA Commission record by ID.
func (s *Svc) GetDSACommissionByID(ctx context.Context, req *revcom.GetDSACommissionByIDRequest) (*revcom.DSACommission, error) {
	if req == nil || req.CommID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	log := logging.FromContext(ctx)
	res, err := s.dsaCommissionStore.GetDSACommissionByID(ctx, req.CommID)

	return toDSACommissionResponse(log, res, err)
}

// ListDSACommission List DSA Commission records.
func (s *Svc) ListDSACommission(ctx context.Context, r *emptypb.Empty) (*revcom.ListDSACommissionResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.dsaCommissionStore.ListDSACommission(ctx)
	if err != nil {
		logging.WithError(err, log).Error("dsa commission store error")
		return nil, util.HandleServiceErr(err)
	}

	list := make([]*revcom.DSACommission, 0, len(res))
	for _, c := range res {
		list = append(list, toDSACommission(&c))
	}

	return &revcom.ListDSACommissionResponse{
		CommissionList: list,
	}, nil
}

// DeleteDSACommission Deletes DSA Commission record by ID.
func (s *Svc) DeleteDSACommission(ctx context.Context, req *revcom.DeleteDSACommissionRequest) (*revcom.DSACommission, error) {
	if req == nil || req.CommID <= 0 {
		return nil, status.Error(codes.InvalidArgument, "id is required")
	}

	log := logging.FromContext(ctx)
	res, err := s.dsaCommissionStore.DeleteDSACommission(ctx, req.CommID)

	return toDSACommissionResponse(log, res, err)
}

func toDSACommissionResponse(log *logrus.Entry, c *revcom_int.DSACommission, err error) (*revcom.DSACommission, error) {
	if err != nil {
		logging.WithError(err, log).Error("dsa commission store error")
		return nil, util.HandleServiceErr(err)
	}

	return toDSACommission(c), nil
}

func toDSACommission(c *revcom_int.DSACommission) *revcom.DSACommission {
	commID, _ := c.ID.Int64()
	tierID, _ := c.TierID.Int64()
	var (
		createdAt *timestamppb.Timestamp
		updatedAt *timestamppb.Timestamp
	)
	if c.CreatedAt != nil {
		createdAt = timestamppb.New(*c.CreatedAt)
	}

	if c.UpdatedAt != nil {
		updatedAt = timestamppb.New(*c.UpdatedAt)
	}

	var effectiveDate *timestamppb.Timestamp
	if c.EffectiveDate != "" {
		t, _ := time.Parse("2006-01-02", c.EffectiveDate)
		effectiveDate = timestamppb.New(t)
	}

	return &revcom.DSACommission{
		CommID:             uint32(commID),
		DsaCode:            c.DsaCode,
		CommissionType:     fromIntegrationCommissionType(c.CommissionType),
		TierID:             uint32(tierID),
		CommissionAmount:   c.CommissionAmount.String(),
		CommissionCurrency: c.CommissionCurrency.String(),
		UpdatedBy:          c.UpdatedBy,
		CreatedAt:          createdAt,
		UpdatedAt:          updatedAt,
		EffectiveDate:      effectiveDate,
		TrxType:            c.TrxType,
		RemitType:          c.RemitType,
	}
}
