package iremit

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

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
		logging.WithError(err, log).Error("get remit cache db error for ir")
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
	dsaTrxType := util.GetPerahubTrxType(ctx)
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
	creq.Remitter.FName, creq.Remitter.MdName, creq.Remitter.LName = perahub.FormatName(cres.Result.SenderName)

	cn := perahub.CombinedName(creq.Receiver.FName, creq.Receiver.MdName, creq.Receiver.LName)
	t := time.Now()
	if _, err = s.ph.IRemitPayout(ctx, perahub.IRPayoutRequest{
		// internal/static
		RiskScore:     "1", // from reymar: this is currently internal to Perahub, we can leave this as blank for the meantime
		RiskCriteria:  "1",
		TxnType:       "2",
		ServiceCharge: "0",
		DstAmount:     "0",
		TxnDate:       t.Format("2006-01-02"),
		MCRate:        "0",
		BBAmt:         "0",
		RateCat:       "0",
		MCRateID:      "0",
		FormType:      "0",
		FormNumber:    "0",
		PayoutType:    "1", // todo
		RmtLocID:      json.Number(creq.Agent.LocationID),
		LocationID:    "0",
		LocationName:  creq.Agent.LocationName,

		// inquire response cache
		RemcoID:     json.Number(cres.RemcoID),
		RefNo:       cres.Result.RefNo,
		SenderName:  cres.Result.SenderName,
		RcvName:     cn,
		TotalAmount: creq.SourceAmount.Amount.Round().Number(),
		PnplAmt:     json.Number(creq.SourceAmount.Amount.Round().Number()),

		// userinput
		ControlNo:     creq.ControlNo,
		ClientRefNo:   creq.DsaOrderID,
		CustomerName:  cn,
		ContactNumber: creq.Receiver.Phone.Number,
		CustomerID:    json.Number(strconv.Itoa(creq.Agent.UserID)),
		PurposeTxn:    creq.TxnPurpose,
		SourceFund:    creq.Receiver.SourceFunds,
		Occupation:    creq.Receiver.Employment.Occupation,
		RelationTo:    creq.Receiver.ReceiverRelation,
		BirthDate:     fmt.Sprintf("%s-%s-%s", creq.Receiver.BirthDate.Year, creq.Receiver.BirthDate.Month, creq.Receiver.BirthDate.Day),
		BirthPlace:    creq.Receiver.BirthPlace,
		BirthCountry:  creq.Receiver.BirthCountry,
		IDType:        creq.Receiver.PrimaryID.IDType,
		IDNumber:      creq.Receiver.PrimaryID.Number,
		Address:       perahub.FormatAddress(creq.Receiver.Address.Address1, creq.Receiver.Address.Address2),
		City:          creq.Receiver.Address.City,
		Province:      creq.Receiver.Address.Province,
		ZipCode:       creq.Receiver.Address.PostalCode,
		Country:       creq.Receiver.Address.Country,
		Barangay:      creq.Receiver.Address.Zone,
		CurrencyID:    json.Number(json.Number(perahub.CurrencyNumber(creq.DestAmount.CurrencyCode()))),
		UserID:        json.Number(strconv.Itoa(creq.Agent.UserID)),
		IsDomestic:    json.Number(strconv.Itoa(creq.TransactionDetails.IsDomestic)),
		RmtUserID:     json.Number(strconv.Itoa(creq.Agent.UserID)),
		RmtIPAddr:     creq.Agent.IPAddress,
		OrgnCtry:      creq.TransactionDetails.SrcCtry,
		DestCtry:      creq.TransactionDetails.DestCtry,
		DsaCode:       dsaCode,
		DsaTrxType:    dsaTrxType,
	}); err != nil {
		logging.WithError(err, log).Error("ir payout error")
		return nil, handleIRemitError(err)
	}

	r.Processed = t
	r.ControlNumber = creq.ControlNo
	return &r, nil
}
