package cebuana

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

		logging.WithError(err, log).Error("get remit cache db error for ceb")
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
	rn := perahub.CombinedName(creq.Receiver.FName, creq.Receiver.MdName, creq.Receiver.LName)
	sn := perahub.CombinedName(creq.Remitter.FName, creq.Remitter.MdName, creq.Remitter.LName)
	ttl, err := creq.SourceAmount.Add(creq.Charge)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Failed to convert to total amount from service fee: %s and principal amount %s", cres.Result.ServiceFee, creq.SourceAmount.Amount.Round().Number()))
	}
	resS, err := s.ph.CebuanaSendMoney(ctx, perahub.CebuanaSendRequest{
		// internal/static
		RemcoID:            "9",
		PayoutType:         "1",
		BbAmount:           "0",
		RiskScore:          "1",
		LocationID:         "0",
		RemoteLocationID:   json.Number(creq.Agent.LocationID),
		McRateID:           "1",
		RateCtg:            "1",
		McRate:             "0",
		FormType:           "OAR",
		FormNumber:         "MAL0101011",
		TrxDate:            t.Format("2006-01-02"),
		TrxType:            "1",
		CurrencyID:         "1",
		SendCurrencyID:     "6",
		LocationName:       "MALOLOS",
		OriginatingCountry: creq.TransactionDetails.SrcCtry,
		DestinationCountry: creq.TransactionDetails.DestCtry,

		// userinput
		BeneficiaryID:      creq.Receiver.RecipientID,
		Barangay:           creq.Remitter.Address.Zone,
		IsDomestic:         strconv.Itoa(creq.TransactionDetails.IsDomestic),
		UserID:             creq.Agent.UserID,
		CustomerID:         strconv.Itoa(creq.Agent.UserID),
		CustomerName:       rn,
		DstAmount:          creq.DestAmount.Amount.Round().Number(),
		TotalAmount:        ttl.Amount.Round().Number(),
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
		City:               creq.Remitter.Address.City,
		Province:           creq.Remitter.Address.Province,
		ZipCode:            creq.Remitter.Address.PostalCode,
		Country:            creq.Remitter.Address.Country,
		ContactNumber:      creq.Remitter.Phone.Number,
		CurrentAddress: perahub.NonexAddress{
			Address1: creq.Remitter.Address.Address1,
			Address2: creq.Remitter.Address.Address2,
			Barangay: creq.Remitter.Address.Zone,
			City:     creq.Remitter.Address.City,
			Province: creq.Remitter.Address.Province,
			ZipCode:  creq.Remitter.Address.PostalCode,
			Country:  creq.Remitter.Address.Country,
		},
		PermanentAddress: perahub.NonexAddress{
			Address1: creq.Remitter.Address.Address1,
			Address2: creq.Remitter.Address.Address2,
			Barangay: creq.Remitter.Address.Zone,
			City:     creq.Remitter.Address.City,
			Province: creq.Remitter.Address.Province,
			ZipCode:  creq.Remitter.Address.PostalCode,
			Country:  creq.Remitter.Address.Country,
		},
		SenderName:      sn,
		ReceiverName:    rn,
		ClientRefNo:     creq.DsaOrderID,
		ServiceCharge:   cres.Result.ServiceFee,
		AgentCode:       creq.Agent.AgentCode,
		PrincipalAmount: creq.SourceAmount.Amount.Round().Number(),
		DsaCode:         phmw.GetDsaCode(ctx),
		DsaTrxType:      util.GetPerahubTrxType(ctx),
	})
	if err != nil {
		logging.WithError(err, log).Error("CEB send error")
		return nil, handleCebuanaError(err)
	}

	r.Processed = t
	creq.ControlNo = resS.Result.ControlNo
	r.ControlNumber = resS.Result.ControlNo
	return &r, nil
}

