package ayannah

import (
	"context"
	"encoding/json"

	"github.com/bojanz/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/random"
)

func (s *Svc) Search(ctx context.Context, r core.SearchRemit) (*core.SearchRemit, error) {
	log := logging.FromContext(ctx)

	locationID := "0"         // todo: this will be given by petnet
	locationName := "locname" // todo: this will be given by petnet
	agentCode := "MAL"        // todo: this will be given by petnet
	deviceID := r.DeviceID
	if phmw.GetTerminalID(ctx) != "" {
		deviceID = phmw.GetTerminalID(ctx)
	}
	rmt, err := s.ph.AYANNAHInquire(ctx, perahub.AYANNAHInquireRequest{
		Branch:          locationName,
		OutletCode:      agentCode,
		ReferenceNumber: random.NumberString(18),
		ControlNumber:   r.ControlNo,
		LocationID:      json.Number(locationID),
		UserID:          json.Number(r.PtnrUserID),
		LocationName:    locationName,
		DeviceID:        deviceID,
		AgentCode:       agentCode,
		LocationCode:    agentCode,
		// static
		Currency: "PHP",
	})
	if err != nil {
		logging.WithError(err, log).Error("inquire error for search")
		return nil, handleAyannahError(err)
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
	delete(objmap, "response_code")

	b, err = json.Marshal(objmap)
	if err != nil {
		logging.WithError(err, log).Error("marshalling object map error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	amt, err := currency.NewAmount(
		rmt.Result.PrincipalAmount,
		rmt.Result.Currency,
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
		RemitAmount:    amtMinor,
		DisburseAmount: amtMinor,
		Status:         rmt.Result.ResponseCode,
		OtherInfo:      b,
	}, nil
}
