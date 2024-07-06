package ayannah

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
	"brank.as/petnet/svcutil/random"
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

		logging.WithError(err, log).Error("get remit cache db error for aya")
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
			res, err = s.send(ctx, r)
		case rc.RemType == core.DisburseRemType:
			log.Trace("disbursing")
			res, err = s.disburse(ctx, r)
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

func (s *Svc) send(ctx context.Context, r core.ProcessRemit) (rres *core.ProcessRemit, cErr error) {
	log := logging.FromContext(ctx)
	c := cacheSendRemit{}
	var err error
	if err = json.Unmarshal(r.RemitCache.Remit, &c); err != nil {
		logging.WithError(err, log).Error("unmarshal remit cache error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	creq := c.StageReq
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
	rn := perahub.CombinedName(creq.Receiver.FName, creq.Receiver.MdName, creq.Receiver.LName)
	sn := perahub.CombinedName(creq.Remitter.FName, creq.Remitter.MdName, creq.Remitter.LName)
	resS, err := s.ph.AYANNAHSendMoney(ctx, perahub.AYANNAHSendRequest{
		// internal/static
		CurrencyID:    "1",
		RemcoID:       "22",
		TrxType:       "1",
		ServiceCharge: "0",
		RiskScore:     "1",
		PayoutType:    "1",
		ControlNumber: "0",
		McRate:        "1",
		McRateID:      "1",
		RateCategory:  "required",
		FormType:      "0",
		FormNumber:    "0",
		BuyBackAmount: "0",
		DsaCode:       phmw.GetDsaCode(ctx),
		DsaTrxType:    util.GetPerahubTrxType(ctx),

		// userinput
		LocationID:         "0",
		LocationName:       creq.Agent.LocationName,
		UserID:             json.Number(strconv.Itoa(creq.Agent.UserID)),
		TrxDate:            t.Format("2006-01-02"),
		CustomerID:         strconv.Itoa(creq.Agent.UserID),
		IsDomestic:         strconv.Itoa(creq.TransactionDetails.IsDomestic),
		CustomerName:       sn,
		RemoteLocationID:   json.Number(creq.Agent.LocationID),
		DstAmount:          creq.DestAmount.Amount.Round().Number(),
		TotalAmount:        creq.SourceAmount.Amount.Round().Number(),
		RemoteUserID:       json.Number(strconv.Itoa(creq.Agent.UserID)),
		RemoteIPAddress:    creq.Agent.IPAddress,
		PurposeTransaction: creq.TxnPurpose,
		SourceFund:         creq.Remitter.SourceFunds,
		Occupation:         creq.Remitter.Employment.Occupation,
		RelationTo:         creq.Remitter.ReceiverRelation,
		BirthDate:          fmt.Sprintf("%s-%s-%s", creq.Remitter.BirthDate.Year, creq.Remitter.BirthDate.Month, creq.Remitter.BirthDate.Day),
		BirthPlace:         creq.Remitter.BirthPlace,
		BirthCountry:       creq.Remitter.BirthCountry,
		IDType:             creq.Remitter.PrimaryID.IDType,
		IDNumber:           creq.Remitter.PrimaryID.Number,
		Address:            perahub.FormatAddress(creq.Remitter.Address.Address1, creq.Remitter.Address.Address2),
		Barangay:           creq.Remitter.Address.Zone,
		City:               creq.Remitter.Address.City,
		Province:           creq.Remitter.Address.Province,
		ZipCode:            creq.Remitter.Address.PostalCode,
		Country:            creq.Remitter.Address.Country,
		ContactNumber:      creq.Remitter.Phone.Number,
		SenderName:         sn,
		ReceiverName:       rn,
		PrincipalAmount:    json.Number(creq.SourceAmount.Amount.Round().Number()),
		ClientReferenceNo:  creq.DsaOrderID,
		OriginatingCountry: creq.TransactionDetails.SrcCtry,
		DestinationCountry: creq.TransactionDetails.DestCtry,
		IPAddress:          creq.Agent.IPAddress,
		ReferenceNumber:    random.NumberString(18),
	})
	if err != nil {
		logging.WithError(err, log).Error("AYA send error")
		return nil, handleAyannahError(err)
	}

	r.Processed = t
	r.ControlNumber = resS.Result.ControlNumber
	creq.ControlNo = resS.Result.ControlNumber
	return &r, nil
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
	t := time.Now()
	if _, err = s.ph.AYANNAHPayout(ctx, perahub.AYANNAHPayoutRequest{
		// inquire response cache
		ControlNumber:   cres.Result.ControlNumber,
		PrincipalAmount: json.Number(cres.Result.PrincipalAmount),
		SenderName:      cres.Result.SenderName,
		ReceiverName:    cres.Result.ReceiverName,
		Address:         cres.Result.Address,
		City:            cres.Result.City,
		Country:         cres.Result.Country,
		ReferenceNumber: cres.Result.ReferenceNumber,
		TotalAmount:     cres.Result.PrincipalAmount,

		// userinput
		ZipCode:            creq.Receiver.Address.PostalCode,
		ContactNumber:      creq.Receiver.Phone.Number,
		Province:           creq.Receiver.Address.Province,
		Barangay:           creq.Receiver.Address.Zone,
		LocationID:         "0",
		LocationName:       creq.Agent.LocationName,
		UserID:             json.Number(strconv.Itoa(creq.Agent.UserID)),
		CustomerID:         strconv.Itoa(creq.Agent.UserID),
		CurrencyID:         json.Number(perahub.CurrencyNumber(creq.DestAmount.CurrencyCode())),
		IsDomestic:         json.Number(strconv.Itoa(creq.TransactionDetails.IsDomestic)),
		RemoteUserID:       json.Number(strconv.Itoa(creq.Agent.UserID)),
		RemoteIPAddress:    creq.Agent.IPAddress,
		PurposeTransaction: creq.TxnPurpose,
		SourceFund:         creq.Receiver.SourceFunds,
		Occupation:         creq.Receiver.Employment.Occupation,
		RelationTo:         creq.Receiver.ReceiverRelation,
		BirthDate:          fmt.Sprintf("%s-%s-%s", creq.Receiver.BirthDate.Year, creq.Receiver.BirthDate.Month, creq.Receiver.BirthDate.Day),
		BirthPlace:         creq.Receiver.BirthPlace,
		BirthCountry:       creq.Receiver.BirthCountry,
		IDType:             creq.Receiver.PrimaryID.IDType,
		IDNumber:           creq.Receiver.PrimaryID.Number,
		DstAmount:          creq.DestAmount.Amount.Round().Number(),
		CustomerName:       cn,
		OriginatingCountry: creq.TransactionDetails.SrcCtry,
		DestinationCountry: creq.TransactionDetails.DestCtry,
		ClientReferenceNo:  creq.DsaOrderID,
		IPAddress:          creq.Agent.IPAddress,
		RemoteLocationID:   json.Number(creq.Agent.LocationID),
		DsaCode:            dsaCode,
		DsaTrxType:         dsaTrxType,

		// internal/static
		TrxDate:       t.Format("2006-01-02"),
		BuyBackAmount: "0",
		RemcoID:       "22",
		TrxType:       "2",
		RiskScore:     "1",
		PayoutType:    "1",
		FormType:      "0",
		FormNumber:    "0",
		McRate:        "0",
		McRateID:      "1",
		ServiceCharge: "0",
		RateCategory:  "required",
	}); err != nil {
		logging.WithError(err, log).Error("aya payout error")
		return nil, handleAyannahError(err)
	}

	r.Processed = t
	r.ControlNumber = creq.ControlNo
	return &r, nil
}
