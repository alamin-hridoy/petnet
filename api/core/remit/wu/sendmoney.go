package wu

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/bojanz/currency"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/api/util"
	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/serviceutil/logging"
)

type cacheSendRemit struct {
	StageReq    core.Remittance         `json:"stage_request"`
	StageRes    core.RemitResponse      `json:"stage_response"`
	PHStageReq  perahub.SMVRequest      `json:"ph_stage_request"`
	PHStageResp perahub.SMVResponseBody `json:"ph_stage_response"`
}

// stageSendMoney ...
// Two error objects returned to keep original error and core error to respond
func (s *Svc) stageSendMoney(ctx context.Context, r *core.Remittance, fee perahub.FIResponseBody) (rres *core.RemitResponse, err error, cErr error) {
	log := logging.FromContext(ctx)

	log.WithField("remittance", r).Trace("received")

	fmFlag := "N"
	if r.TargetDest {
		fmFlag = "F"
	}
	gnd := ""
	switch r.Remitter.Gender {
	case "Male":
		gnd = "M"
	case "Female":
		gnd = "F"
	default:
		err = status.Error(codes.InvalidArgument, "gender is required.")
		return nil, err, err
	}

	iso := s.iso(ctx)
	smvReq := perahub.SMVRequest{
		TransactionType: r.SendRemitType.Code,
		UserCode:        strconv.Itoa(r.Agent.UserID),
		FrgnRefNo:       r.DsaOrderID,

		SenderFirstName:  r.Remitter.FName,
		SenderMiddleName: r.Remitter.MdName,
		SenderLastName:   r.Remitter.LName,

		SenderContactPhone:      r.Remitter.Phone.CtyCode + r.Remitter.Phone.Number,
		SenderMobileCountryCode: r.Remitter.Mobile.CtyCode,
		SenderMobileNo:          r.Remitter.Mobile.Number,
		Email:                   r.Remitter.Email,

		SenderAddrLine1:        r.Remitter.Address.Address1,
		SenderAddrLine2:        r.Remitter.Address.Address2,
		SenderCity:             r.Remitter.Address.City,
		SenderState:            r.Remitter.Address.State,
		SenderPostalCode:       r.Remitter.Address.PostalCode,
		SenderAddrCountryCode:  r.Remitter.Address.Country,
		SenderAddrCurrencyCode: r.SourceAmount.CurrencyCode(),
		SenderAddrCountryName:  iso(r.Remitter.Address.Country),

		MyWUNumber: r.MyWUNumber,
		IDType:     r.Remitter.PrimaryID.IDType,
		IDCountry:  iso(r.Remitter.PrimaryID.Country),
		IDNumber:   r.Remitter.PrimaryID.Number,
		IDIssued:   r.Remitter.PrimaryID.Issued.String(),
		IDExpiry:   r.Remitter.PrimaryID.Expiry.String(),

		DateOfBirth:    r.Remitter.BirthDate.String(),
		CountryOfBirth: iso(r.Remitter.BirthCountry),
		Nationality:    iso(r.Remitter.Nationality),

		PromoCode:            r.Promo,
		ReceiverRelationship: r.Remitter.ReceiverRelation,
		Message:              strings.Split(r.Message, "\n"),
		TransactionPurpose:   r.TxnPurpose,
		SourceOfFunds:        r.Remitter.SourceFunds,
		SenderReason:         r.SendReason,

		Occupation:              r.Remitter.Employment.Occupation,
		EmploymentPositionLevel: r.Remitter.Employment.PositionLevel,
		Gender:                  gnd,

		AckFlag:       "X",
		MyWUEnrollTag: "none",

		ReceiverFirstName:  r.Receiver.FName,
		ReceiverMiddleName: r.Receiver.MdName,
		ReceiverLastName:   r.Receiver.LName,

		ReceiverContactPhone:      r.Receiver.Phone.Number,
		ReceiverMobileCountryCode: r.Receiver.Mobile.CtyCode,
		ReceiverMobileNo:          r.Receiver.Mobile.Number,

		ReceiverAddrLine1:        r.Receiver.Address.Address1,
		ReceiverAddrLine2:        r.Receiver.Address.Address2,
		ReceiverCity:             r.Receiver.Address.City,
		ReceiverState:            r.Receiver.Address.State,
		ReceiverPostalCode:       r.Receiver.Address.PostalCode,
		ReceiverAddrCountryCode:  r.Receiver.Address.Country,
		ReceiverAddrCurrencyCode: r.DestAmount.CurrencyCode(),

		PrincipalAmount: fee.OrigPrincipal,
		FixedAmountFlag: fmFlag,

		DestStateCode: r.DestState,
		DestCity:      r.DestCity,
		DestCurrency:  r.DestAmount.CurrencyCode(),
		DestCountry:   r.TransactionDetails.DestCtry,

		// todo(robin): making this static for now, live this should be provided by user
		// TerminalID:       phmw.GetTerminalID(ctx),
		// OperatorID:       phmw.GetOperatorID(ctx),
		TerminalID:       getTerminalID(ctx),
		OperatorID:       operatorID,
		RemoteTerminalID: remoteTerminalID,
		RemoteOperatorID: remoteOperatorID,
	}

	smvRes, err := s.ph.SendMoneyValidate(ctx, smvReq)
	if err != nil {
		logging.WithError(err, log).Error("send money validate error")
		return nil, err, handleWUError(err)
	}

	if r.DsaID == "" {
		r.DsaID = hydra.OrgID(ctx)
	}

	tx, err := currency.NewMinor("0", r.SourceAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid source amount currency")
		return nil, err, coreerror.NewCoreError(codes.Internal, "invalid source amount currency")
	}

	pcpl, err := currency.NewMinor(smvRes.Fin.OrigPnplAmt.String(), r.SourceAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid origin principle amount response")
		return nil, err, coreerror.NewCoreError(codes.Internal, "invalid origin principle amount response")
	}
	rmt, err := currency.NewMinor(smvRes.Fin.DestPcplAmt.String(), r.DestAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid dest principle amount response")
		return nil, err, coreerror.NewCoreError(codes.Internal, "invalid dest principle amount response")
	}

	gross, err := currency.NewMinor(smvRes.Fin.GrossTotal.String(), r.SourceAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid gross total response")
		return nil, err, coreerror.NewCoreError(codes.Internal, "invalid gross total response")
	}
	chg, err := currency.NewMinor(smvRes.Fin.Charges.String(), r.SourceAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid charges amount response")
		return nil, err, coreerror.NewCoreError(codes.Internal, "invalid charges amount total response")
	}

	txs := map[string]currency.Minor{}
	for k, v := range map[string]json.Number{
		"State":     smvRes.Fin.Taxes.StateTax,
		"County":    smvRes.Fin.Taxes.CountyTax,
		"Municipal": smvRes.Fin.Taxes.MuniTax,
	} {
		if v.String() == "" {
			continue
		}
		t, err := currency.NewMinor(v.String(), r.SourceAmount.CurrencyCode())
		if err != nil {
			emsg := fmt.Sprintf("invalid %s tax amount response", k)
			logging.WithError(err, log).
				WithFields(logrus.Fields{
					"amt": v.String(),
					"cur": r.SourceAmount.CurrencyCode(),
					"tax": k,
				}).Error(emsg)
			return nil, err, coreerror.NewCoreError(codes.Internal, emsg)
		}
		tx, _ = tx.Add(t)
		txs[k] = t
	}

	chgs := map[string]currency.Minor{static.WUCode: chg}
	res := &core.RemitResponse{
		PrincipalAmount:  pcpl,
		RemitAmount:      rmt,
		Taxes:            txs,
		Tax:              tx,
		Charges:          chgs,
		Charge:           chg,
		GrossTotal:       gross,
		PromoDescription: smvRes.Promotions.PromoDesc,
		PromoMessage:     smvRes.Promotions.PromoMessage,
	}
	r.GrossTotal = gross
	r.Tax = tx
	r.Charge = chg

	remit, err := json.Marshal(cacheSendRemit{
		StageReq:    *r,
		StageRes:    *res,
		PHStageReq:  smvReq,
		PHStageResp: *smvRes,
	})
	if err != nil {
		logging.WithError(err, log).Error("send money validate marshal request response error")
		return nil, err, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	rcReq := storage.RemitCache{
		DsaID:          r.DsaID,
		UserID:         r.UserID,
		RemcoID:        r.RemitPartner,
		RemType:        core.CreateRemType,
		PtnrRemType:    r.SendRemitType.Code,
		RemcoMemberID:  r.MyWUNumber,
		RemcoControlNo: smvRes.NewDetails.MTCN,
		Step:           storage.StageStep,
		Remit:          remit,
	}
	rc, err := s.st.CreateRemitCache(ctx, rcReq)
	if err != nil {
		logging.WithError(err, log).Error("send money validate create cache error")
		return nil, err, coreerror.NewCoreError(codes.Internal, coreerror.MsgDatabaseError)
	}

	res.TransactionID = rc.TxnID
	r.TransactionID = rc.TxnID
	return res, nil, nil
}

func (s *Svc) sendMoney(ctx context.Context, r core.ProcessRemit) (rres *core.ProcessRemit, cErr error) {
	log := logging.FromContext(ctx)
	var err error
	c := cacheSendRemit{}
	if err = json.Unmarshal(r.RemitCache.Remit, &c); err != nil {
		logging.WithError(err, log).Error("send money unmarshal remit cache error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}
	creq := c.StageReq
	phcreq := c.PHStageReq
	phcres := c.PHStageResp

	defer func() {
		_, err := util.RecordConfirmTxn(ctx, s.st.(*postgres.Storage), creq, util.ConfirmTxnOpts{
			TxnID:       r.TransactionID,
			TxnType:     storage.SendType,
			TxnErr:      err,
			PtnrRemType: creq.SendRemitType.Code,
		})
		if err != nil {
			log.Error(err)
		}
	}()

	sent, err := s.ph.SendMoneyStore(ctx, perahub.SMStoreRequest{
		FrgnRefNo:                 phcreq.FrgnRefNo,
		UserCode:                  strconv.Itoa(creq.Agent.UserID),
		CustomerCode:              phcreq.CustomerCode,
		SenderFirstName:           phcreq.SenderFirstName,
		SenderMiddleName:          phcreq.SenderMiddleName,
		SenderLastName:            phcreq.SenderLastName,
		SenderAddrCountryCode:     phcreq.SenderAddrCountryCode,
		SenderAddrCurrencyCode:    phcreq.SenderAddrCurrencyCode,
		SenderContactPhone:        phcreq.SenderContactPhone,
		SenderMobileCountryCode:   phcreq.SenderMobileCountryCode,
		SenderMobileNo:            phcreq.SenderMobileNo,
		SenderReason:              phcreq.SenderReason,
		Email:                     phcreq.Email,
		ReceiverNameType:          phcreq.ReceiverNameType,
		ReceiverFirstName:         phcreq.ReceiverFirstName,
		ReceiverMiddleName:        phcreq.ReceiverMiddleName,
		ReceiverLastName:          phcreq.ReceiverLastName,
		ReceiverAddrLine1:         phcreq.ReceiverAddrLine1,
		ReceiverAddrLine2:         phcreq.ReceiverAddrLine2,
		ReceiverCity:              phcreq.ReceiverCity,
		ReceiverState:             phcreq.ReceiverState,
		ReceiverPostalCode:        phcreq.ReceiverPostalCode,
		ReceiverAddrCountryCode:   phcreq.ReceiverAddrCountryCode,
		ReceiverAddrCurrencyCode:  phcreq.ReceiverAddrCurrencyCode,
		ReceiverContactPhone:      phcreq.ReceiverContactPhone,
		ReceiverMobileCountryCode: phcreq.ReceiverMobileCountryCode,
		ReceiverMobileNo:          phcreq.ReceiverMobileNo,
		DestinationCountryCode:    phcreq.DestCountry,
		DestinationCurrencyCode:   phcreq.DestCurrency,
		DestinationStateCode:      phcreq.DestStateCode,
		DestinationCityName:       phcreq.DestCity,
		TransactionType:           phcreq.TransactionType,
		PrincipalAmount:           phcres.Fin.DestPcplAmt,
		FixedAmountFlag:           phcreq.FixedAmountFlag,
		PromoCode:                 phcreq.PromoCode,
		Message:                   phcreq.Message,
		AddlServiceChg:            phcres.ServiceCode.AddlSvcChg,
		ComplianceDataBuffer:      phcres.Compliance.ComplianceBuf,
		OriginatingCity:           phcres.PaymentDetails.OrigCity,
		OriginatingState:          phcres.PaymentDetails.OrigState,
		MTCN:                      phcres.NewDetails.MTCN,
		NewMTCN:                   phcres.NewDetails.NewMTCN,
		ExchangeRate:              phcreq.ExchangeRate,
		Fin:                       phcres.Fin,
		Promo:                     phcres.Promotions,
		ComplianceDetails: perahub.ComplianceDetails{
			IDDetails: perahub.IDDetails{
				IDType:    phcreq.IDType,
				IDCountry: phcreq.IDCountry,
				IDNumber:  phcreq.IDNumber,
			},
			Birthdate:         phcreq.DateOfBirth,
			Occupation:        phcreq.Occupation,
			TransactionReason: phcreq.TransactionPurpose,
			CurrentAddress: perahub.Address{
				AddrLine1:  phcreq.SenderAddrLine1,
				AddrLine2:  phcreq.SenderAddrLine2,
				City:       phcreq.SenderCity,
				State:      phcreq.SenderState,
				PostalCode: phcreq.SenderPostalCode,
				Country:    phcreq.SenderAddrCountryName,
			},
			TxnRelationship:         phcreq.ReceiverRelationship,
			EmploymentPositionLevel: phcreq.EmploymentPositionLevel,
		},
		BankName:      phcreq.BankName,
		BankLocation:  phcreq.BankLocation,
		AccountNumber: phcreq.AccountNumber,
		BankCode:      phcreq.BankCode,
		AccountSuffix: phcreq.AccountSuffix,
		MyWUNumber:    phcres.PreferredCustomer.MyWUNumber,
		// todo(robin): making this static for now, live this should be provided by user
		// TerminalID:       phmw.GetTerminalID(ctx),
		// OperatorID:       phmw.GetOperatorID(ctx),
		TerminalID:       getTerminalID(ctx),
		OperatorID:       operatorID,
		RemoteTerminalID: remoteTerminalID,
		RemoteOperatorID: remoteOperatorID,
	})
	if err != nil {
		logging.WithError(err, log).Error("send money store error")
		return nil, handleWUError(err)
	}

	r.ControlNumber = sent.MTCN
	creq.ControlNo = sent.MTCN
	return &r, nil
}

func (s *Svc) iso(ctx context.Context) func(string) string {
	return func(c string) string {
		ct, err := s.stc.GetISO(ctx, c)
		if err != nil {
			return c
		}
		return ct.Country
	}
}
