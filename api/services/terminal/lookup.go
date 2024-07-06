package terminal

import (
	"context"
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/structpb"
	tspb "google.golang.org/protobuf/types/known/timestamppb"

	"brank.as/petnet/api/core"
	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/core/static"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/util"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) LookupRemit(ctx context.Context, req *tpb.LookupRemitRequest) (*tpb.LookupRemitResponse, error) {
	log := logging.FromContext(ctx)

	pn := req.GetRemitPartner()

	orgType := phmw.GetOrgInfo(ctx)

	if orgType == phmw.Provider {
		pn = static.PerahubRemit
	}

	// todo(robin): if we will be using other countries change country to dynamic
	if !static.PartnerExists(pn, "PH") {
		log.Error("partner doesn't exist")
		return nil, util.HandleServiceErr(status.Error(codes.NotFound, coreerror.MsgPartnerDoesntExist))
	}

	r, err := s.validators[pn].LookupRemitValidate(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	rmt, err := s.remit.SearchRemit(ctx, *r, pn)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	var tax, chg map[string]*tpb.Amount
	if len(rmt.Taxes) > 0 {
		tax = make(map[string]*tpb.Amount, len(rmt.Taxes))
		for k, v := range rmt.Taxes {
			tax[k] = &tpb.Amount{
				Amount:   v.Number(),
				Currency: v.CurrencyCode(),
			}
		}
	}
	if len(rmt.Charges) > 0 {
		chg = make(map[string]*tpb.Amount, len(rmt.Charges))
		for k, v := range rmt.Charges {
			chg[k] = &tpb.Amount{
				Amount:   v.Number(),
				Currency: v.CurrencyCode(),
			}
		}
	}

	tm := func(t time.Time) *tspb.Timestamp {
		if t.IsZero() {
			return nil
		}
		return tspb.New(t)
	}

	if pn != "WU" {
		l := &structpb.Value{}
		protojson.Unmarshal(rmt.OtherInfo, l)
		return &tpb.LookupRemitResponse{
			ControlNumber: rmt.ControlNo,
			Status:        rmt.Status,
			RemitAmount: &tpb.Amount{
				Amount:   rmt.RemitAmount.Number(),
				Currency: rmt.RemitAmount.CurrencyCode(),
			},
			DisburseAmount: &tpb.Amount{
				Amount:   rmt.DisburseAmount.Number(),
				Currency: rmt.DisburseAmount.CurrencyCode(),
			},
			OtherInfo: l,
		}, nil
	} else {
		return &tpb.LookupRemitResponse{
			ControlNumber: rmt.ControlNo,
			Remitter: &tpb.Contact{
				FirstName:  rmt.Remitter.FirstName,
				MiddleName: rmt.Remitter.MiddleName,
				LastName:   rmt.Remitter.LastName,
				Address: &tpb.Address{
					Address1:   rmt.Remitter.Address.Address1,
					Address2:   rmt.Remitter.Address.Address2,
					City:       rmt.Remitter.Address.City,
					State:      rmt.Remitter.Address.State,
					PostalCode: rmt.Remitter.Address.PostalCode,
					Country:    rmt.Remitter.Address.Country,
				},
				Phone: &ppb.PhoneNumber{
					CountryCode: rmt.Remitter.Phone.CtyCode,
					Number:      rmt.Remitter.Phone.Number,
				},
				Mobile: &ppb.PhoneNumber{
					CountryCode: rmt.Remitter.Mobile.CtyCode,
					Number:      rmt.Remitter.Mobile.Number,
				},
			},
			Receiver: &tpb.Contact{
				FirstName:  rmt.Receiver.FirstName,
				MiddleName: rmt.Receiver.MiddleName,
				LastName:   rmt.Receiver.LastName,
				Address: &tpb.Address{
					Address1:   rmt.Receiver.Address.Address1,
					Address2:   rmt.Receiver.Address.Address2,
					City:       rmt.Receiver.Address.City,
					State:      rmt.Receiver.Address.State,
					PostalCode: rmt.Receiver.Address.PostalCode,
					Country:    rmt.Receiver.Address.Country,
				},
				Phone: &ppb.PhoneNumber{
					CountryCode: rmt.Receiver.Phone.CtyCode,
					Number:      rmt.Receiver.Phone.Number,
				},
				Mobile: &ppb.PhoneNumber{
					CountryCode: rmt.Receiver.Mobile.CtyCode,
					Number:      rmt.Receiver.Mobile.Number,
				},
			},
			SourceCountry:    rmt.SentCountry,
			DestinationCity:  rmt.DestCity,
			DestinationState: rmt.DestState,
			Status:           rmt.Status,
			RemitAmount: &tpb.Amount{
				Amount:   rmt.RemitAmount.Number(),
				Currency: rmt.RemitAmount.CurrencyCode(),
			},
			Taxes: tax,
			TotalTax: &tpb.Amount{
				Amount:   rmt.Tax.CurrencyCode(),
				Currency: rmt.Tax.Number(),
			},
			Charges: chg,
			TotalCharges: &tpb.Amount{
				Amount:   rmt.Charge.Number(),
				Currency: rmt.Charge.CurrencyCode(),
			},
			DisburseAmount: &tpb.Amount{
				Amount:   rmt.DisburseAmount.Number(),
				Currency: rmt.DisburseAmount.CurrencyCode(),
			},
			TransactionStagedTime:    tm(rmt.TxnStagedTime),
			TransactionCompletedTime: tm(rmt.TxnCompletedTime),
		}, nil
	}
}

func (s *WUVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.PayoutCurrency, required, is.CurrencyCode),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric, validation.Length(10, 10)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		DestCurrency: req.PayoutCurrency,
		ControlNo:    req.ControlNumber,
	}
	return r, nil
}

