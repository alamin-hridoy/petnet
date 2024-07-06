package ussc

import (
	"context"
	"encoding/json"
	"fmt"

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
)

type cacheSendRemit struct {
	StageReq core.Remittance                `json:"stage_request"`
	StageRes perahub.USSCFeeInquiryRespBody `json:"stage_response"`
}

// Source Amounts must be in PHP.
func (s *Svc) StageCreateRemit(ctx context.Context, r core.Remittance) (rres *core.RemitResponse, cErr error) {
	log := logging.FromContext(ctx)
	var err error
	defer func() {
		_, err := util.RecordStageTxn(ctx, s.st.(*postgres.Storage), r, util.StageTxnOpts{
			TxnID:       r.TransactionID,
			TxnType:     storage.SendType,
			TxnErr:      err,
			PtnrRemType: static.Sendout,
		})
		if err != nil {
			logging.WithError(err, log).Error("record transaction error")
		}
	}()

	if s.st.OrderIDExists(ctx, r.DsaOrderID) {
		log.Error("order already exist")
		cErr = coreerror.NewCoreError(codes.AlreadyExists, fmt.Sprintf("a transaction with order id: %s, already exists", r.DsaOrderID))
		err = cErr
		return nil, cErr
	}

	BranchCode := "branch1" // todo: will change, gotten from petnet
	fee, err := s.ph.USSCFeeInquiry(ctx, perahub.USSCFeeInquiryRequest{
		Amount:     r.SourceAmount.Amount.Round().Number(),
		BranchCode: BranchCode, // todo: will change, gotten from petnet

		// static
		Panalokard: "",
		USSCPromo:  "",
	})
	if err != nil {
		logging.WithError(err, log).Error("inquire error for create")
		return nil, handleUSSCError(err)
	}

	if fee.Message != perahub.USSCFeeInquiryOK {
		log.Error("not awaiting disbursement")
		cErr = coreerror.NewCoreError(codes.InvalidArgument, "Remittance is not awaiting disbursement")
		err = cErr
		return nil, cErr
	}

	zamt, err := currency.NewMinor("0", r.SourceAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid currency")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid currency")
	}

	r.Agent.OutletCode = BranchCode
	pcpl, err := currency.NewAmount(fee.Result.PnplAmount, r.SourceAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid Principle amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid Principle amount")
	}

	chg, err := currency.NewAmount(fee.Result.ServiceCharge, r.SourceAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid ServiceCharge amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid ServiceCharge amount")
	}

	ttl, err := pcpl.Add(chg)
	if err != nil {
		logging.WithError(err, log).Error("add charge to principle amount error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	r.Tax = zamt
	r.Charge = currency.ToMinor(chg.Round())
	r.GrossTotal = currency.ToMinor(ttl.Round())
	r.Agent.LocationID = core.LocationID // todo: this will be given by petnet

	res := &core.RemitResponse{
		PrincipalAmount: currency.ToMinor(pcpl.Round()),
		Charge:          currency.ToMinor(chg.Round()),
		GrossTotal:      currency.ToMinor(ttl.Round()),
	}
	cache, err := json.Marshal(cacheSendRemit{StageReq: r, StageRes: *fee})
	if err != nil {
		logging.WithError(err, log).Error("caching request marshal error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	rcReq, err := s.st.CreateRemitCache(ctx, storage.RemitCache{
		DsaID:         r.DsaID,
		UserID:        r.UserID,
		RemcoID:       r.RemitPartner,
		RemType:       core.CreateRemType,
		PtnrRemType:   static.USSCCode,
		RemcoMemberID: r.Receiver.PartnerMemberID,
		Step:          storage.StageStep,
		Remit:         cache,
	})
	if err != nil {
		logging.WithError(err, log).Error("creating remit db cache error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDatabaseError)
	}

	res.TransactionID = rcReq.TxnID
	r.TransactionID = rcReq.TxnID
	return res, nil
}
