package fees

import (
	"context"
	"database/sql"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fpb "brank.as/petnet/gunk/dsa/v2/fees"
	"brank.as/petnet/profile/storage"
)

func (s *Svc) UpsertFee(ctx context.Context, req *fpb.UpsertFeeRequest) (*fpb.UpsertFeeResponse, error) {
	if err := validation.ValidateStruct(req.Fee,
		validation.Field(&req.Fee.OrgID, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	fee := req.GetFee()

	feeType := 2
	if req.Fee.Type == fpb.FeeType_TypeFee {
		feeType = 1
	}

	f, err := s.core.UpsertFee(ctx, storage.FeeCommission{
		ID:           fee.ID,
		OrgID:        fee.OrgID,
		Type:         feeType,
		OrgProfileID: "10000000-0000-0000-0000-000000000000", // setting to avoid error, this column should be deleted once refactor of org profile has been done
		StartDate: sql.NullTime{
			Time:  fee.GetSchedule().GetStartDate().AsTime(),
			Valid: fee.GetSchedule().GetStartDate().IsValid(),
		},
		EndDate: sql.NullTime{
			Time:  fee.GetSchedule().GetEndDate().AsTime(),
			Valid: fee.GetSchedule().GetEndDate().IsValid(),
		},
		Deleted: sql.NullTime{
			Time:  fee.GetDeleted().AsTime(),
			Valid: fee.GetDeleted().IsValid(),
		},
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to store fee entry")
	}

	return &fpb.UpsertFeeResponse{ID: f.ID}, nil
}

func (s *Svc) UpsertRate(ctx context.Context, req *fpb.UpsertRateRequest) (*fpb.UpsertRateResponse, error) {
	if err := validation.ValidateStruct(req.Rate,
		validation.Field(&req.Rate.FeeComID, validation.Required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	f, err := s.core.UpsertRate(ctx, storage.Rate{
		ID:              req.Rate.ID,
		FeeCommissionID: req.Rate.GetFeeComID(),
		MinVolume:       req.Rate.GetMinVolume(),
		MaxVolume:       req.Rate.GetMaxVolume(),
		TxnRate:         req.Rate.GetTxnRate(),
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to store rate entry")
	}

	return &fpb.UpsertRateResponse{ID: f.ID}, nil
}
