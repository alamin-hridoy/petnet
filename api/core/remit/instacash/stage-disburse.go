package instacash

import (
	"context"
	"encoding/json"

	"github.com/bojanz/currency"
	"google.golang.org/grpc/codes"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/api/util"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/random"
)

type cacheDisburse struct {
	StageReq  core.Remittance                      `json:"stage_request"`
	StageResp perahub.InstaCashInquireResponseBody `json:"stage_response"`
}

func (s *Svc) StageDisburseRemit(ctx context.Context, rmt core.Remittance) (rres *core.Remittance, cErr error) {
	log := logging.FromContext(ctx)
	// err variable used in defer to save error in db
	var err error
	defer func() {
		_, err := util.RecordStageTxn(ctx, s.st.(*postgres.Storage), rmt, util.StageTxnOpts{
			TxnID:       rmt.TransactionID,
			TxnType:     storage.DisburseType,
			TxnErr:      err,
			PtnrRemType: static.Payout,
		})
		if err != nil {
			logging.WithError(err, log).Error("record transaction error")
		}
	}()

	rmt.Agent.LocationID = "371"                                 // todo: this will be given by petnet
	rmt.Agent.LocationName = "Information Technology Department" // todo: this will be given by petnet
	rmt.Agent.AgentID = "84424911"                               // todo: this will be given by petnet
	rmt.Agent.AgentCode = "HOA"
	rmt.Agent.DeviceID = "a2c91f7155d18ba34c08782c86b07338b2fd1a6f06e30d1243" // todo: this will be given by petnet
	rmt.GeneratedRefNo = random.NumberString(18)
	rmt.Agent.UserID = 5188
	if phmw.GetTerminalID(ctx) != "" {
		rmt.Agent.DeviceID = phmw.GetTerminalID(ctx)
	}

	res, err := s.ph.InstaCashInquire(ctx, perahub.InstaCashInquireRequest{
		Branch:          rmt.Agent.LocationName,
		OutletCode:      rmt.Agent.LocationID,
		ReferenceNumber: rmt.GeneratedRefNo,
		ControlNumber:   rmt.ControlNo,
		LocationID:      "0",
		UserID:          rmt.Agent.UserID,
		LocationName:    rmt.Agent.LocationName,
		DeviceId:        rmt.Agent.DeviceID,
		AgentId:         rmt.Agent.AgentID,
		AgentCode:       rmt.Agent.AgentCode,
		BranchCode:      rmt.Agent.LocationID,
		LocationCode:    rmt.Agent.LocationID,
	})
	if err != nil {
		logging.WithError(err, log).Error("inquire error for disburse")
		return nil, handleInstaCashError(err)
	}

	if res.Message != perahub.InstaCashAwaitPayment {
		log.Error("not awaiting disbursement")
		cErr = coreerror.NewCoreError(codes.InvalidArgument, "Remittance is not awaiting disbursement")
		err = cErr
		return nil, cErr
	}

	if res.Result.Currency != "PHP" && res.Result.Currency != "USD" {
		log.Error("invalid currency. not php or usd")
		cErr = coreerror.NewCoreError(codes.InvalidArgument, "Invalid principal currency for principal amount, only PHP and USD are accepted")
		err = cErr
		return nil, cErr
	}

	amt, err := currency.NewAmount(string(res.Result.PrincipalAmount), res.Result.Currency)
	if err != nil {
		logging.WithError(err, log).Error("invalid principal amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid principal amount")
	}

	zamt, err := currency.NewMinor("0", res.Result.Currency)
	if err != nil {
		logging.WithError(err, log).Error("invalid currency")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid currency")
	}

	rmt.GrossTotal = currency.ToMinor(amt.Round())
	rmt.SourceAmount = currency.ToMinor(amt.Round())
	rmt.DestAmount = zamt
	rmt.Tax = zamt
	rmt.Charge = zamt
	rmt.Remitter.FName, rmt.Remitter.MdName, rmt.Remitter.LName = perahub.FormatName(res.Result.SenderName)
	cache, err := json.Marshal(cacheDisburse{StageReq: rmt, StageResp: *res})
	if err != nil {
		logging.WithError(err, log).Error("caching request marshal error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	rc, err := s.st.CreateRemitCache(ctx, storage.RemitCache{
		DsaID:          rmt.DsaID,
		UserID:         rmt.UserID,
		RemcoID:        rmt.RemitPartner,
		RemType:        core.DisburseRemType,
		PtnrRemType:    static.Payout,
		RemcoMemberID:  rmt.Receiver.PartnerMemberID,
		RemcoControlNo: rmt.ControlNo,
		Step:           storage.StageStep,
		Remit:          cache,
	})
	if err != nil {
		logging.WithError(err, log).Error("creating remit db cache error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDatabaseError)
	}

	rmt.TransactionID = rc.TxnID
	return &rmt, nil
}
