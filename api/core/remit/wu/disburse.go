package wu

import (
	"context"
	"encoding/json"
	"strconv"
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
	"brank.as/petnet/svcutil/random"
)

type cacheDisburse struct {
	StageReq  core.Remittance              `json:"stage_request"`
	StageResp perahub.RMSearchResponseBody `json:"stage_response"`
}

func (s *Svc) disburse(ctx context.Context, r core.ProcessRemit) (res *core.ProcessRemit, cErr error) {
	log := logging.FromContext(ctx)

	var err error
	c := cacheDisburse{}
	if err = json.Unmarshal(r.RemitCache.Remit, &c); err != nil {
		logging.WithError(err, log).Error("disburse unmarshal remit cache error for wu")
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

	creq.RemcoAltControlNo = cres.Txn.NewMtcn
	creq.Remitter.FName = cres.Txn.Sender.Name.FirstName
	creq.Remitter.LName = cres.Txn.Sender.Name.LastName
	creq.Remitter.Address.Address1 = cres.Txn.Sender.Address.Street
	creq.Remitter.Address.Address2 = ""
	creq.Remitter.Address.City = cres.Txn.Sender.Address.City
	creq.Remitter.Address.State = cres.Txn.Sender.Address.State
	creq.Remitter.Address.PostalCode = cres.Txn.Sender.Address.PostalCode
	creq.Remitter.Address.Country = cres.Txn.Sender.Address.CountryCode.IsoCode.Country
	creq.Remitter.Phone.Number = cres.Txn.Sender.Phone
	creq.Remitter.Mobile.CtyCode = cres.Txn.Sender.MobileDetails.CountryCode
	creq.Remitter.Mobile.Number = cres.Txn.Sender.MobileDetails.Number
	creq.ControlNo = cres.Txn.Mtcn

	var wuID, wuPt string
	if r.RemitCache.RemcoMemberID != "" {
		k, kErr := s.ph.KYCLookup(ctx, perahub.KYCLookupRequest{
			RefNo:      random.InvitationCode(20),
			SearchType: "by_wu_card",
			TrxType:    "PAY",
			MyWUNumber: r.RemitCache.RemcoMemberID,
			FirstName:  creq.Receiver.FName,
			LastName:   creq.Receiver.LName,
			// todo(robin): making this static for now, live this should be provided by user
			// OperatorID:       phmw.GetOperatorID(ctx),
			TerminalID: getTerminalID(ctx),
			OperatorID: operatorID,
			UserCode:   strconv.Itoa(creq.Agent.UserID),
		})
		if kErr != nil {
			logging.WithError(kErr, log).Error("kyc lookup")
			k = &perahub.KYCLookupBody{}
		}
		wuID = k.Customer.MyWUDetails.MyWUNumber
		wuPt = k.Customer.MyWUDetails.CurrentYrPts.String()
	}

	iso := s.iso(ctx)
	gnd := ""
	switch creq.Receiver.Gender {
	case "Male":
		gnd = "M"
	case "Female":
		gnd = "F"
	default:
		err = status.Error(codes.InvalidArgument, "gender is required.")
		return nil, err
	}
	rmt, err := s.ph.RecMoneyPay(ctx, perahub.RecMoneyPayRequest{
		FrgnRefNo: creq.DsaOrderID,
		// TODO(Chad): Figure out lookup
		UserCode: strconv.Itoa(creq.Agent.UserID),
		// CustomerCode:       k.Customer,
		ReceiverFirstName:  creq.Receiver.FName,
		ReceiverMiddleName: creq.Receiver.MdName,
		ReceiverLastName:   creq.Receiver.LName,

		AddrLine1:  creq.Receiver.Address.Address1,
		AddrLine2:  creq.Receiver.Address.Address2,
		PostalCode: creq.Receiver.Address.PostalCode,
		Country:    iso(creq.Receiver.Address.Country),
		CurrCity:   creq.Receiver.Address.City,
		CurrState:  creq.Receiver.Address.State,

		ReceiverAddress1:     creq.Receiver.Address.Address1,
		ReceiverCity:         creq.Receiver.Address.City,
		ReceiverState:        creq.Receiver.Address.State,
		ReceiverStateZip:     creq.Receiver.Address.PostalCode,
		ReceiverCountryCode:  creq.Receiver.Address.Country,
		ReceiverCurrencyCode: creq.DestAmount.CurrencyCode(),
		ReceiverHasPhone: func() string {
			if creq.Receiver.Phone.Number != "" {
				return "Y"
			}
			return ""
		}(),
		RecMobCountryCode: creq.Receiver.Mobile.CtyCode,
		PhoneNumber:       creq.Receiver.Mobile.Number,
		PhoneCityCode:     creq.Receiver.Phone.CtyCode,
		ContactPhone:      creq.Receiver.Phone.Number,
		Email:             creq.Receiver.Email,
		Gender:            gnd,
		Birthdate:         creq.Receiver.BirthDate.String(),
		BirthCountry:      iso(creq.Receiver.BirthCountry),
		Nationality:       iso(creq.Receiver.Nationality),

		IDType:      creq.Receiver.PrimaryID.IDType,
		IDCountry:   iso(creq.Receiver.PrimaryID.Country),
		IDNumber:    creq.Receiver.PrimaryID.Number,
		IDIssueDate: creq.Receiver.PrimaryID.Issued.String(),
		IDHasExpiry: func() string {
			if creq.Receiver.PrimaryID.Expiry == nil {
				return "N"
			}
			return "Y"
		}(),
		IDExpirationDate: creq.Receiver.PrimaryID.Expiry.String(),

		FundSource: creq.Receiver.SourceFunds,
		Occupation: creq.Receiver.Employment.Occupation,
		// TODO(Chad): additional field
		EmployerName:         creq.Receiver.Employment.Employer,
		PositionLevel:        creq.Receiver.Employment.PositionLevel,
		TxnPurpose:           creq.TxnPurpose,
		ReceiverRelationship: creq.Receiver.ReceiverRelation,

		SenderFirstName:    cres.Txn.Sender.Name.FirstName,
		SenderLastName:     cres.Txn.Sender.Name.LastName,
		PdState:            cres.Txn.Payment.ExpectedPayoutLocation.State,
		PdCity:             cres.Txn.Payment.ExpectedPayoutLocation.City,
		PdDestCountry:      cres.Txn.Payment.DestCountry.IsoCode.Country,
		PdDestCurrency:     cres.Txn.Payment.DestCountry.IsoCode.Currency,
		PdOriginatingCity:  cres.Txn.Payment.OriginatingCity,
		PdOrigCountryCode:  cres.Txn.Payment.OrigCountry.IsoCode.Country,
		PdOrigCurrencyCode: cres.Txn.Payment.OrigCountry.IsoCode.Currency,
		PdTransactionType:  cres.Txn.Payment.TransactionType,
		PdExchangeRate:     cres.Txn.Payment.ExchangeRate,
		PdOrigDestCountry:  cres.Txn.Payment.DestCountry.IsoCode.Country,
		PdOrigDestCurrency: cres.Txn.Payment.DestCountry.IsoCode.Currency,

		GrossTotal:       cres.Txn.Financials.GrossTotal,
		PayAmount:        cres.Txn.Financials.PayAmount,
		Principal:        cres.Txn.Financials.Principal,
		Charges:          cres.Txn.Financials.Charges,
		Tolls:            cres.Txn.Financials.Tolls,
		RealPrincipal:    cres.GrossPayout.String(),
		DstAmount:        cres.DST.String(),
		RealNet:          cres.NetPayout.String(),
		FilingTime:       cres.Txn.FilingTime,
		FilingDate:       cres.Txn.FilingDate,
		MoneyTransferKey: cres.Txn.MoneyTransferKey,

		Mtcn:    cres.Txn.Mtcn,
		NewMtcn: cres.Txn.NewMtcn,
		// TODO(Chad): missing data from search
		Message: nil,
		// todo(robin): making this static for now, live this should be provided by user
		// TerminalID:       phmw.GetTerminalID(ctx),
		// OperatorID:       phmw.GetOperatorID(ctx),
		TerminalID:       getTerminalID(ctx),
		OperatorID:       operatorID,
		RemoteTerminalID: remoteOperatorID,
		RemoteOperatorID: remoteOperatorID,

		MyWUNumber: wuID,
		MyWUPoints: wuPt,
		HasLoyalty: func() string {
			if wuID != "" {
				return "Y"
			}
			return ""
		}(),
	})
	if err != nil {
		logging.WithError(err, log).Error("wu rec money pay error")
		return nil, handleWUError(err)
	}

	proc, _ := time.Parse("01-02-06 15:04", rmt.PaidDateTime)
	if proc.IsZero() {
		proc = time.Now()
	}

	r.Processed = proc
	r.ControlNumber = cres.Txn.Mtcn
	return &r, nil
}
