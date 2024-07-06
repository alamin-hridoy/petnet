package ussc

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
		logging.WithError(err, log).Error("get remit cache db error for ussc")
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

	t := time.Now()
	if _, err = s.ph.USSCPayout(ctx, perahub.USSCPayoutRequest{
		// userinput
		LocationID:        "0",
		LocationName:      creq.Agent.LocationName,
		UserID:            creq.Agent.UserID,
		ClientReferenceNo: creq.DsaOrderID,
		RemoteLocationID:  json.Number(creq.Agent.LocationID),
		RemoteUserID:      json.Number(strconv.Itoa(creq.Agent.UserID)),
		DstAmount:         json.Number(creq.DestAmount.Amount.Round().Number()),
		RemoteIPAddress:   creq.Agent.IPAddress,
		CustomerID:        strconv.Itoa(creq.Agent.UserID),
		SourceFund:        creq.Receiver.SourceFunds,
		Occupation:        creq.Receiver.Employment.Occupation,
		RelationTo:        creq.Receiver.ReceiverRelation,
		BirthPlace:        creq.Receiver.BirthPlace,
		BirthCountry:      creq.Receiver.BirthCountry,
		IDType:            creq.Receiver.PrimaryID.IDType,
		IDNumber:          creq.Receiver.PrimaryID.Number,
		Barangay:          creq.Receiver.Address.Zone,
		City:              creq.Receiver.Address.City,
		Province:          creq.Receiver.Address.Province,
		ZipCode:           creq.Receiver.Address.PostalCode,
		Country:           creq.Receiver.Address.Country,
		ContactNumber:     creq.Receiver.Phone.Number,
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
		BirthDate: fmt.Sprintf("%s-%s-%s", creq.Receiver.BirthDate.Year, creq.Receiver.BirthDate.Month, creq.Receiver.BirthDate.Day),
		Address:   perahub.FormatAddress(creq.Receiver.Address.Address1, creq.Receiver.Address.Address2),

		// inquire response
		RemcoID:            cres.RemcoID,
		CustomerName:       cres.Result.RcvName,
		ServiceCharge:      cres.Result.ServiceCharge,
		TotalAmount:        cres.Result.TotalAmount,
		PurposeTransaction: cres.Result.PurposeTransaction,
		SenderName:         cres.Result.SenderName,
		ReceiverName:       cres.Result.RcvName,
		SenderFirstName:    cres.Result.SenderFirstName,
		SenderLastName:     cres.Result.SenderLastName,
		ReceiverFirstName:  cres.Result.ReceiverFirstName,
		ReceiverLastName:   cres.Result.ReceiverLastName,

		PrincipalAmount: json.Number(cres.Result.PrincipalAmount),
		ControlNumber:   cres.Result.ControlNo,
		ReferenceNumber: cres.Result.RefNo,
		BranchCode:      creq.Agent.LocationID,

		// internal/static data
		TrxDate:            t.Format("2006-01-02"),
		OriginatingCountry: "Philippines",
		DestinationCountry: "PH",
		IsDomestic:         strconv.Itoa(creq.TransactionDetails.IsDomestic),
		TrxType:            "2",
		PayoutType:         "1",
		RiskScore:          "1",
		CurrencyID:         "1",
		FormType:           "0",
		FormNumber:         "0",
		McRate:             "0",
		BuyBackAmount:      "0",
		RateCategory:       "0",
		McRateId:           "0",
		DsaCode:            dsaCode,
		DsaTrxType:         dsaTrxType,
	}); err != nil {
		logging.WithError(err, log).Error("ussc payout error")
		return nil, handleUSSCError(err)
	}
	creq.Remitter.FName, creq.Remitter.MdName, creq.Remitter.LName = perahub.FormatName(cres.Result.SenderName)
	r.Processed = t
	r.ControlNumber = cres.Result.ControlNo
	return &r, nil
}

func (s *Svc) send(ctx context.Context, r core.ProcessRemit) (rres *core.ProcessRemit, cErr error) {
	log := logging.FromContext(ctx)
	var err error
	c := cacheSendRemit{}
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

	addr := perahub.NonexAddress{
		Address1: creq.Remitter.Address.Address1,
		Address2: creq.Remitter.Address.Address2,
		Barangay: creq.Remitter.Address.Zone,
		City:     creq.Remitter.Address.City,
		Province: creq.Remitter.Address.Province,
		ZipCode:  creq.Remitter.Address.PostalCode,
		Country:  creq.Remitter.Address.Country,
	}
	uscSendRes, err := s.ph.USSCSendMoney(ctx, perahub.USSCSendRequest{
		// internal/static
		TrxDate:            t.Format("2006-01-02"),
		RemcoID:            "10",
		TrxType:            "2",
		IsDomestic:         "1",
		LocationID:         "0",
		LocationName:       "loc-name", // todo: will change, gotten from petnet
		RemoteLocationID:   json.Number(creq.Agent.LocationID),
		DstAmount:          "0",
		RiskScore:          "1",
		RiskCrt:            "1",
		FormType:           "OR",
		FormNumber:         "MAL00001",
		PayoutType:         "1",
		CurrencyID:         "1",
		ControlNo:          random.NumberString(10),
		BranchCode:         creq.Agent.OutletCode,
		McRate:             "0",
		BbAmount:           "0",
		RateCtg:            "0",
		McRateID:           "0",
		OriginatingCountry: creq.TransactionDetails.SrcCtry,
		DestinationCountry: creq.TransactionDetails.DestCtry,

		// inquire response cache
		ReferenceNo:     cres.Result.RefNo,
		PrincipalAmount: cres.Result.PnplAmount,
		TotalAmount:     cres.Result.TotAmount,
		ServiceCharge:   cres.Result.ServiceCharge,

		// userinput
		UserID:             json.Number(strconv.Itoa(creq.Agent.UserID)),
		CustomerID:         strconv.Itoa(creq.Agent.UserID),
		CustomerName:       rn,
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
		CurrentAddress:     addr,
		PermanentAddress:   addr,
		SenderName:         sn,
		ReceiverName:       rn,
		ClientRefNo:        creq.DsaOrderID,
		SendFName:          creq.Remitter.FName,
		SendMName:          creq.Remitter.MdName,
		SendLName:          creq.Remitter.LName,
		RecFName:           creq.Receiver.FName,
		RecMName:           creq.Receiver.MdName,
		RecLName:           creq.Receiver.LName,
		RecConNo:           creq.Receiver.Phone.Number,
		KycVer:             true,
		Gender:             perahub.GenderChar(creq.Remitter.Gender),
		DsaCode:            phmw.GetDsaCode(ctx),
		DsaTrxType:         util.GetPerahubTrxType(ctx),
	})
	if err != nil {
		logging.WithError(err, log).Error("ussc send money error")
		return nil, handleUSSCError(err)
	}
	r.Processed = t
	r.ControlNumber = uscSendRes.Result.ControlNo
	creq.ControlNo = uscSendRes.Result.ControlNo
	return &r, nil
}
