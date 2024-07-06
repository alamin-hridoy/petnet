package perahubremit

import (
	"context"
	"encoding/json"
	"time"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/api/util"
	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (s *Svc) ProcessRemit(ctx context.Context, r core.ProcessRemit) (*core.ProcessRemit, error) {
	log := logging.FromContext(ctx).WithField("txID", r.TransactionID)
	log.Trace("processing")
	rc, err := s.st.GetRemitCache(ctx, r.TransactionID)
	if err != nil {
		if err == storage.ErrNotFound {
			logging.WithError(err, log).Error("remit not found")
			return nil, status.Error(codes.NotFound, "remittance not found")
		}

		logging.WithError(err, log).Error("get remit cache db error for ic")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDatabaseError)
	}
	r.RemitCache = *rc
	log = log.WithField("remit_cache", *rc)
	log.Debug("from cache")

	if rc.Step != storage.StageStep {
		return nil, status.Error(codes.InvalidArgument, "invalid remittance type")
	}
	log.Trace("disbursing")
	return s.disburse(ctx, r)
}

func (s *Svc) disburse(ctx context.Context, r core.ProcessRemit) (rres *core.ProcessRemit, cErr error) {
	log := logging.FromContext(ctx)
	dsaCode := phmw.GetDsaCode(ctx)
	var err error
	c := cacheDisburse{}
	if err = json.Unmarshal(r.RemitCache.Remit, &c); err != nil {
		logging.WithError(err, log).Error("unmarshal remit cache error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}
	creq := c.StageReq
	cres := c.StageResp

	defer func() {
		_, err := util.RecordConfirmTxn(ctx, s.st.(*postgres.Storage), creq, util.ConfirmTxnOpts{
			TxnID:       r.TransactionID,
			TxnType:     storage.DisburseType,
			TxnErr:      err,
			PtnrRemType: static.Payout,
		})
		if err != nil {
			log.Error(err)
		}
	}()

	if _, err = s.ph.PerahubRemitconfirm(ctx, perahub.PerahubRemitConfirmRequest{
		ID:              cres.Result.ID,
		ReferenceNumber: cres.Result.ReferenceNumber,
		PartnerCode:     dsaCode,
	}); err != nil {
		logging.WithError(err, log).Error("confirm remit error")
		return nil, handlePerahubRemitError(err)
	}
	t := time.Now()

	r.Processed = t
	r.ControlNumber = creq.ControlNo
	return &r, nil
}
