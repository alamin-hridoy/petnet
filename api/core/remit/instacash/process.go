package instacash

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

		logging.WithError(err, log).Error("get remit cache db error for ic")
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
	if _, err = s.ph.InstaCashPayout(ctx, perahub.InstaCashPayoutRequest{
		// inquire response cache
		ControlNumber:      cres.Result.ControlNumber,
		ReferenceNumber:    cres.Result.ReferenceNumber,
		OriginatingCountry: cres.Result.OriginatingCountry,
		DestinationCountry: cres.Result.DestinationCountry,
		SenderName:         cres.Result.SenderName,
		ReceiverName:       cres.Result.ReceiverName,
		PrincipalAmount:    string(cres.Result.PrincipalAmount),
		PurposeTransaction: cres.Result.Purpose,

		// internal/static
		RemcoID:          "16",
		PayoutType:       "1",
		DstAmount:        "0",
		TrxType:          "2",
		ServiceCharge:    "0",
		RiskScore:        "1",
		RiskCriteria:     "1",
		TrxDate:          t.Format("2006-01-02"),
		IsDomestic:       "1",
		LocationID:       "0",
		RemoteLocationID: "371",
		McRate:           "0",
		McRateId:         "0",
		RateCategory:     "0",
		BuyBackAmount:    "0",
		FormType:         "0",
		FormNumber:       "0",
		CurrencyID:       "1",

		// userinput
		UserID:        creq.Agent.UserID,
		CustomerID:    strconv.Itoa(creq.Agent.UserID),
		CustomerName:  cn,
		TotalAmount:   string(cres.Result.PrincipalAmount),
		SourceFund:    creq.Receiver.SourceFunds,
		Occupation:    creq.Receiver.Employment.Occupation,
		RelationTo:    creq.Receiver.ReceiverRelation,
		BirthDate:     fmt.Sprintf("%s-%s-%s", creq.Receiver.BirthDate.Year, creq.Receiver.BirthDate.Month, creq.Receiver.BirthDate.Day),
		BirthPlace:    creq.Receiver.BirthPlace,
		BirthCountry:  creq.Receiver.BirthCountry,
		IDType:        creq.Receiver.PrimaryID.IDType,
		IDNumber:      creq.Receiver.PrimaryID.Number,
		Address:       perahub.FormatAddress(creq.Receiver.Address.Address1, creq.Receiver.Address.Address2),
		Barangay:      creq.Receiver.Address.Zone,
		City:          creq.Receiver.Address.City,
		Province:      creq.Receiver.Address.Province,
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

		ClientReferenceNo: creq.DsaOrderID,
		RemoteUserID:      creq.Agent.UserID,
		RemoteIPAddress:   creq.Agent.IPAddress,
		IPAddress:         creq.Agent.IPAddress,
		ZipCode:           creq.Receiver.Address.PostalCode,
		LocationName:      creq.Agent.LocationName,
		DsaCode:           dsaCode,
		DsaTrxType:        dsaTrxType,
	}); err != nil {
		logging.WithError(err, log).Error("ic payout error")
		return nil, handleInstaCashError(err)
	}

	r.Processed = t
	r.ControlNumber = creq.ControlNo
	return &r, nil
}
