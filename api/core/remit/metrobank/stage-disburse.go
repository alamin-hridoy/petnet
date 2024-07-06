package metrobank

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/bojanz/currency"
	"google.golang.org/grpc/codes"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/api/util"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/random"
)

type cacheDisburse struct {
	StageReq  core.Remittance               `json:"stage_request"`
	StageResp perahub.MBInquireResponseBody `json:"stage_response"`
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

	if s.st.OrderIDExists(ctx, rmt.DsaOrderID) {
		log.Error("order already exist")
		cErr = coreerror.NewCoreError(codes.AlreadyExists, fmt.Sprintf("a transaction with order id: %s, already exists", rmt.DsaOrderID))
		err = cErr
		return nil, cErr
	}

	rmt.Agent.LocationID = "191"       // todo: this will be given by petnet
	rmt.Agent.LocationName = "MALOLOS" // todo: this will be given by petnet
	res, err := s.ph.MBInquire(ctx, perahub.MBInquireRequest{
		RefNo:        random.NumberString(18),
		ControlNo:    rmt.ControlNo,
		LocationID:   "0",
		UserID:       json.Number(strconv.Itoa(rmt.Agent.UserID)),
		LocationName: rmt.Agent.LocationName,
	})
	if err != nil {
		logging.WithError(err, log).Error("inquire error for disburse")
		return nil, handleMetroBankError(err)
	}

	if res.Msg != perahub.MBAwaitPayment {
		log.Error("not awaiting disbursement")
		cErr = coreerror.NewCoreError(codes.InvalidArgument, "Remittance is not awaiting disbursement")
		err = cErr
		return nil, cErr
	}

	amt, err := currency.NewAmount(string(res.Result.PrincipalAmount), "PHP")
	if err != nil {
		logging.WithError(err, log).Error("invalid principal amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid principal amount")
	}

	zamt, err := currency.NewMinor("0", "PHP")
	if err != nil {
		logging.WithError(err, log).Error("invalid currency")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid currency")
	}

	rmt.GrossTotal = currency.ToMinor(amt.Round())
	rmt.SourceAmount = currency.ToMinor(amt.Round())
	rmt.DestAmount = zamt
	rmt.Tax = zamt
	rmt.Charge = zamt
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
