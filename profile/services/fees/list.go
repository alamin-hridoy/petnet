package fees

import (
	"context"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	fpb "brank.as/petnet/gunk/dsa/v2/fees"
	"brank.as/petnet/profile/storage"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

func (s *Svc) ListFees(ctx context.Context, req *fpb.ListFeesRequest) (*fpb.ListFeesResponse, error) {
	if err := validation.ValidateStruct(req,
		validation.Field(&req.OrgID, validation.Required, is.UUIDv4),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	reqType := int(fpb.FeeType_UnknownFeeType)
	if req.GetType() == fpb.FeeType_TypeCommission.String() {
		reqType = int(fpb.FeeType_TypeCommission)
	} else if req.GetType() == fpb.FeeType_TypeFee.String() {
		reqType = int(fpb.FeeType_TypeFee)
	}

	fs, err := s.core.ListFees(ctx, req.OrgID, storage.LimitOffsetFilter{
		Limit:  req.GetLimit(),
		Offset: req.GetOffset(),
		Type:   reqType,
	})
	if err != nil {
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to store fee entry")
	}

	fcl := []*fpb.Fee{}
	for _, fee := range fs {
		lr, err := s.core.ListRates(ctx, fee.ID)
		if err != nil {
			return nil, err
		}

		rl := []*fpb.Rate{}
		for _, rate := range lr {
			nrl := &fpb.Rate{
				ID:        rate.ID,
				FeeComID:  rate.FeeCommissionID,
				MinVolume: rate.MinVolume,
				MaxVolume: rate.MaxVolume,
				TxnRate:   rate.TxnRate,
			}
			rl = append(rl, nrl)
		}
		pf := &fpb.Fee{
			ID:    fee.ID,
			Type:  fpb.FeeType(fee.Type),
			OrgID: fee.OrgID,
			Rates: rl,
			Schedule: &fpb.Schedule{
				Status: fpb.FeeStatus(fee.FeeStatus),
			},
		}
		if fee.StartDate.Valid {
			pf.Schedule.StartDate = tspb.New(fee.StartDate.Time)
		}
		if fee.EndDate.Valid {
			pf.Schedule.EndDate = tspb.New(fee.EndDate.Time)
		}
		fcl = append(fcl, pf)
	}

	var tot int32
	if len(fcl) > 0 {
		tot = int32(fs[0].Count)
	}
	return &fpb.ListFeesResponse{Fees: fcl, Total: tot}, nil
}
