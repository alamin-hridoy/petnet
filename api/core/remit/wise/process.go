package wise

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/api/util"
	"brank.as/petnet/serviceutil/logging"
	"github.com/bojanz/currency"
)

func (s *Svc) ProcessRemit(ctx context.Context, r core.ProcessRemit) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx).WithField("txID", r.TransactionID)
	log.Trace("processing")
	// TODO: Need to implement OTP validation
	rc, err := s.st.GetRemitCache(ctx, r.TransactionID)
	if err != nil {
		if err == storage.ErrNotFound {
			logging.WithError(err, log).Error("remit not found")
			return nil, status.Error(codes.NotFound, "remittance not found")
		}
		logging.WithError(err, log).Error("get remit cache db error for wise")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDatabaseError)
	}
	r.RemitCache = *rc
	log = log.WithField("remit_cache", *rc)
	log.Debug("from cache")
	var res *core.ProcessRemit
	switch rc.Step {
	case storage.StageStep:
		log.Trace("sending")
		switch {
		case rc.RemType == core.CreateRemType:
			res, err = s.sendmoney(ctx, r)
		default:
			return nil, status.Error(codes.InvalidArgument, "invalid remittance type")
		}
	default:
		return nil, status.Error(codes.InvalidArgument, "invalid remittance type")
	}
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (s *Svc) sendmoney(ctx context.Context, r core.ProcessRemit) (rres *core.ProcessRemit, cErr error) {
	log := logging.FromContext(ctx)
	var err error
	c := cacheSendRemit{}
	if err = json.Unmarshal(r.RemitCache.Remit, &c); err != nil {
		logging.WithError(err, log).Error("unmarshal remit cache error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}
	creq := c.StageReq
	cres := c.StageRes

	defer func() {
		_, err := util.RecordConfirmTxn(ctx, s.st.(*postgres.Storage), creq, util.ConfirmTxnOpts{
			TxnID:       r.TransactionID,
			TxnType:     storage.SendType,
			TxnErr:      err,
			PtnrRemType: static.Sendout,
		})
		if err != nil {
			log.Error(err)
		}
	}()

	t := time.Now()
	res, err := s.ph.WISEProceedTransfer(ctx, perahub.WISEProceedTransferReq{
		Email:       creq.Remitter.Email,
		RecipientID: json.Number(creq.Receiver.RecipientID),
		Details: perahub.WISEPCDetails{
			Reference: creq.Message,
		},
	})
	if err != nil {
		logging.WithError(err, log).Error("wise transfer error")
		return nil, handleWiseError(err)
	}

	chg, err := currency.NewAmount(string(cres.UpdatedQuoteSummary.TotalFee), cres.UpdatedQuoteSummary.SourceCurrency)
	if err != nil {
		logging.WithError(err, log).Error("invalid total fee amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid total fee amount")
	}

	creq.Charge = currency.ToMinor(chg.Round())
	creq.Tax = core.MustMinor("0", cres.UpdatedQuoteSummary.SourceCurrency)
	creq.CustomerTxnID = res.CustomerTxnID

	r.Processed = t
	r.ControlNumber = creq.ControlNo
	return &r, nil
}
