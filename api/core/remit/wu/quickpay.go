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
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/random"
)

type cacheQuickPay struct {
	StageReq    core.Remittance     `json:"stage_request"`
	PHStageReq  perahub.QPVRequest  `json:"ph_stage_request"`
	PHStageResp perahub.QPVResponse `json:"ph_stage_response"`
}

// stageQuickPay ...
// Two error objects returned to keep original error and core error to respond
func (s *Svc) stageQuickPay(ctx context.Context, r *core.Remittance, fee perahub.FIResponseBody) (rres *core.RemitResponse, err error, cErr error) {
	log := logging.FromContext(ctx)

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

	qpReq := perahub.QPVRequest{
		FrgnRefNo: r.DsaOrderID,

		UserCode:     strconv.Itoa(r.Agent.UserID),
		CustomerCode: "",

		TransactionType: r.SendRemitType.Code,

		SenderFirstName:  r.Remitter.FName,
		SenderMiddleName: r.Remitter.MdName,
		SenderLastName:   r.Remitter.LName,
		SenderAddrLine1:  r.Remitter.Address.Address1,
		SenderAddrLine2:  r.Remitter.Address.Address2,
		SenderCity:       r.Remitter.Address.City,
		SenderState:      r.Remitter.Address.State,
		SenderPostalCode: r.Remitter.Address.PostalCode,

		SenderAddrCountry:  r.Remitter.Address.Country,
		SenderAddrCurrency: "",

		SenderAddrCountryName:   r.Remitter.Address.Country,
		SenderContactPhone:      r.Remitter.Phone.CtyCode + r.Remitter.Phone.Number,
		SenderMobileCountryCode: r.Remitter.Mobile.CtyCode,
		SenderMobileNo:          r.Remitter.Mobile.Number,
		Email:                   r.Remitter.Email,
		Gender:                  gnd,

		BirthCountry: r.Remitter.BirthCountry,
		Birthdate:    r.Remitter.BirthDate.String(),

		Nationality: r.Remitter.Nationality,
		IDType:      r.Remitter.PrimaryID.IDType,
		IDCountry:   r.Remitter.PrimaryID.Country,
		IDNumber:    r.Remitter.PrimaryID.Number,
		IDIssued:    r.Remitter.PrimaryID.Issued.String(),
		IDExpiry:    r.Remitter.PrimaryID.Expiry.String(),

		Occupation:              r.Remitter.Employment.Occupation,
		EmploymentPositionLevel: r.Remitter.Employment.PositionLevel,

		FundSource:           r.Remitter.SourceFunds,
		ReceiverRelationship: r.Remitter.ReceiverRelation,
		TransactionPurpose:   r.TxnPurpose,
		SendingReason:        r.SendReason,
		MyWUNumber:           r.MyWUNumber,
		MyWUEnrollTag:        "none",
		CompanyName:          r.Business.Name,
		// TODO
		// CompanyCode:             r.Business.,
		CompanyAccountCode: r.Business.Account,
		ReferenceNo:        r.Business.ControlNo,
		PromoCode:          r.Promo,
		PrincipalAmount:    fee.OrigPrincipal,
		ExchangeRate:       fee.ExchangeRate,
		FixedAmountFlag:    fmFlag,
		DestCountry:        r.Receiver.Phone.CtyCode,
		DestCurrency:       r.DestAmount.CurrencyCode(),
		DestState:          r.DestState,
		DestCity:           r.DestCity,
		Message:            strings.Split(r.Message, "\n"),
		BankName:           r.DestAccount.BIC,
		BankLocation:       r.DestAccount.BIC,
		BankCode:           r.DestAccount.BIC,
		AccountSuffix:      r.DestAccount.AcctSfx,
		AccountNumber:      r.DestAccount.AcctNo,
		IsOnBehalf:         "",
		AckFlag:            "X",
		// todo(robin): making this static for now, live this should be provided by user
		// TerminalID:       phmw.GetTerminalID(ctx),
		// OperatorID:       phmw.GetOperatorID(ctx),
		TerminalID:       getTerminalID(ctx),
		OperatorID:       operatorID,
		RemoteTerminalID: remoteTerminalID,
		RemoteOperatorID: remoteOperatorID,
	}

	qpRes, err := s.ph.QuickPayValidate(ctx, qpReq)
	if err != nil {
		logging.WithError(err, log).Error("quick pay validate error")
		return nil, err, handleWUError(err)
	}

	tx, err := currency.NewMinor("0", r.SourceAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid source amount currency")
		return nil, err, coreerror.NewCoreError(codes.Internal, "invalid source amount currency")
	}

	pcpl, err := currency.NewMinor(qpRes.Fin.OrigPnplAmt.String(), r.SourceAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid origin principle amount response")
		return nil, err, coreerror.NewCoreError(codes.Internal, "invalid origin principle amount response")
	}
	rmt, err := currency.NewMinor(qpRes.Fin.DestPcplAmt.String(), r.DestAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid dest principle amount response")
		return nil, err, coreerror.NewCoreError(codes.Internal, "invalid dest principle amount response")
	}
	gross, err := currency.NewMinor(qpRes.Fin.GrossTotal.String(), r.DestAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid gross total response")
		return nil, err, coreerror.NewCoreError(codes.Internal, "invalid gross total response")
	}
	chg, err := currency.NewMinor(qpRes.Fin.Charges.String(), r.SourceAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid charges amount response")
		return nil, err, coreerror.NewCoreError(codes.Internal, "invalid charges amount total response")
	}

	txs := map[string]currency.Minor{}
	r.Charge = chg
	r.Tax = tx
	for k, v := range map[string]json.Number{
		"State":     qpRes.Fin.Taxes.StateTax,
		"County":    qpRes.Fin.Taxes.CountyTax,
		"Municipal": qpRes.Fin.Taxes.MuniTax,
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
				}).Error("parsing amount")
			return nil, err, coreerror.NewCoreError(codes.Internal, emsg)
		}
		tx, _ = tx.Add(t)
		txs[k] = t
	}
	chgs := map[string]currency.Minor{static.WUCode: chg}

	r.GrossTotal = gross
	remit, err := json.Marshal(cacheQuickPay{
		StageReq:    *r,
		PHStageReq:  qpReq,
		PHStageResp: *qpRes,
	})
	if err != nil {
		logging.WithError(err, log).Error("quick pay validate marshal request response error")
		return nil, err, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	rcReq := storage.RemitCache{
		DsaID:          r.DsaID,
		UserID:         r.UserID,
		RemcoID:        r.RemitPartner,
		RemType:        core.CreateRemType,
		PtnrRemType:    r.SendRemitType.Code,
		RemcoMemberID:  r.MyWUNumber,
		RemcoControlNo: qpRes.NewDetails.MTCN,
		Step:           storage.StageStep,
		Remit:          remit,
	}
	rc, err := s.st.CreateRemitCache(ctx, rcReq)
	if err != nil {
		logging.WithError(err, log).Error("quick pay validate create cache error")
		return nil, err, coreerror.NewCoreError(codes.Internal, coreerror.MsgDatabaseError)
	}

	r.TransactionID = rc.TxnID
	return &core.RemitResponse{
		PrincipalAmount:  pcpl,
		RemitAmount:      rmt,
		Taxes:            txs,
		Tax:              tx,
		Charges:          chgs,
		Charge:           chg,
		GrossTotal:       gross,
		PromoDescription: qpRes.Promotions.PromoDesc,
		PromoMessage:     qpRes.Promotions.Message,
		TransactionID:    rc.TxnID,
	}, nil, nil
}

