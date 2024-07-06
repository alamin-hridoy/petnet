package wu

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/bojanz/currency"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/random"
)

func (s *Svc) Search(ctx context.Context, r core.SearchRemit) (*core.SearchRemit, error) {
	log := logging.FromContext(ctx)

	rmt, err := s.ph.RMSearch(ctx, perahub.RMSearchRequest{
		FrgnRefNo:    random.InvitationCode(20),
		MTCN:         r.ControlNo,
		DestCurrency: r.DestCurrency,
		TerminalID:   getTerminalID(ctx),
		OperatorID:   "drp",
	})
	if err != nil {
		logging.WithError(err, log).Error("inquire error for search")
		return nil, handleWUError(err)
	}

	srch := core.SearchRemit{
		DsaID:        r.DsaID,
		UserID:       r.UserID,
		RemitPartner: r.RemitPartner,
		RemitType:    r.RemitType,
		DestCurrency: r.DestCurrency,
		ControlNo:    r.ControlNo,
		Remitter: core.Contact{
			FirstName: rmt.Txn.Sender.Name.FirstName,
			LastName:  rmt.Txn.Sender.Name.LastName,
			Address: core.Address{
				Address1:   rmt.Txn.Sender.Address.Street,
				City:       rmt.Txn.Sender.Address.City,
				State:      rmt.Txn.Sender.Address.State,
				PostalCode: rmt.Txn.Sender.Address.PostalCode,
				Country:    rmt.Txn.Sender.Address.CountryCode.IsoCode.Country,
			},
			Phone: core.PhoneNumber{
				Number: rmt.Txn.Sender.Phone,
			},
			Mobile: core.PhoneNumber{
				CtyCode: rmt.Txn.Sender.MobileDetails.CountryCode,
				Number:  rmt.Txn.Sender.MobileDetails.Number,
			},
		},
		Receiver: core.Contact{
			FirstName: rmt.Txn.Receiver.Name.FirstName,
			LastName:  rmt.Txn.Receiver.Name.LastName,
			Address: core.Address{
				Address1:   rmt.Txn.Receiver.Address.Street,
				City:       rmt.Txn.Receiver.Address.City,
				State:      rmt.Txn.Receiver.Address.State,
				PostalCode: rmt.Txn.Receiver.Address.PostalCode,
				Country:    rmt.Txn.Receiver.Address.CountryCode.IsoCode.Country,
			},
			Phone: core.PhoneNumber{
				Number: rmt.Txn.Receiver.Phone,
			},
			Mobile: core.PhoneNumber{
				CtyCode: rmt.Txn.Receiver.MobileDetails.CountryCode,
				Number:  rmt.Txn.Receiver.MobileDetails.Number,
			},
		},
		SentCountry: rmt.Txn.Payment.OrigCountry.IsoCode.Country,
		DestCity:    rmt.Txn.Payment.ExpectedPayoutLocation.City,
		DestState:   rmt.Txn.Payment.ExpectedPayoutLocation.State,
		Status:      rmt.Txn.PayStatus,
	}

	spl := strings.Fields(rmt.Txn.FilingTime)
	if len(spl) == 2 {
		// Splitting enables us to turn "0204A" (WU format) into "0204AM" so time can parse it.
		// Without resorting to violence like regex.
		tt, tz := spl[0], spl[1]
		tm, err := time.Parse("01-02-06 0304PM MST", rmt.Txn.FilingDate+tt+"M "+tz)
		if err != nil {
			logging.WithError(err, log).Error("filing time parse error")
			return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
		}
		srch.TxnCompletedTime = tm
	}

	remit, err := currency.NewMinor(
		rmt.Txn.Financials.GrossTotal.String(),
		rmt.Txn.Payment.OrigCountry.IsoCode.Currency,
	)
	if err != nil {
		logging.WithError(err, log).Error("pay amount parsing")
		return nil, status.Error(codes.Internal, "failed to load remittance data")
	}
	srch.RemitAmount = remit

	srch.DisburseAmount, err = currency.NewMinor(
		rmt.Txn.Financials.PayAmount.String(),
		rmt.Txn.Payment.SenderDestCountry.IsoCode.Currency,
	)
	if err != nil {
		logging.WithError(err, log).Error("pay amount parsing")
		return nil, status.Error(codes.Internal, "failed to load remittance data")
	}

	tot, chg, err := fillCharges(ctx, rmt.Txn.Financials, rmt.DST,
		rmt.Txn.Payment.DestCountry.IsoCode.Currency)
	if err != nil {
		return nil, err
	}
	srch.Charges = chg
	srch.Charge = tot

	return &srch, nil
}

func fillCharges(ctx context.Context, fin perahub.RMSFinancials, dst json.Number, curr string,
) (currency.Minor, map[string]currency.Minor, error) {
	log := logging.FromContext(ctx)

	tot, err := currency.NewMinor("0", curr)
	if err != nil {
		logging.WithError(err, log).WithField("currency", curr).Error("parsing")
		return tot, nil, status.Error(codes.Internal, "failed to load remittance data")
	}
	c := make(map[string]currency.Minor)
	for k, v := range map[string]string{
		static.WUCode: fin.Charges.String(),
		"Tolls":       fin.Tolls.String(),
		"DST":         dst.String(),
	} {
		amt, err := currency.NewMinor(v, curr)
		if err != nil {
			logging.WithError(err, log).WithField("value", v).WithField("charge", k).Error("parsing")
			return tot, nil, status.Error(codes.Internal, "failed to load remittance data")
		}
		c[k] = amt
		tot, _ = tot.Add(amt)
	}
	return tot, c, nil
}
