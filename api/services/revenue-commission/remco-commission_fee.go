package revenue_commission

import (
	"context"
	"encoding/json"
	"fmt"

	revcom_int "brank.as/petnet/api/integration/revenue-commission"
	"brank.as/petnet/api/util"
	revcom "brank.as/petnet/gunk/drp/v1/revenue-commission"
	"brank.as/petnet/serviceutil/logging"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// CreateRemcoCommissionFee Creates Remco Commission Fee record.
func (s *Svc) CreateRemcoCommissionFee(ctx context.Context, req *revcom.CreateRemcoCommissionFeeRequest) (*revcom.RemcoCommissionFee, error) {
	log := logging.FromContext(ctx)
	createReq, err := createRemcoCommissionFeeValidate(req)
	if err != nil {
		logging.WithError(err, log).Error("failed validating request")
		return nil, util.HandleServiceErr(err)
	}

	res, err := s.commissionFeeStore.CreateRemcoCommissionFee(ctx, createReq)

	return toRemcoCommissionFeeResponse(res, err, log)
}

// UpdateRemcoCommissionFee Updates Remco Commission Fee record.
func (s *Svc) UpdateRemcoCommissionFee(ctx context.Context, req *revcom.UpdateRemcoCommissionFeeRequest) (*revcom.RemcoCommissionFee, error) {
	log := logging.FromContext(ctx)
	updateReq, err := updateRemcoCommissionFeeValidate(req)
	if err != nil {
		logging.WithError(err, log).Error("failed validating request")
		return nil, util.HandleServiceErr(err)
	}

	res, err := s.commissionFeeStore.UpdateRemcoCommissionFee(ctx, updateReq)

	return toRemcoCommissionFeeResponse(res, err, log)
}

// GetRemcoCommissionFeeByID Gets Remco Commission Fee record by ID.
func (s *Svc) GetRemcoCommissionFeeByID(ctx context.Context, req *revcom.GetRemcoCommissionFeeByIDRequest) (*revcom.RemcoCommissionFee, error) {
	log := logging.FromContext(ctx)
	if req.FeeID == 0 {
		log.Error("commission fee id required")
		return nil, util.HandleServiceErr(status.New(codes.InvalidArgument, "commission fee id required").Err())
	}

	res, err := s.commissionFeeStore.GetRemcoCommissionFeeByID(ctx, req.FeeID)

	return toRemcoCommissionFeeResponse(res, err, log)
}

// ListRemcoCommissionFee Gets List Of All Remco Commission Fee records.
func (s *Svc) ListRemcoCommissionFee(ctx context.Context, r *emptypb.Empty) (*revcom.ListRemcoCommissionFeeResponse, error) {
	log := logging.FromContext(ctx)
	res, err := s.commissionFeeStore.ListRemcoCommissionFee(ctx)

	return toListRemcoCommissionFeeResponse(res, err, log)
}

// DeleteRemcoCommissionFee Deletes Remco Commission Fee record by ID.
func (s *Svc) DeleteRemcoCommissionFee(ctx context.Context, req *revcom.DeleteRemcoCommissionFeeRequest) (*revcom.RemcoCommissionFee, error) {
	log := logging.FromContext(ctx)
	if req.FeeID == 0 {
		log.Error("commission fee id required")
		return nil, util.HandleServiceErr(status.New(codes.InvalidArgument, "commission fee id required").Err())
	}

	res, err := s.commissionFeeStore.DeleteRemcoCommissionFee(ctx, req.FeeID)

	return toRemcoCommissionFeeResponse(res, err, log)
}

func createRemcoCommissionFeeValidate(req *revcom.CreateRemcoCommissionFeeRequest) (*revcom_int.SaveRemcoCommissionFeeRequest, error) {
	err := validation.ValidateStruct(req,
		validation.Field(&req.RemcoID, validation.Required),
		validation.Field(&req.MinAmount, validation.Required, is.Float),
		validation.Field(&req.MaxAmount, validation.Required, is.Float),
		validation.Field(&req.ServiceFee, validation.Required, is.Float),
		validation.Field(&req.CommissionAmount, validation.Required, is.Float),
		validation.Field(&req.UpdatedBy, validation.Required),
		validation.Field(&req.CommissionAmountOTC, validation.Required, is.Float),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	trxType := revcom_int.TrxTypeInbound
	if req.TrxType == revcom.TrxType_TrxTypeOutbound {
		trxType = revcom_int.TrxTypeOutbound
	}

	return &revcom_int.SaveRemcoCommissionFeeRequest{
		RemcoID:             json.Number(fmt.Sprintf("%d", req.RemcoID)),
		MinAmount:           json.Number(req.MinAmount),
		MaxAmount:           json.Number(req.MaxAmount),
		ServiceFee:          json.Number(req.ServiceFee),
		CommissionType:      toIntegrationCommissionType(req.GetCommissionType()),
		CommissionAmount:    json.Number(req.CommissionAmount),
		CommissionAmountOTC: json.Number(req.CommissionAmountOTC),
		TrxType:             trxType,
		UpdatedBy:           req.UpdatedBy,
	}, nil
}

func updateRemcoCommissionFeeValidate(req *revcom.UpdateRemcoCommissionFeeRequest) (*revcom_int.SaveRemcoCommissionFeeRequest, error) {
	err := validation.ValidateStruct(req,
		validation.Field(&req.FeeID, validation.Required),
		validation.Field(&req.RemcoID, validation.Required),
		validation.Field(&req.MinAmount, validation.Required, is.Float),
		validation.Field(&req.MaxAmount, validation.Required, is.Float),
		validation.Field(&req.ServiceFee, validation.Required, is.Float),
		validation.Field(&req.CommissionAmount, validation.Required, is.Float),
		validation.Field(&req.UpdatedBy, validation.Required),
		validation.Field(&req.CommissionAmountOTC, validation.Required, is.Float),
	)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	trxType := revcom_int.TrxTypeInbound
	if req.TrxType == revcom.TrxType_TrxTypeOutbound {
		trxType = revcom_int.TrxTypeOutbound
	}

	return &revcom_int.SaveRemcoCommissionFeeRequest{
		FeeID:               req.FeeID,
		RemcoID:             json.Number(fmt.Sprintf("%d", req.RemcoID)),
		MinAmount:           json.Number(req.MinAmount),
		MaxAmount:           json.Number(req.MaxAmount),
		ServiceFee:          json.Number(req.ServiceFee),
		CommissionType:      toIntegrationCommissionType(req.GetCommissionType()),
		CommissionAmount:    json.Number(req.CommissionAmount),
		CommissionAmountOTC: json.Number(req.CommissionAmountOTC),
		TrxType:             trxType,
		UpdatedBy:           req.UpdatedBy,
	}, nil
}

func toRemcoCommissionFeeResponse(d *revcom_int.RemcoCommissionFee, err error, log *logrus.Entry) (*revcom.RemcoCommissionFee, error) {
	if err != nil {
		logging.WithError(err, log).Error("commission fee store error")
		return nil, util.HandleServiceErr(err)
	}

	if d == nil {
		return nil, util.HandleServiceErr(status.New(codes.Unknown, "empty commission fee response").Err())
	}

	var (
		feeID, remcoID       int64
		createdAt, updatedAt *timestamppb.Timestamp
		trxType              = revcom.TrxType_TrxTypeInbound
	)

	feeID, err = d.FeeID.Int64()
	if err != nil {
		logging.WithError(err, log).Error("commission fee id is not numeric")
		return nil, util.HandleServiceErr(err)
	}

	remcoID, _ = d.RemcoID.Int64()
	if d.CreatedAt != nil {
		createdAt = timestamppb.New(*d.CreatedAt)
	}

	if d.UpdatedAt != nil {
		updatedAt = timestamppb.New(*d.UpdatedAt)
	}

	if d.TrxType == revcom_int.TrxTypeOutbound {
		trxType = revcom.TrxType_TrxTypeOutbound
	}

	return &revcom.RemcoCommissionFee{
		FeeID:               uint32(feeID),
		RemcoID:             uint32(remcoID),
		MinAmount:           d.MinAmount.String(),
		MaxAmount:           d.MaxAmount.String(),
		ServiceFee:          d.ServiceFee.String(),
		CommissionAmount:    d.CommissionAmount.String(),
		CommissionAmountOTC: d.CommissionAmountOTC.String(),
		CommissionType:      fromIntegrationCommissionType(d.CommissionType),
		TrxType:             trxType,
		UpdatedBy:           d.UpdatedBy,
		CreatedAt:           createdAt,
		UpdatedAt:           updatedAt,
	}, nil
}

func toListRemcoCommissionFeeResponse(dl []revcom_int.RemcoCommissionFee, err error, log *logrus.Entry) (*revcom.ListRemcoCommissionFeeResponse, error) {
	if err != nil {
		logging.WithError(err, log).Error("remco commission fee store error")
		return nil, util.HandleServiceErr(err)
	}

	if len(dl) == 0 {
		return &revcom.ListRemcoCommissionFeeResponse{RemcoCommissionFeeList: []*revcom.RemcoCommissionFee{}}, nil
	}

	feeList := make([]*revcom.RemcoCommissionFee, 0, len(dl))
	for _, d := range dl {
		dsa, _ := toRemcoCommissionFeeResponse(&d, nil, log)
		feeList = append(feeList, dsa)
	}

	return &revcom.ListRemcoCommissionFeeResponse{
		RemcoCommissionFeeList: feeList,
	}, nil
}

func toIntegrationCommissionType(commType revcom.CommissionType) revcom_int.CommissionType {
	switch commType {
	case revcom.CommissionType_CommissionTypeRange:
		return revcom_int.CommissionTypeRange
	case revcom.CommissionType_CommissionTypePercent:
		return revcom_int.CommissionTypePercent
	default:
		return revcom_int.CommissionTypeAbsolute
	}
}

func fromIntegrationCommissionType(commType revcom_int.CommissionType) revcom.CommissionType {
	switch commType {
	case revcom_int.CommissionTypeRange:
		return revcom.CommissionType_CommissionTypeRange
	case revcom_int.CommissionTypePercent:
		return revcom.CommissionType_CommissionTypePercent
	default:
		return revcom.CommissionType_CommissionTypeAbsolute
	}
}
