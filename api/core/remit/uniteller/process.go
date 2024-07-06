package uniteller

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
		logging.WithError(err, log).Error("get remit cache db error for unt")
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
	sfname, smname, slname := perahub.FormatName(cres.Result.SenderName)
	t := time.Now()
	if phmw.GetTerminalID(ctx) != "" {
		creq.Agent.DeviceID = phmw.GetTerminalID(ctx)
	}
	if _, err = s.ph.UNTPayout(ctx, perahub.UNTPayoutRequest{
		// inquire response cache
		ControlNumber:      cres.Result.ControlNumber,
		PrincipalAmount:    cres.Result.PrincipalAmount,
		Currency:           cres.Result.Currency,
		SenderName:         cres.Result.SenderName,
		ReceiverName:       cres.Result.ReceiverName,
		Address:            cres.Result.Address,
		City:               cres.Result.City,
		Country:            cres.Result.Country,
		ZipCode:            cres.Result.ZipCode,
		OriginatingCountry: cres.Result.OriginatingCountry,
		DestinationCountry: cres.Result.DestinationCountry,
		ContactNumber:      cres.Result.ContactNumber,

		// internal/static
		RiskScore:        "0",
		RiskCriteria:     "0",
		TrxType:          "2",
		ServiceCharge:    "0",
		DstAmount:        "0",
		FormType:         "0",
		FormNumber:       "0",
		PayoutType:       "1",
		BuyBackAmount:    "0",
		McRateId:         "0",
		McRate:           "0",
		TrxDate:          t.Format("2006-01-02"),
		RemcoID:          "20",
		RateCategory:     "0",
		RemoteLocationID: "3711",
		CurrencyID:       "1",

		// userinput
		LocationCode:      creq.Agent.LocationID,
		LocationID:        "0",
		IDIssuedBY:        creq.Receiver.PrimaryID.Country,
		Barangay:          creq.Receiver.Address.Zone,
		BirthCountry:      creq.Receiver.BirthCountry,
		LocationName:      creq.Agent.LocationName,
		ClientReferenceNo: creq.DsaOrderID,
		ReferenceNumber:   random.NumberString(18),
		Gender:            perahub.GenderChar(creq.Receiver.Gender),
		IDNumber:          creq.Receiver.PrimaryID.Number,
		IDType:            creq.Receiver.PrimaryID.IDType,
		IDDOIssue: fmt.Sprintf("%s-%s-%s",
			creq.Receiver.PrimaryID.Expiry.Year,
			creq.Receiver.PrimaryID.Expiry.Month,
			creq.Receiver.PrimaryID.Expiry.Day),
		IDExpDate: fmt.Sprintf("%s-%s-%s",
			creq.Receiver.PrimaryID.Issued.Year,
			creq.Receiver.PrimaryID.Issued.Month,
			creq.Receiver.PrimaryID.Issued.Day),
		Province:    creq.Receiver.Address.Province,
		State:       creq.Receiver.Address.State,
		Nationality: creq.Receiver.Nationality,
		BirthDate:   fmt.Sprintf("%s-%s-%s", creq.Receiver.BirthDate.Year, creq.Receiver.BirthDate.Month, creq.Receiver.BirthDate.Day),

		Occupation:         creq.Receiver.Employment.Occupation,
		UserID:             creq.Agent.UserID,
		CustomerID:         strconv.Itoa(creq.Agent.UserID),
		IsDomestic:         strconv.Itoa(creq.TransactionDetails.IsDomestic),
		CustomerName:       cn,
		TotalAmount:        creq.SourceAmount.Number(),
		RemoteIPAddress:    creq.Agent.IPAddress,
		RemoteUserID:       creq.Agent.UserID,
		PurposeTransaction: creq.TxnPurpose,
		SourceFund:         creq.Receiver.SourceFunds,
		RelationTo:         creq.Receiver.ReceiverRelation,
		BirthPlace:         creq.Receiver.BirthPlace,
		SendFName:          sfname,
		SendMName:          smname,
		SendLName:          slname,
		RecFName:           creq.Receiver.FName,
		RecMName:           creq.Receiver.MdName,
		RecLName:           creq.Receiver.LName,
		DeviceID:           creq.Agent.DeviceID,
		AgentID:            creq.Agent.AgentID,
		AgentCode:          creq.Agent.AgentCode,
		IPAddress:          creq.Agent.IPAddress,
		DsaCode:            dsaCode,
		DsaTrxType:         dsaTrxType,
		// clarify
		OrderNumber: "1",
	}); err != nil {
		logging.WithError(err, log).Error("unt payout error")
		return nil, handleUnitellerError(err)
	}

	creq.Remitter.FName, creq.Remitter.MdName, creq.Remitter.LName = perahub.FormatName(cres.Result.SenderName)
	r.Processed = t
	r.ControlNumber = creq.ControlNo
	return &r, nil
}