func (s *Svc) sendQuickPay(ctx context.Context, r core.ProcessRemit) (rres *core.ProcessRemit, cErr error) {
	log := logging.FromContext(ctx)

	var err error
	c := cacheQuickPay{}
	if err = json.Unmarshal(r.RemitCache.Remit, &c); err != nil {
		logging.WithError(err, log).Error("send quick pay unmarshal remit cache error")
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

	sent, err := s.ph.QuickPayStore(ctx, perahub.QPSRequest{
		FrgnRefNo:               phcreq.FrgnRefNo,
		UserCode:                strconv.Itoa(creq.Agent.UserID),
		CustomerCode:            phcreq.CustomerCode,
		SenderFirstName:         phcreq.SenderFirstName,
		SenderMiddleName:        phcreq.SenderMiddleName,
		SenderLastName:          phcreq.SenderLastName,
		SenderAddrCountryCode:   phcreq.SenderAddrCountry,
		SenderAddrCurrencyCode:  phcreq.SenderAddrCurrency,
		SenderContactPhone:      phcreq.SenderContactPhone,
		SenderMobileCountryCode: phcreq.SenderMobileCountryCode,
		SenderMobileNo:          phcreq.SenderMobileNo,
		SendingReason:           phcreq.SendingReason,
		Email:                   phcreq.Email,
		CompanyName:             phcreq.CompanyName,
		CompanyCode:             phcreq.CompanyCode,
		CompanyAccountCode:      phcreq.CompanyAccountCode,
		ReferenceNo:             phcreq.ReferenceNo,
		DestCountry:             phcreq.DestCountry,
		DestCurrency:            phcreq.DestCurrency,
		DestState:               phcreq.DestState,
		DestCity:                phcreq.DestCity,
		TransactionType:         phcreq.TransactionType,
		PrincipalAmount:         phcres.Fin.OrigPnplAmt,
		FixedAmountFlag:         phcreq.FixedAmountFlag,
		PromoCode:               phcreq.PromoCode,
		Message:                 phcreq.Message,
		AddlServiceCharges:      phcres.ServiceCode.AddlServiceCharges,
		ComplianceDataBuffer:    phcres.Compliance.ComplianceDataBuffer,
		OrigCity:                phcres.PaymentDetails.OrigCity,
		OrigState:               phcres.PaymentDetails.OrigState,
		MTCN:                    phcres.NewDetails.MTCN,
		NewMTCN:                 phcres.NewDetails.NewMTCN,
		ExchangeRate:            phcreq.ExchangeRate,
		Fin:                     perahub.Financials(phcres.Fin),
		Promo: perahub.Promotions{
			PromoDesc:       phcres.Promotions.PromoDesc,
			PromoMessage:    phcres.Promotions.Message,
			SenderPromoCode: phcres.Promotions.SenderPromoCode,
		},
		Compliance: perahub.ComplianceDetails{
			IDDetails: perahub.IDDetails{
				IDType:    phcreq.IDType,
				IDCountry: phcreq.IDCountry,
				IDNumber:  phcreq.IDNumber,
			},
			IDIssued:          phcreq.IDIssued,
			IDExpiry:          phcreq.IDExpiry,
			Birthdate:         phcreq.Birthdate,
			BirthCountry:      phcreq.BirthCountry,
			Nationality:       phcreq.Nationality,
			SourceFunds:       phcreq.FundSource,
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
			IActOnMyBehalf:          phcreq.IsOnBehalf,
			GalacticID:              *phcreq.GalacticID,
			EmploymentPositionLevel: phcreq.EmploymentPositionLevel,
		},
		BankName:      phcreq.BankName,
		BankLocation:  phcreq.BankLocation,
		AccountNumber: phcreq.AccountNumber,
		BankCode:      phcreq.BankCode,
		AccountSuffix: phcreq.AccountSuffix,
		MyWUNumber:    phcres.PreferredCustomer.MyWUNumber,
		MyWuPoints:    "",
		MyWUEnrollTag: phcreq.MyWUEnrollTag,
		// todo(robin): making this static for now, live this should be provided by user
		// TerminalID:       phmw.GetTerminalID(ctx),
		// OperatorID:       phmw.GetOperatorID(ctx),
		TerminalID:        getTerminalID(ctx),
		OperatorID:        operatorID,
		RemoteTerminalID:  remoteTerminalID,
		RemoteOperatorID:  remoteOperatorID,
		ClientReferenceNo: random.InvitationCode(20),
	})
	if err != nil {
		logging.WithError(err, log).Error("quickpay store")
		r.RemitCache.Step = storage.TxnStep(storage.FailStatus)
		if _, dbErr := s.st.UpdateRemitCache(ctx, r.RemitCache); dbErr != nil {
			logging.WithError(dbErr, log).Error("send quick pay update cache error")
		}

		logging.WithError(err, log).Error("send quick pay error")
		return nil, handleWUError(err)
	}

	r.ControlNumber = sent.MTCN
	creq.ControlNo = sent.MTCN
	return &r, nil
}
