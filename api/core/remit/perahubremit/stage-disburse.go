package perahubremit

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
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/api/util"
	"brank.as/petnet/serviceutil/logging"
)

type cacheDisburse struct {
	StageReq  core.Remittance                          `json:"stage_request"`
	StageResp perahub.PerahubRemitValidateResponseBody `json:"stage_response"`
}

func (s *Svc) StageDisburseRemit(ctx context.Context, rmt core.Remittance) (rres *core.Remittance, cErr error) {
	log := logging.FromContext(ctx)
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
	res, err := s.ph.PerahubRemitInquire(ctx, perahub.PerahubRemitInquireRequest{
		ControlNumber: rmt.ControlNo,
		LocationID:    371, // todo provided by perahub
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub remit inquire error")
		return nil, handlePerahubRemitError(err)
	}

	if res != nil {
		rmt.SourceAmount, _ = currency.NewMinor(strconv.Itoa(res.Result.PrincipalAmount), res.Result.IsoCurrency)
		rmt.Remitter.FName = res.Result.SenderFirstName
		rmt.Remitter.LName = res.Result.SenderLastName
		rmt.Remitter.MdName = res.Result.SenderMiddleName
		rmt.Receiver.FName = res.Result.ReceiverFirstName
		rmt.Receiver.LName = res.Result.ReceiverLastName
		rmt.Receiver.MdName = res.Result.ReceiverMiddleName
	}

	val, err := s.ph.PerahubRemitValidate(ctx, perahub.PerahubRemitValidateRequest{
		ControlNumber: rmt.ControlNo,
		PrincipalAmount: func() int {
			amt, err := rmt.SourceAmount.Int64()
			if err != nil {
				logging.WithError(err, log).Error("invalid source amount")
				return 0
			}
			return int(amt)
		}(),
		SenderFirstName:    rmt.Remitter.FName,
		SenderLastName:     rmt.Remitter.LName,
		SenderMiddleName:   rmt.Remitter.MdName,
		ReceiverFirstName:  rmt.Remitter.FName,
		ReceiverLastName:   rmt.Remitter.LName,
		ReceiverMiddleName: rmt.Remitter.MdName,
		PartnerCode:        phmw.GetDsaCode(ctx),
	})
	if err != nil {
		logging.WithError(err, log).Error("perahub remit validate error")
		return nil, handlePerahubRemitError(err)
	}

	if val != nil {
		rmt.GeneratedRefNo = val.Result.ReferenceNumber
		rmt.CustomerTxnID = strconv.Itoa(val.Result.ID)
	}

	amt, err := currency.NewAmount(strconv.Itoa(res.Result.PrincipalAmount), res.Result.IsoCurrency)
	if err != nil {
		logging.WithError(err, log).Error("invalid principal amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid principal amount")
	}

	rmt.SourceAmount = currency.ToMinor(amt.Round())

	cache, err := json.Marshal(cacheDisburse{StageReq: rmt, StageResp: *val})
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
