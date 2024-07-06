package remitly

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
		logging.WithError(err, log).Error("get remit cache db error for rm")
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
	if _, err = s.ph.RMPayout(ctx, perahub.RMPayoutRequest{
		// clarify
		OrderNumber: "1",

		// internal/static
		RefNo:         random.InvitationCode(20),
		LocationCode:  creq.Agent.LocationID,
		LocationID:    "0",
		LocationName:  creq.Agent.LocationName,
		TxnDate:       t.Format("2006-01-02"),
		TxnType:       "2",
		ServiceCharge: "0",
		RmtLocID:      json.Number(creq.Agent.LocationID),
		DstAmount:     "0",
		BBAmt:         "0",
		MCRateID:      "0",
		MCRate:        "0",
		RateCat:       "0",
		RiskScore:     "1", // from reymar: this is currently internal to Perahub, we can leave this as blank for the meantime
		RiskCriteria:  "1",
		FormType:      "0",
		FormNumber:    "0",
		PayoutType:    "1",

		// inquire response cache
		PnplAmt:     creq.SourceAmount.Amount.Round().Number(),
		RemcoID:     json.Number(cres.RemcoID),
		TotalAmount: creq.SourceAmount.Amount.Round().Number(),
		SenderName:  cres.Result.SenderName,
		RcvName:     cn,
		SenderFName: creq.Remitter.FName,
		SenderMName: creq.Remitter.MdName,
		SenderLName: creq.Remitter.LName,
		RcvFName:    creq.Receiver.FName,
		RcvMName:    creq.Receiver.MdName,
		RcvLName:    creq.Receiver.LName,
		AgentID:     creq.Agent.AgentID,
		AgentCode:   creq.Agent.AgentCode,

		// userinput
		ClientRefNo:  creq.DsaOrderID,
		Gender:       perahub.GenderChar(creq.Receiver.Gender),
		ControlNo:    creq.ControlNo,
		CurrencyCode: creq.DestAmount.CurrencyCode(),
		IDNumber:     creq.Receiver.PrimaryID.Number,
		IDType:       creq.Receiver.PrimaryID.IDType,
		IDIssBy:      creq.Receiver.PrimaryID.Country,
		IDIssueDate: fmt.Sprintf("%s-%s-%s",
			creq.Receiver.PrimaryID.Issued.Year,
			creq.Receiver.PrimaryID.Issued.Month,
			creq.Receiver.PrimaryID.Issued.Day),
		IDExpDate: fmt.Sprintf("%s-%s-%s",
			creq.Receiver.PrimaryID.Expiry.Year,
			creq.Receiver.PrimaryID.Expiry.Month,
			creq.Receiver.PrimaryID.Expiry.Day),
		ContactNumber: creq.Receiver.Phone.Number,
		Address:       perahub.FormatAddress(creq.Receiver.Address.Address1, creq.Receiver.Address.Address2),
		City:          creq.Receiver.Address.City,
		Province:      creq.Receiver.Address.Province,
		Country:       creq.Receiver.Address.Country,
		ZipCode:       creq.Receiver.Address.PostalCode,
		State:         creq.Receiver.Address.State,
		Natl:          creq.Receiver.Nationality,
		BirthDate:     fmt.Sprintf("%s-%s-%s", creq.Receiver.BirthDate.Year, creq.Receiver.BirthDate.Month, creq.Receiver.BirthDate.Day),
		BirthCountry:  creq.Receiver.BirthCountry,
		Occupation:    creq.Receiver.Employment.Occupation,
		UserID:        strconv.Itoa(creq.Agent.UserID),
		CustomerID:    json.Number(strconv.Itoa(creq.Agent.UserID)),
		CurrencyID:    json.Number(perahub.CurrencyNumber(creq.DestAmount.CurrencyCode())),
		IsDomestic:    json.Number(strconv.Itoa(creq.TransactionDetails.IsDomestic)),
		CustomerName:  cn,
		RmtIPAddr:     creq.Agent.IPAddress,
		RmtUserID:     strconv.Itoa(creq.Agent.UserID),
		OrgnCtry:      creq.TransactionDetails.SrcCtry,
		DestCtry:      creq.TransactionDetails.DestCtry,
		PurposeTxn:    creq.TxnPurpose,
		SourceFund:    creq.Receiver.SourceFunds,
		RelationTo:    creq.Receiver.ReceiverRelation,
		BirthPlace:    creq.Receiver.BirthPlace,
		Barangay:      creq.Receiver.Address.Zone,
		IPAddr:        creq.Agent.IPAddress,
		DsaCode:       dsaCode,
		DsaTrxType:    dsaTrxType,
	}); err != nil {
		logging.WithError(err, log).Error("rm payout error")
		return nil, handleRemitlyError(err)
	}

	creq.Remitter.FName, creq.Remitter.MdName, creq.Remitter.LName = perahub.FormatName(cres.Result.SenderName)
	r.Processed = t
	r.ControlNumber = creq.ControlNo
	return &r, nil
}
