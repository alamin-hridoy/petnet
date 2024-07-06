package terminal

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/integration/perahub"

	pnpb "brank.as/petnet/gunk/drp/v1/partner"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
)

func (s *Svc) ListRemit(ctx context.Context, req *tpb.ListRemitRequest) (*tpb.ListRemitResponse, error) {
	var ft time.Time
	var ut time.Time
	if err := validation.ValidateStruct(req,
		validation.Field(&req.Offset, validation.Min(0)),
		validation.Field(&req.Limit, validation.Min(0)),
		validation.Field(&req.SortOrder, validation.Min(0)),
		validation.Field(&req.SortByColumn, validation.Min(0)),
		validation.Field(&req.ControlNumbers),
		validation.Field(&req.From, validation.When(req.From != "", validation.By(func(interface{}) error {
			var err error
			ft, err = time.Parse("2006-01-02", req.From)
			if err != nil {
				return fmt.Errorf("'from' should be in the format YYYY-MM-DD")
			}
			return nil
		}))),
		validation.Field(&req.Until, validation.When(req.Until != "", validation.By(func(interface{}) error {
			var err error
			ut, err = time.Parse("2006-01-02", req.Until)
			if err != nil {
				return fmt.Errorf("'until' should be in the format YYYY-MM-DD")
			}
			return nil
		}))),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	res, err := s.remit.ListRemit(ctx, core.FilterList{
		From:           ft,
		Until:          ut,
		Limit:          int(req.Limit),
		Offset:         int(req.Offset),
		SortOrder:      tpb.SortOrder_name[int32(req.SortOrder)],
		SortByColumn:   tpb.SortByColumn_name[int32(req.SortByColumn)],
		ControlNumbers: req.ControlNumbers,
		ExcludePartner: req.ExcludePartner,
		ExcludeType:    req.ExcludeType,
	})
	if err != nil {
		switch t := err.(type) {
		case *perahub.Error:
			if t.Type == perahub.PartnerError {
				return nil, perahub.GRPCError(t.GRPCCode, "partner error", &pnpb.Error{
					Code:    t.Code,
					Message: t.Msg,
				})
			}
			return nil, status.Errorf(codes.Internal, "internal error occurred")
		}
		if status.Code(err) != codes.Unknown {
			return nil, err
		}
		return nil, status.Error(codes.Internal, "failed to retrieve historical transactions.")
	}

	tm := func(t time.Time) *tspb.Timestamp {
		if t.IsZero() {
			return nil
		}
		return tspb.New(t)
	}

	rmt := make([]*tpb.Remittance, len(res.SearchRemits))
	for i, t := range res.SearchRemits {
		rmt[i] = &tpb.Remittance{
			ControlNumber: t.ControlNo,
			RemitPartner:  t.RemitPartner,
			RemitType:     t.RemitType,
			GrossAmount: &tpb.Amount{
				Amount:   t.RemitAmount.Number(),
				Currency: t.RemitAmount.CurrencyCode(),
			},
			RemitAmount: &tpb.Amount{
				Amount:   t.DisburseAmount.Number(),
				Currency: t.DisburseAmount.CurrencyCode(),
			},
			Remitter: &tpb.Contact{
				FirstName:  t.Remitter.FirstName,
				MiddleName: t.Remitter.MiddleName,
				LastName:   t.Remitter.LastName,
				Email:      t.Remitter.Email,
				Address: &tpb.Address{
					Address1:   t.Remitter.Address.Address1,
					Address2:   t.Remitter.Address.Address2,
					City:       t.Remitter.Address.City,
					State:      t.Remitter.Address.State,
					Province:   t.Remitter.Address.Province,
					PostalCode: t.Remitter.Address.PostalCode,
					Country:    t.Remitter.Address.Country,
					Zone:       t.Remitter.Address.Zone,
				},
				Phone: &ppb.PhoneNumber{
					CountryCode: t.Remitter.Phone.CtyCode,
					Number:      t.Remitter.Phone.Number,
				},
				Mobile: &ppb.PhoneNumber{
					CountryCode: t.Remitter.Mobile.CtyCode,
					Number:      t.Remitter.Mobile.Number,
				},
			},
			Receiver: &tpb.Contact{
				FirstName:  t.Receiver.FirstName,
				MiddleName: t.Receiver.MiddleName,
				LastName:   t.Receiver.LastName,
				Email:      t.Receiver.Email,
				Address: &tpb.Address{
					Address1:   t.Receiver.Address.Address1,
					Address2:   t.Receiver.Address.Address2,
					City:       t.Receiver.Address.City,
					State:      t.Receiver.Address.State,
					Province:   t.Receiver.Address.Province,
					PostalCode: t.Receiver.Address.PostalCode,
					Country:    t.Receiver.Address.Country,
					Zone:       t.Receiver.Address.Zone,
				},
				Phone: &ppb.PhoneNumber{
					CountryCode: t.Receiver.Phone.CtyCode,
					Number:      t.Receiver.Phone.Number,
				},
				Mobile: &ppb.PhoneNumber{
					CountryCode: t.Receiver.Mobile.CtyCode,
					Number:      t.Receiver.Mobile.Number,
				},
			},
			TransactionStagedTime:    tm(t.TxnStagedTime),
			TransactionCompletedTime: tm(t.TxnCompletedTime),
		}
	}

	return &tpb.ListRemitResponse{
		Next:        req.Offset + int32(len(rmt)),
		Remittances: rmt,
		Total:       int32(res.Total),
	}, nil
}
