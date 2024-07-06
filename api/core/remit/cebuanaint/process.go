package cebuanaint

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	"github.com/bojanz/currency"
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

		logging.WithError(err, log).Error("get remit cache db error for cebInt")
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
	t := time.Now()

	svCh, err := currency.NewAmount(string(cres.Result.ServiceCharge), cres.Result.Currency)
	if err != nil {
		logging.WithError(err, log).Error("invalid Service Charge amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid service charge amount")
	}
	pnpAm, err := currency.NewAmount(string(cres.Result.PrincipalAmount), cres.Result.Currency)
	if err != nil {
		logging.WithError(err, log).Error("invalid principal amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid principal amount")
	}

	ttlAm, err := currency.Amount.Add(svCh, pnpAm)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, fmt.Sprintf("Failed to convert to total amount from service fee: %s and principal amount %s", svCh, pnpAm))
	}

	if _, err = s.ph.CEBINTPayout(ctx, perahub.CEBINTPayoutRequest{
		// inquire response cache
		IsDomestic:        string(cres.Result.IsDomestic),
		ClientReferenceNo: string(cres.Result.ClientReferenceNo),
		ControlNumber:     cres.Result.ControlNumber,
		SenderName:        cres.Result.SenderName,
		ReceiverName:      cres.Result.ReceiverName,
		PrincipalAmount:   cres.Result.PrincipalAmount,
		ServiceCharge:     string(cres.Result.ServiceCharge),
		BirthDate:         cres.Result.BirthDate,
		BeneficiaryID:     string(cres.Result.BeneficiaryID),

		// internal/static
		TrxDate:                  t.Format("2006-01-02"),
		RemcoID:                  "19",
		TrxType:                  "2",
		RiskScore:                "1",
		RiskCriteria:             "1",
		FormType:                 "0",
		FormNumber:               "0",
		PayoutType:               "1",
		IdentificationTypeID:     "11",
		PassportIDIssuedCountry:  "166",
		McRate:                   "0",
		BuyBackAmount:            "0",
		RateCategory:             "0",
		McRateID:                 "0",
		IDIssuedState:            "Null",
		IDIssuingState:           "Null",
		InternationalPartnerCode: "PL0005",

		// userinput
		ContactNumber:      creq.Receiver.Phone.Number,
		LocationID:         "0",
		LocationName:       creq.Agent.LocationName,
		UserID:             json.Number(strconv.Itoa(creq.Agent.UserID)),
		CustomerID:         strconv.Itoa(creq.Agent.UserID),
		CurrencyID:         perahub.CurrencyNumber(creq.DestAmount.CurrencyCode()),
		CustomerName:       cn,
		RemoteLocationID:   json.Number(creq.Agent.LocationID),
		DstAmount:          creq.DestAmount.Amount.Round().Number(),
		TotalAmount:        ttlAm.Number(),
		RemoteUserID:       json.Number(strconv.Itoa(creq.Agent.UserID)),
		RemoteIPAddress:    creq.Agent.IPAddress,
		OriginatingCountry: creq.TransactionDetails.SrcCtry,
		DestinationCountry: creq.TransactionDetails.DestCtry,
		PurposeTransaction: creq.TxnPurpose,
		SourceFund:         creq.Receiver.SourceFunds,
		Occupation:         creq.Receiver.Employment.Occupation,
		RelationTo:         creq.Receiver.ReceiverRelation,
		BirthPlace:         creq.Receiver.BirthPlace,
		BirthCountry:       creq.Receiver.BirthCountry,
		IDType:             creq.Receiver.PrimaryID.IDType,
		IDNumber:           creq.Receiver.PrimaryID.Number,
		Address:            perahub.FormatAddress(creq.Receiver.Address.Address1, creq.Receiver.Address.Address2),
		Barangay:           creq.Receiver.Address.Zone,
		City:               creq.Receiver.Address.City,
		Province:           creq.Receiver.Address.Province,
		ZipCode:            creq.Receiver.Address.PostalCode,
		Country:            creq.Receiver.Address.Country,
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
		ReferenceNumber: random.NumberString(18),
		AgentCode:       creq.Agent.AgentCode,
		IDIssuedBy:      creq.Receiver.PrimaryID.Country,
		IDDateOfIssue: fmt.Sprintf("%s-%s-%s",
			creq.Receiver.PrimaryID.Issued.Year,
			creq.Receiver.PrimaryID.Issued.Month,
			creq.Receiver.PrimaryID.Issued.Day),
		IDIssuingCountry: creq.Receiver.PrimaryID.Country,
		IDIssuedCity:     creq.Receiver.PrimaryID.City,
		DsaCode:          dsaCode,
		DsaTrxType:       dsaTrxType,
	}); err != nil {
		logging.WithError(err, log).Error("cebInt payout error")
		return nil, handleCebuanaIntError(err)
	}

	r.Processed = t
	r.ControlNumber = creq.ControlNo
	return &r, nil
}
