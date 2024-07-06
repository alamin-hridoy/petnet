package wu

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/bojanz/currency"
	"google.golang.org/grpc/codes"

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

const (
	terminalID       = "PH259AMT001A" // petnet will provide this operator id
	operatorID       = "drp"          // petnet will provide this terminal id
	remoteOperatorID = "1"
	remoteTerminalID = "2"
	galacticID       = "1000000000028668380"
	exchangeRate     = "1"
	frgnRefNo        = "4d0bad65d8a55d3f81f6"
)

func getTerminalID(ctx context.Context) string {
	// TODO(vitthal): getting error if drp terminal ID is passed to WU
	/*tID := phmw.GetTerminalID(ctx)
	if tID == "" {
		return terminalID
	}

	return tID*/

	return terminalID
}

func (s *Svc) StageDisburseRemit(ctx context.Context, rmt core.Remittance) (rres *core.Remittance, cErr error) {
	log := logging.FromContext(ctx)
	// err variable used in defer to save error in db
	var err error
	defer func() {
		_, err := util.RecordStageTxn(ctx, s.st.(*postgres.Storage), rmt, util.StageTxnOpts{
			TxnID:       rmt.TransactionID,
			TxnType:     storage.DisburseType,
			TxnErr:      err,
			PtnrRemType: static.Payout,
		})
		if err != nil {
			logging.WithError(err, log).Error("record transaction error")
		}
	}()
	if s.st.OrderIDExists(ctx, rmt.DsaOrderID) {
		log.Error("order already exist")
		cErr = coreerror.NewCoreError(codes.AlreadyExists, fmt.Sprintf("a transaction with order id: %s, already exists", rmt.DsaOrderID))
		err = cErr
		return nil, cErr
	}

	res, err := s.ph.RMSearch(ctx, perahub.RMSearchRequest{
		FrgnRefNo:    random.InvitationCode(20),
		MTCN:         rmt.ControlNo,
		DestCurrency: rmt.DestAmount.CurrencyCode(),
		// todo(robin): making this static for now, live this should be provided by user
		// OperatorID:       phmw.GetOperatorID(ctx),

		TerminalID: getTerminalID(ctx),
		OperatorID: operatorID,
	})
	if err != nil {
		logging.WithError(err, log).Error("inquire error for disburse")
		return nil, handleWUError(err)
	}

	if res.Txn.PayStatus != perahub.AwaitPayment {
		log.Error("not awaiting disbursement")
		cErr = coreerror.NewCoreError(codes.InvalidArgument, "Remittance is not awaiting disbursement")
		err = cErr
		return nil, cErr
	}

	gamt, err := currency.NewAmount(res.GrossPayout.String(), rmt.DestAmount.CurrencyCode())
	if err != nil {
		logging.WithError(err, log).Error("invalid gross amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid gross amount")
	}

	rmt.GrossTotal = currency.ToMinor(gamt.Round())
	amt, err := currency.NewMinor(
		res.Txn.Financials.Principal.String(), rmt.DestAmount.CurrencyCode(),
	)
	if err != nil {
		logging.WithError(err, log).Error("invalid principal amount")
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "Invalid principal amount")
	}

	rmt.SourceAmount = amt
	zamt := core.MustMinor("0", "PHP")
	rmt.Tax = zamt
	rmt.Charge = zamt

	rmt.Remitter.FName, rmt.Remitter.LName = res.Txn.Sender.Name.FirstName, res.Txn.Sender.Name.LastName
	cache, err := json.Marshal(cacheDisburse{StageReq: rmt, StageResp: *res})
	if err != nil {
		logging.WithError(err, log).Error("caching request marshal error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	rc, err := s.st.CreateRemitCache(ctx, storage.RemitCache{
		DsaID:             rmt.DsaID,
		UserID:            rmt.UserID,
		RemcoID:           rmt.RemitPartner,
		RemType:           core.DisburseRemType,
		PtnrRemType:       static.Payout,
		RemcoMemberID:     rmt.Receiver.PartnerMemberID,
		RemcoControlNo:    res.Txn.Mtcn,
		RemcoAltControlNo: res.Txn.NewMtcn,
		Step:              storage.StageStep,
		Remit:             cache,
	})
	if err != nil {
		logging.WithError(err, log).Error("creating remit db cache error")
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDatabaseError)
	}

	rmt.TransactionID = rc.TxnID
	return &rmt, nil
}
