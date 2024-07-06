package metrobank

import (
	"context"
	"encoding/json"

	"github.com/bojanz/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/random"
)

func (s *Svc) Search(ctx context.Context, r core.SearchRemit) (*core.SearchRemit, error) {
	log := logging.FromContext(ctx)

	locationID := "0"         // todo: this will be given by petnet
	locationName := "MALOLOS" // todo: this will be given by petnet
	refNo := random.NumberString(18)
	rmt, err := s.ph.MBInquire(ctx, perahub.MBInquireRequest{
		RefNo:        refNo,
		ControlNo:    r.ControlNo,
		LocationID:   json.Number(locationID),
		UserID:       json.Number(r.PtnrUserID),
		LocationName: locationName,
	})
	if err != nil {
		logging.WithError(err, log).Error("inquire error for search")
		return nil, handleMetroBankError(err)
	}

	b, err := json.Marshal(rmt.Result)
	if err != nil {
		logging.WithError(err, log).Error("marshalling inquire response error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	var objmap map[string]interface{}
	if err = json.Unmarshal(b, &objmap); err != nil {
		logging.WithError(err, log).Error("unmarshalling to object map error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	delete(objmap, "control_number")
	delete(objmap, "principal_amount")
	delete(objmap, "currency")
	delete(objmap, "status_text")

	b, err = json.Marshal(objmap)
	if err != nil {
		logging.WithError(err, log).Error("marshalling object map error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	cr := rmt.Result.Currency
	if rmt.Result.Currency == "" {
		cr = "PHP"
	}
	amt, err := currency.NewAmount(
		string(rmt.Result.PrincipalAmount),
		cr,
	)
	if err != nil {
		logging.WithError(err, log).Error("pay amount parsing")
		return nil, status.Error(codes.Internal, "failed to load remittance data")
	}

	amtMinor := currency.ToMinor(amt.Round())
	return &core.SearchRemit{
		DsaID:          r.DsaID,
		UserID:         r.UserID,
		RemitPartner:   r.RemitPartner,
		RemitType:      r.RemitType,
		ControlNo:      r.ControlNo,
		Status:         rmt.Result.StatusText,
		RemitAmount:    amtMinor,
		DisburseAmount: amtMinor,
		OtherInfo:      b,
	}, nil
}
