package uniteller

import (
	"context"
	"encoding/json"
	"strconv"

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

	locationID := "371"                                 // todo: this will be given by petnet
	locationName := "Information Technology Department" // todo: this will be given by petnet
	agentID := "84424911"                               // todo: this will be given by petnet
	agentCode := "HOA"
	deviceID := r.DeviceID
	if phmw.GetTerminalID(ctx) != "" {
		deviceID = phmw.GetTerminalID(ctx)
	}
	refNo := random.NumberString(18)

	uid, _ := strconv.Atoi(r.PtnrUserID)
	rmt, err := s.ph.UNTInquire(ctx, perahub.UNTInquireRequest{
		Branch:          locationName,
		OutletCode:      locationID,
		LocationID:      "0",
		BranchCode:      locationID,
		LocationCode:    locationID,
		ReferenceNumber: refNo,
		ControlNumber:   r.ControlNo,
		LocationName:    locationName,
		UserID:          uid,
		DeviceId:        deviceID,
		AgentId:         agentID,
		AgentCode:       agentCode,
		// should be set statically
		Currency: "PHP",
	})
	if err != nil {
		logging.WithError(err, log).Error("inquire error for search")
		return nil, handleUnitellerError(err)
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

	delete(objmap, "response_code")
	delete(objmap, "control_number")
	delete(objmap, "principal_amount")
	delete(objmap, "currency")

	b, err = json.Marshal(objmap)
	if err != nil {
		logging.WithError(err, log).Error("marshalling object map error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	amt, err := currency.NewAmount(
		rmt.Result.PrincipalAmount.String(),
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
		Status:         rmt.Result.ResponseCode,
		RemitAmount:    amtMinor,
		DisburseAmount: amtMinor,
		OtherInfo:      b,
	}, nil
}
