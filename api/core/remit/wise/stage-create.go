package wise

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
	StageReq core.Remittance                 `json:"stage_request"`
	StageRes perahub.WISEPrepareTransferResp `json:"stage_response"`
}

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
			log.Error()
		}
	}()

	if s.st.OrderIDExists(ctx, r.DsaOrderID) {
		log.Error("order already exist")
		cErr = coreerror.NewCoreError(codes.AlreadyExists, fmt.Sprintf("a transaction with order id: %s, already exists", r.DsaOrderID))
		err = cErr
		return nil, cErr
	}

	log.WithField("receiver", r.Receiver).Debug("prepare transfer")
	resp, err := s.ph.WISEPrepareTransfer(ctx, perahub.WISEPrepareTransferReq{
		Email:               r.Remitter.Email,
		RecipientID:         r.Receiver.RecipientID,
		AccountHolderName:   r.Receiver.AccountHolderName,
		SourceAccountNumber: r.Receiver.SourceAccountNumber,
	})
	if err != nil {
		logging.WithError(err, log).Error("inquire error for create")
		return nil, handleWiseError(err)
	}

	zamt, err := currency.NewMinor("0", resp.UpdatedQuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("invalid currency")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid currency")
	}

	pcpl, err := currency.NewAmount(string(resp.UpdatedQuoteSummary.SourceAmount), resp.UpdatedQuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("invalid SourceAmount amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid SourceAmount")
	}

	chg, err := currency.NewAmount(string(resp.UpdatedQuoteSummary.TotalFee), resp.UpdatedQuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("invalid TotalFee amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid TotalFee amount")
	}

	ramt, err := currency.NewAmount(string(resp.UpdatedQuoteSummary.TransferAmount), resp.UpdatedQuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("invalid TotalFee amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid TransferAmount")
	}

	tamt, err := currency.NewAmount(string(resp.UpdatedQuoteSummary.TargetAmount), resp.UpdatedQuoteSummary.TargetCurrency)
	if err != nil {
		logging.WithError(err, log).Error("invalid TotalFee amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid TargetAmount")
	}
	r.DestAmount = currency.ToMinor(tamt.Round())

	r.GrossTotal = currency.ToMinor(pcpl.Round())
	r.SourceAmount = currency.ToMinor(pcpl.Round())
	r.Tax = zamt
	r.Charge = currency.ToMinor(chg.Round())
	r.Agent.LocationID = core.LocationID // todo: this will be given by petnet

	res := &core.RemitResponse{
		PrincipalAmount: currency.ToMinor(pcpl.Round()),
		RemitAmount:     currency.ToMinor(ramt.Round()),
		Charge:          currency.ToMinor(chg.Round()),
		Tax:             zamt,
		GrossTotal:      currency.ToMinor(pcpl.Round()),
	}

	cache, err := json.Marshal(cacheSendRemit{StageReq: r, StageRes: *resp})
	if err != nil {
		logging.WithError(err, log).Error("caching request marshal error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	rcReq, err := s.st.CreateRemitCache(ctx, storage.RemitCache{
		DsaID:         r.DsaID,
		RemcoID:       r.RemitPartner,
		RemType:       core.CreateRemType,
		PtnrRemType:   resp.UpdatedQuoteSummary.PayOut,
		RemcoMemberID: r.Receiver.RecipientID,
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