func (s *Svc) disburse(ctx context.Context, r core.ProcessRemit) (rres *core.ProcessRemit, cErr error) {
	log := logging.FromContext(ctx)
	c := cacheDisburse{}
	dsaCode := phmw.GetDsaCode(ctx)
	dsaTrxType := util.GetPerahubTrxType(ctx)

	var err error
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
	ttl, err := creq.SourceAmount.Add(creq.Charge)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Failed to convert to total amount from service fee: %s and principal amount %s", cres.Result.ServiceCharge, cres.Result.PnplAmt.String()))
	}
	if _, err = s.ph.CEBPayout(ctx, perahub.CEBPayoutRequest{
		// inquire response cache
		ClientRefNo:   cres.Result.ClientReferenceNo,
		ControlNo:     cres.Result.ControlNo,
		SenderName:    cres.Result.SenderName,
		RcvName:       cres.Result.RcvName,
		PnplAmt:       cres.Result.PnplAmt,
		ServiceCharge: cres.Result.ServiceCharge,
		BirthDate:     cres.Result.BirthDate,

		// internal/static
		BuyBackAmount: "0",
		MCRate:        "0",
		RateCat:       "0",
		MCRateID:      "0",
		FormType:      "0",
		FormNumber:    "0",
		TxnDate:       t.Format("2006-01-02"),
		RemcoID:       "9",
		RiskScore:     "1",
		PayoutType:    "1",
		IDTypeID:      "1",
		TxnType:       "2",

		// userinput
		LocationID:    "0",
		RmtLocID:      json.Number(creq.Agent.LocationID),
		LocationName:  creq.Agent.LocationName,
		TotalAmount:   ttl.Amount.Round().Number(),
		DstAmount:     creq.DestAmount.Amount.Round().Number(),
		UserID:        json.Number(strconv.Itoa(creq.Agent.UserID)),
		CustomerID:    strconv.Itoa(creq.Agent.UserID),
		CurrencyID:    json.Number(perahub.CurrencyNumber(creq.DestAmount.CurrencyCode())),
		BeneficiaryID: json.Number(creq.Receiver.RecipientID),
		IsDomestic:    strconv.Itoa(creq.TransactionDetails.IsDomestic),
		CustomerName:  cn,
		RmtUserID:     json.Number(strconv.Itoa(creq.Agent.UserID)),
		RmtIpADD:      creq.Agent.IPAddress,
		OrgnCtry:      creq.TransactionDetails.SrcCtry,
		DestCtry:      creq.TransactionDetails.DestCtry,
		PurposeTxn:    creq.TxnPurpose,
		SourceFund:    creq.Receiver.SourceFunds,
		Occupation:    creq.Receiver.Employment.Occupation,
		RelationTo:    creq.Receiver.ReceiverRelation,
		BirthPlace:    creq.Receiver.BirthPlace,
		BirthCountry:  creq.Receiver.BirthCountry,
		IDType:        creq.Receiver.PrimaryID.IDType,
		IDNumber:      creq.Receiver.PrimaryID.Number,
		Address:       perahub.FormatAddress(creq.Receiver.Address.Address1, creq.Receiver.Address.Address2),
		Barangay:      creq.Receiver.Address.Zone,
		City:          creq.Receiver.Address.City,
		Province:      creq.Receiver.Address.Province,
		ZipCode:       creq.Receiver.Address.PostalCode,
		Country:       creq.Receiver.Address.Country,
		ContactNumber: creq.Receiver.Phone.Number,
		CurrentAddress: perahub.NonexAddress{
			Address1: creq.Receiver.Address.Address1,
			Address2: creq.Receiver.Address.Address2,
			Barangay: creq.Receiver.Address.Zone,
			City:     creq.Receiver.Address.City,
			Province: creq.Receiver.Address.Province,
			ZipCode:  creq.Receiver.Address.PostalCode,
			Country:  creq.Receiver.Address.Country,
		},
		PermanentAddress: perahub.NonexAddress{
			Address1: creq.Receiver.Address.Address1,
			Address2: creq.Receiver.Address.Address2,
			Barangay: creq.Receiver.Address.Zone,
			City:     creq.Receiver.Address.City,
			Province: creq.Receiver.Address.Province,
			ZipCode:  creq.Receiver.Address.PostalCode,
			Country:  creq.Receiver.Address.Country,
		},
		RefNo:      random.NumberString(18),
		DsaCode:    dsaCode,
		DsaTrxType: dsaTrxType,
	}); err != nil {
		logging.WithError(err, log).Error("ceb payout error")
		return nil, handleCebuanaError(err)
	}

	r.Processed = t
	r.ControlNumber = creq.ControlNo
	return &r, nil
}
