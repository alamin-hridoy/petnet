package wu

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
	"brank.as/petnet/api/storage/postgres"
	"brank.as/petnet/api/util"
	"brank.as/petnet/serviceutil/logging"
)

// Source Amounts must be in PHP.
func (s *Svc) StageCreateRemit(ctx context.Context, r core.Remittance) (rres *core.RemitResponse, cErr error) {
	log := logging.FromContext(ctx)

	r2 := &r
	var err error
	defer func() {
		_, err := util.RecordStageTxn(ctx, s.st.(*postgres.Storage), *r2, util.StageTxnOpts{
			TxnID:       r.TransactionID,
			TxnType:     storage.SendType,
			TxnErr:      err,
			PtnrRemType: r.SendRemitType.Code,
		})
		if err != nil {
			log.Error()
		}
	}()

	if s.st.OrderIDExists(ctx, r.DsaOrderID) {
		log.Error("order already exist")
		err = coreerror.NewCoreError(codes.AlreadyExists, fmt.Sprintf("a transaction with order id: %s, already exists", r.DsaOrderID))
		return nil, err
	}

	fmFlag := "N"
	if r.TargetDest {
		fmFlag = "F"
	}

	usrMsg := strings.Split(r.Message, "\n")
	fee, err := s.ph.FeeInquiry(ctx, perahub.FIRequest{
		FrgnRefNo:       r.DsaOrderID,
		PrincipalAmount: json.Number(r.SourceAmount.MinorUnits().String()),
		FixedAmountFlag: fmFlag,
		DestCountry:     r.TransactionDetails.DestCtry,
		DestCurrency:    r.DestAmount.CurrencyCode(),
		TransactionType: static.WUTxType(r.SendRemitType.Code),
		PromoCode:       r.Promo,
		Message:         usrMsg,
		// todo(robin): enable dynamic terminalid and operatorid once petnet has set it up for
		// sandbox
		// TerminalID:      phmw.GetTerminalID(ctx),
		// OperatorID:      phmw.GetOperatorID(ctx),
		TerminalID: getTerminalID(ctx),
		OperatorID: operatorID,
		UserCode:   strconv.Itoa(r.Agent.UserID),
	})
	if err != nil {
		logging.WithError(err, log).Error("inquire error for create")
		return nil, handleWUError(err)
	}

	var resp *core.RemitResponse
	switch r.SendRemitType.Code {
	case static.WUQuickPay:
		resp, err, cErr = s.stageQuickPay(ctx, r2, *fee)
	case static.WUSendMoney, static.WUDirectBank, static.WUMobileTransfer:
		resp, err, cErr = s.stageSendMoney(ctx, r2, *fee)
	default:
		err = coreerror.NewCoreError(codes.InvalidArgument, "Invalid remittance type")
		return nil, err
	}

	if cErr != nil {
		return nil, cErr
	}

	return resp, nil
}