func (s *TFVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
	}
	return r, nil
}

func (s *ICVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required, is.Digit),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
	}
	return r, nil
}

func (s *UNTVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required, is.Digit),
		validation.Field(&req.DeviceID, required, is.Digit),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
		DeviceID:     req.DeviceID,
	}
	return r, nil
}

func (s *RIAVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required, is.Digit),
		validation.Field(&req.DeviceID, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
		DeviceID:     req.DeviceID,
	}
	return r, nil
}

func (s *RMVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required, is.Digit),
		validation.Field(&req.DeviceID, required),
		validation.Field(&req.PayoutCurrency, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
		DeviceID:     req.DeviceID,
		DestCurrency: req.PayoutCurrency,
	}
	return r, nil
}

func (s *USSCVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required, is.Digit),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
	}
	return r, nil
}

func (s *JPRVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required, is.Digit),
		validation.Field(&req.DeviceID, required),
		validation.Field(&req.PayoutCurrency, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
		DeviceID:     req.DeviceID,
		DestCurrency: req.PayoutCurrency,
	}
	return r, nil
}

func (s *IRVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required, is.Digit),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
	}
	return r, nil
}

func (s *MBVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required, is.Digit),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
		DeviceID:     req.DeviceID,
	}
	return r, nil
}

func (s *BPIVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required, is.Digit),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
		DeviceID:     req.DeviceID,
	}
	return r, nil
}

func (s *WISEVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	return nil, fmt.Errorf("service not available for Wise")
}

func (s *CEBVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.UserID, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
		DeviceID:     req.DeviceID,
	}
	return r, nil
}

func (s *CEBIVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.UserID, required),
		validation.Field(&req.DeviceID, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
		DeviceID:     req.DeviceID,
	}
	return r, nil
}

func (s *AYAVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.UserID, required),
		validation.Field(&req.DeviceID, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
		DeviceID:     req.DeviceID,
	}
	return r, nil
}

func (s *IEVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required, is.Digit),
		validation.Field(&req.DeviceID, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
		DeviceID:     req.DeviceID,
	}
	return r, nil
}

func (s *PHUBVal) LookupRemitValidate(ctx context.Context, req *tpb.LookupRemitRequest) (*core.SearchRemit, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required),
		validation.Field(&req.UserID, required, is.Digit),
		validation.Field(&req.DeviceID, required),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	r := &core.SearchRemit{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		RemitPartner: req.RemitPartner,
		ControlNo:    req.ControlNumber,
		PtnrUserID:   req.UserID,
		DeviceID:     req.DeviceID,
	}
	return r, nil
}
