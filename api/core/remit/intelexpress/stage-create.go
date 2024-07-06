package intelexpress

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bojanz/currency"
	"google.golang.org/grpc/codes"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/api/util"
	"brank.as/petnet/serviceutil/logging"
)

type cacheSendRemit struct {
	StageReq core.Remittance `json:"stage_request"`
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
			log.Error()
		}
	}()

	if s.st.OrderIDExists(ctx, r.DsaOrderID) {
		log.Error("order already exist")
		cErr = coreerror.NewCoreError(codes.AlreadyExists, fmt.Sprintf("a transaction with order id: %s, already exists", r.DsaOrderID))
		err = cErr
		return nil, cErr
	}

	zamt, err := currency.NewMinor("0", r.SourceAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid currency")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid currency")
	}
	r.GrossTotal = r.SourceAmount
	r.Tax = zamt
	r.Charge = zamt
	r.Agent.LocationID = core.LocationID // todo: this will be given by petnet

	cache, err := json.Marshal(cacheSendRemit{StageReq: r})
	if err != nil {
		logging.WithError(err, log).Error("caching request marshal error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	res := &core.RemitResponse{
		PrincipalAmount: r.SourceAmount,
		RemitAmount:     r.SourceAmount,
		GrossTotal:      r.SourceAmount,
		Tax:             zamt,
		Charge:          zamt,
	}

	rcReq, err := s.st.CreateRemitCache(ctx, storage.RemitCache{
		DsaID:          r.DsaID,
		UserID:         r.UserID,
		RemcoID:        r.RemitPartner,
		RemType:        core.CreateRemType,
		PtnrRemType:    static.Sendout,
		RemcoMemberID:  r.Remitter.PartnerMemberID,
		RemcoControlNo: r.ControlNo,
		Step:           storage.StageStep,
		Remit:          cache,
	})
	if err != nil {
		logging.WithError(err, log).Error("creating remit db cache error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDatabaseError)
	}

	res.TransactionID = rcReq.TxnID
	r.TransactionID = rcReq.TxnID
	return res, nil
}
