package metrobank

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
		logging.WithError(err, log).Error("get remit cache db error for mb")
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

	cn := perahub.CombinedName(creq.Receiver.FName, creq.Receiver.MdName, creq.Receiver.LName)
	sn := perahub.CombinedName(creq.Remitter.FName, creq.Remitter.MdName, creq.Remitter.LName)
	t := time.Now()
	if _, err = s.ph.MBPayout(ctx, perahub.MBPayoutRequest{
		// internal/static
		Currency:      "PHP",
		CurrencyId:    "1",
		TrxType:       "2",
		ServiceCharge: "0",
		DstAmount:     "0",
		RiskScore:     "1",
		PayoutType:    "1",
		SenderName:    sn,
		McRate:        "0",
		RateCategory:  "0",
		BuyBackAmount: "0",
		McRateId:      "0",
		FormType:      "0",
		FormNumber:    "0",

		// inquire response cache
		ReceiverName:      cn,
		RemcoId:           string(cres.RemcoID),
		TotalAmount:       creq.SourceAmount.Amount.Round().Number(),
		PrincipalAmount:   cres.Result.PrincipalAmount,
		ReferenceNumber:   cres.Result.RefNo,
		ClientReferenceNo: creq.DsaOrderID,

		// userinput
		LocationId:         "0",
		LocationName:       creq.Agent.LocationName,
		UserId:             json.Number(strconv.Itoa(int(creq.Agent.UserID))),
		TrxDate:            t.Format("2006-01-02"),
		CustomerId:         strconv.Itoa(creq.Agent.UserID),
		IsDomestic:         json.Number(strconv.Itoa(int(creq.TransactionDetails.IsDomestic))),
		CustomerName:       cn,
		RemoteLocationId:   json.Number(creq.Agent.LocationID),
		RemoteUserId:       json.Number(strconv.Itoa(int(creq.Agent.UserID))),
		RemoteIpAddress:    creq.Agent.IPAddress,
		PurposeTransaction: creq.TxnPurpose,
		SourceFund:         creq.Receiver.SourceFunds,
		Occupation:         creq.Receiver.Employment.Occupation,
		RelationTo:         creq.Receiver.ReceiverRelation,
		BirthDate:          fmt.Sprintf("%s-%s-%s", creq.Receiver.BirthDate.Year, creq.Receiver.BirthDate.Month, creq.Receiver.BirthDate.Day),
		BirthPlace:         creq.Receiver.BirthPlace,
		BirthCountry:       creq.Receiver.BirthCountry,
		IdType:             creq.Receiver.PrimaryID.IDType,
		IdNumber:           creq.Receiver.PrimaryID.Number,
		Address:            perahub.FormatAddress(creq.Receiver.Address.Address1, creq.Receiver.Address.Address2),
		Barangay:           creq.Receiver.Address.Zone,
		City:               creq.Receiver.Address.City,
		Province:           creq.Receiver.Address.Province,
		ZipCode:            creq.Receiver.Address.PostalCode,
		Country:            creq.Receiver.Address.Country,
		ContactNumber:      creq.Receiver.Phone.Number,
		ControlNumber:      cres.Result.ControlNo,
		OriginatingCountry: creq.TransactionDetails.SrcCtry,
		DestinationCountry: creq.TransactionDetails.DestCtry,
		DsaCode:            dsaCode,
		DsaTrxType:         dsaTrxType,
	}); err != nil {
		logging.WithError(err, log).Error("mb payout error")
		return nil, handleMetroBankError(err)
	}

	r.Processed = t
	r.ControlNumber = creq.ControlNo
	return &r, nil
}
