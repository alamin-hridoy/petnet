package terminal

import (
	"context"
	"fmt"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
	phmw "brank.as/petnet/api/perahub-middleware"
	"brank.as/petnet/api/util"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	"brank.as/petnet/serviceutil/logging"
)

const defaultCountry = "ZZ"

func (s *Svc) DisburseRemit(ctx context.Context, req *tpb.DisburseRemitRequest) (*tpb.DisburseRemitResponse, error) {
	log := logging.FromContext(ctx)

	pn := req.GetRemitPartner()
	orgType := phmw.GetOrgInfo(ctx)

	if orgType == phmw.Provider {
		pn = static.PerahubRemit
	}
	_, err := s.lk.DisburseRemitType(ctx, pn, req.GetRemitType(), false)
	if err != nil {
		log.Error(err)
		return nil, util.HandleServiceErr(err)
	}

	r, err := s.validators[pn].DisburseRemitValidate(ctx, req)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	rmt, err := s.remit.StageDisburseRemit(ctx, *r, pn)
	if err != nil {
		return nil, util.HandleServiceErr(err)
	}

	return &tpb.DisburseRemitResponse{TransactionID: rmt.TransactionID}, nil
}

func (s *WUVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.DisburseCurrency, required, is.CurrencyCode),
		validation.Field(&req.OrderID, required),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withMobile, withPhoneCtryCode, withState)),
				validation.Field(&r.PartnerMemberID, is.Alphanumeric),
				validation.Field(&r.Employment, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.Nationality, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.SendingReason),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification, withCtry)),
				validation.Field(&r.AlternateID, validation.Each(
					validation.By(func(id interface{}) error {
						i, _ := id.(*ppb.Identification)
						return valID(i, withCtry).Validate(nil)
					}),
				)),
				validation.Field(&r.Email, required, is.EmailFormat),
			)
		})),
		validation.Field(&req.Agent, required, validation.By(func(interface{}) error {
			r := req.Agent
			return validation.ValidateStruct(r,
				validation.Field(&r.UserID, required))
		})),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	rm := &core.Remittance{
		DsaID:  phmw.GetDSA(ctx),
		UserID: phmw.GetUserID(ctx),
		Agent: core.Agent{
			UserID: int(req.Agent.UserID),
		},
		DsaOrderID:   req.OrderID,
		ControlNo:    req.ControlNumber,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		SendReason:   req.Receiver.SendingReason,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		DestAmount:   core.MustMinor("0", req.DisburseCurrency),
		Receiver: core.UserKYC{
			PartnerMemberID: req.Receiver.PartnerMemberID,
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.Receiver.ContactInfo.MiddleName,
			LName:           req.Receiver.ContactInfo.LastName,
			Gender:          req.Receiver.Gender.String(),
			Address:         core.ToAddr(req.Receiver.ContactInfo.Address),
			Phone:           core.ToPhone(req.Receiver.ContactInfo.Phone),
			Mobile:          core.ToPhone(req.Receiver.ContactInfo.Mobile),
			SourceFunds:     req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation:    req.Receiver.Employment.Occupation,
				PositionLevel: req.Receiver.Employment.PositionLevel,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, ""),
			AlternateID: func() []core.Identification {
				ids := make([]core.Identification, len(req.Receiver.AlternateID))
				for i, id := range req.Receiver.AlternateID {
					ids[i] = *ToID(id, "")
				}
				return ids
			}(),
			Email:        req.Receiver.Email,
			BirthDate:    *core.ToDate(req.Receiver.Birthdate),
			BirthCountry: req.Receiver.BirthCountry,
			Nationality:  req.Receiver.Nationality,
		},
	}
	return rm, nil
}

func (s *IRVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	bctry, err := s.q.FindCountryByAlpha(req.Receiver.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	destcntry := ""
	if req.Transaction != nil && req.Transaction.DestinationCountry != "" {
		destctry, err := s.q.FindCountryByAlpha(req.Transaction.DestinationCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid destination country")
		}
		destcntry = destctry.Name.Common
	}

	ctryName := defaultCountry
	if req.Transaction != nil && req.Transaction.SourceCountry != "" {
		srcCtry, err := s.q.FindCountryByAlpha(req.Transaction.SourceCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid originating country")
		}
		ctryName = srcCtry.Name.Common
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, ""),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     bctry.Name.Common,
			BirthPlace:       req.Receiver.BirthPlace,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}
	if req.Transaction != nil {
		rm.TransactionDetails = core.TransactionDetails{
			SrcCtry:    ctryName,
			DestCtry:   destcntry,
			IsDomestic: perahub.IsDomestic(req.Transaction.SourceCountry, req.Transaction.DestinationCountry),
		}
	}
	return rm, nil
}

func (s *TFVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.Gender, required),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.SendingReasonID, required, is.Digit),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.ProofOfAddress, required),
				validation.Field(&r.KYCVerified, required),
				validation.Field(&r.Identification, required, valID(r.Identification, withCtry, withExp)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment, withOccID)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	bctry, err := s.q.FindCountryByAlpha(req.Receiver.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	idctry, err := s.q.FindCountryByAlpha(req.Receiver.Identification.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid id country")
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				OccupationID: req.Receiver.Employment.OccupationID,
				Occupation:   req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, idctry.Name.Common),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     bctry.Name.Common,
			BirthPlace:       req.Receiver.BirthPlace,
			Gender:           req.Receiver.Gender.String(),
			KYCVerified:      req.Receiver.KYCVerified == tpb.Bool_True,
			SendingReasonID:  req.Receiver.SendingReasonID,
			ProofOfAddress:   req.Receiver.ProofOfAddress.String(),
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}
	return rm, nil
}

func (s *RMVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.DisburseCurrency, required, is.CurrencyCode),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				// todo(robin): validate province and state, also move all partners to using province instead of state
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withState, withZone)),
				validation.Field(&r.PartnerMemberID, required, is.Digit),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.Gender, required),
				validation.Field(&r.Nationality, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification, withCtry, withExp)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	ctryName := defaultCountry
	if req.Transaction != nil && req.Transaction.SourceCountry != "" {
		srcCtry, err := s.q.FindCountryByAlpha(req.Transaction.SourceCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid originating country")
		}
		ctryName = srcCtry.Name.Common
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				State:      req.Receiver.ContactInfo.Address.State,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, req.Receiver.Identification.Country),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     req.Receiver.BirthCountry,
			BirthPlace:       req.Receiver.BirthPlace,
			Gender:           req.Receiver.Gender.String(),
			Nationality:      req.Receiver.Nationality,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
			DeviceID:  req.Agent.DeviceID,
		},
	}
	if req.Transaction != nil {
		rm.TransactionDetails = core.TransactionDetails{
			SrcCtry:    ctryName,
			DestCtry:   req.Transaction.DestinationCountry,
			IsDomestic: perahub.IsDomestic(req.Transaction.SourceCountry, req.Transaction.DestinationCountry),
		}
	}
	return rm, nil
}

func (s *RIAVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required, is.Digit),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.Gender, required),
				validation.Field(&r.Nationality, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification, withCtry, withExp)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent, withDevID)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}

	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, req.Receiver.Identification.Country),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     req.Receiver.BirthCountry,
			BirthPlace:       req.Receiver.BirthPlace,
			Gender:           req.Receiver.Gender.String(),
			Nationality:      req.Receiver.Nationality,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
			DeviceID:  req.Agent.DeviceID,
		},
	}
	return rm, nil
}

func (s *MBVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required, is.Digit),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Remitter, required, valContact(req.Remitter, withNameOnly)),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	ctryName := defaultCountry
	if req.Transaction != nil && req.Transaction.SourceCountry != "" {
		srcCtry, err := s.q.FindCountryByAlpha(req.Transaction.SourceCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid originating country")
		}
		ctryName = srcCtry.Name.Common
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, req.Receiver.Identification.Country),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     req.Receiver.BirthCountry,
			BirthPlace:       req.Receiver.BirthPlace,
			Gender:           req.Receiver.Gender.String(),
		},
		Remitter: core.UserKYC{
			FName:  req.Remitter.GetFirstName(),
			MdName: req.Remitter.GetMiddleName(),
			LName:  req.Remitter.GetLastName(),
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}
	if req.Transaction != nil {
		rm.TransactionDetails = core.TransactionDetails{
			SrcCtry:    ctryName,
			DestCtry:   req.Transaction.DestinationCountry,
			IsDomestic: perahub.IsDomestic(req.Transaction.SourceCountry, req.Transaction.DestinationCountry),
		}
	}
	return rm, nil
}

func (s *BPIVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	bctry, err := s.q.FindCountryByAlpha(req.Receiver.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, ""),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     bctry.Name.Common,
			BirthPlace:       req.Receiver.BirthPlace,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}
	return rm, nil
}

func (s *JPRVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.DisburseCurrency, required, is.CurrencyCode),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	bctry, err := s.q.FindCountryByAlpha(req.Receiver.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, ""),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     bctry.Name.Common,
			BirthPlace:       req.Receiver.BirthPlace,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}
	return rm, nil
}

func (s *USSCVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	bctry, err := s.q.FindCountryByAlpha(req.Receiver.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, ""),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     bctry.Name.Common,
			BirthPlace:       req.Receiver.BirthPlace,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}
	return rm, nil
}

func (s *ICVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	bctry, err := s.q.FindCountryByAlpha(req.Receiver.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, ""),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     bctry.Name.Common,
			BirthPlace:       req.Receiver.BirthPlace,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}
	return rm, nil
}

func (s *WISEVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	return nil, fmt.Errorf("service not available for Wise")
}

func (s *UNTVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withState, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.Nationality, required, is.CountryCode2),
				validation.Field(&r.Gender, required),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification, withCtry, withExp)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent, withDevID)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ctryName := defaultCountry
	if req.Transaction != nil && req.Transaction.SourceCountry != "" {
		srcCtry, err := s.q.FindCountryByAlpha(req.Transaction.SourceCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid originating country")
		}
		ctryName = srcCtry.Name.Common
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		DestAmount:   core.MustMinor("0", "PHP"),
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Nationality:     req.Receiver.Nationality,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    req.Receiver.ContactInfo.Address.Country,
				State:      req.Receiver.ContactInfo.Address.State,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, req.Receiver.Identification.Country),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     req.Receiver.BirthCountry,
			BirthPlace:       req.Receiver.BirthPlace,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
			DeviceID:  req.Agent.DeviceID,
		},
	}
	if req.Transaction != nil {
		rm.TransactionDetails = core.TransactionDetails{
			SrcCtry:    ctryName,
			DestCtry:   req.Transaction.DestinationCountry,
			IsDomestic: perahub.IsDomestic(req.Transaction.SourceCountry, req.Transaction.DestinationCountry),
		}
	}
	return rm, nil
}

func (s *CEBVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ctryName := defaultCountry
	if req.Transaction != nil && req.Transaction.SourceCountry != "" {
		srcCtry, err := s.q.FindCountryByAlpha(req.Transaction.SourceCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid originating country")
		}
		ctryName = srcCtry.Name.Common
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		DestAmount:   core.MustMinor("0", "PHP"),
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    req.Receiver.ContactInfo.Address.Country,
				State:      req.Receiver.ContactInfo.Address.State,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, req.Receiver.Identification.Country),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     req.Receiver.BirthCountry,
			BirthPlace:       req.Receiver.BirthPlace,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}
	if req.Transaction != nil {
		rm.TransactionDetails = core.TransactionDetails{
			SrcCtry:    ctryName,
			DestCtry:   req.Transaction.DestinationCountry,
			IsDomestic: perahub.IsDomestic(req.Transaction.SourceCountry, req.Transaction.DestinationCountry),
		}
	}
	return rm, nil
}

func (s *CEBIVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification, withCtry, withCity, withStates, withIstate)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	dctryName := ""
	if req.Transaction != nil && req.Transaction.DestinationCountry != "" {
		dctry, err := s.q.FindCountryByAlpha(req.Transaction.DestinationCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid destination country")
		}
		dctryName = dctry.Name.Common
	}
	bctry, err := s.q.FindCountryByAlpha(req.Receiver.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	ctryName := defaultCountry
	if req.Transaction != nil && req.Transaction.SourceCountry != "" {
		srcCtry, err := s.q.FindCountryByAlpha(req.Transaction.SourceCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid originating country")
		}
		ctryName = srcCtry.Name.Common
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		DestAmount:   core.MustMinor("0", "PHP"),
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				State:      req.Receiver.ContactInfo.Address.State,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, ""),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     req.Receiver.BirthCountry,
			BirthPlace:       bctry.Name.Common,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}
	if req.Transaction != nil {
		rm.TransactionDetails = core.TransactionDetails{
			SrcCtry:    ctryName,
			DestCtry:   dctryName,
			IsDomestic: perahub.IsDomestic(req.Transaction.SourceCountry, req.Transaction.DestinationCountry),
		}
	}
	return rm, nil
}

func (s *AYAVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	bctry, err := s.q.FindCountryByAlpha(req.Receiver.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	dctryName := ""
	if req.Transaction != nil && req.Transaction.DestinationCountry != "" {
		dctry, err := s.q.FindCountryByAlpha(req.Transaction.DestinationCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid destination country")
		}
		dctryName = dctry.Name.Common
	}
	ctryName := defaultCountry
	if req.Transaction != nil && req.Transaction.SourceCountry != "" {
		srcCtry, err := s.q.FindCountryByAlpha(req.Transaction.SourceCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid originating country")
		}
		ctryName = srcCtry.Name.Common
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, ""),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     bctry.Name.Common,
			BirthPlace:       req.Receiver.BirthPlace,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}
	if req.Transaction != nil {
		rm.TransactionDetails = core.TransactionDetails{
			SrcCtry:    ctryName,
			DestCtry:   dctryName,
			IsDomestic: perahub.IsDomestic(req.Transaction.SourceCountry, req.Transaction.DestinationCountry),
		}
	}
	return rm, nil
}

func (s *IEVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.DisburseCurrency, required, is.CurrencyCode),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	ctry, err := s.q.FindCountryByAlpha(req.Receiver.ContactInfo.Address.Country)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid contactinfo country")
	}
	bctry, err := s.q.FindCountryByAlpha(req.Receiver.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	ctryName := defaultCountry
	if req.Transaction != nil && req.Transaction.SourceCountry != "" {
		srcCtry, err := s.q.FindCountryByAlpha(req.Transaction.SourceCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid originating country")
		}
		ctryName = srcCtry.Name.Common
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		DestAmount:   core.MustMinor("0", req.DisburseCurrency),
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    ctry.Name.Common,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, ""),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     bctry.Name.Common,
			BirthPlace:       req.Receiver.BirthPlace,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}

	if req.Transaction != nil {
		rm.TransactionDetails = core.TransactionDetails{
			SrcCtry:    ctryName,
			DestCtry:   req.Transaction.DestinationCountry,
			IsDomestic: perahub.IsDomestic(req.Transaction.SourceCountry, req.Transaction.DestinationCountry),
		}
	}
	return rm, nil
}

func (s *PHUBVal) DisburseRemitValidate(ctx context.Context, req *tpb.DisburseRemitRequest) (*core.Remittance, error) {
	required := validation.Required
	if err := validation.ValidateStruct(req,
		validation.Field(&req.RemitPartner, required, is.Alphanumeric),
		validation.Field(&req.ControlNumber, required, is.Alphanumeric),
		validation.Field(&req.OrderID, required, is.Digit),
		validation.Field(&req.DisburseCurrency, required, is.CurrencyCode),
		validation.Field(&req.Receiver, required, validation.By(func(interface{}) error {
			r := req.Receiver
			return validation.ValidateStruct(r,
				validation.Field(&r.ContactInfo, required, valContact(r.ContactInfo, withProvince, withZone)),
				validation.Field(&r.PartnerMemberID, required),
				validation.Field(&r.Birthdate, required, validateDate(r.Birthdate)),
				validation.Field(&r.BirthCountry, required, is.CountryCode2),
				validation.Field(&r.ReceiverRelation, required),
				validation.Field(&r.TransactionPurpose, required),
				validation.Field(&r.Identification, required, valID(r.Identification)),
				validation.Field(&r.Employment, required, valEmploy(r.Employment)),
				validation.Field(&r.SourceFunds, required),
			)
		})),
		validation.Field(&req.Agent, required, valAgent(req.Agent)),
	); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	bctry, err := s.q.FindCountryByAlpha(req.Receiver.BirthCountry)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, "invalid birthcountry")
	}
	ctryName := defaultCountry
	if req.Transaction != nil && req.Transaction.SourceCountry != "" {
		srcCtry, err := s.q.FindCountryByAlpha(req.Transaction.SourceCountry)
		if err != nil {
			return nil, status.Error(codes.InvalidArgument, "invalid originating country")
		}
		ctryName = srcCtry.Name.Common
	}
	rm := &core.Remittance{
		DsaID:        phmw.GetDSA(ctx),
		UserID:       phmw.GetUserID(ctx),
		ControlNo:    req.ControlNumber,
		DsaOrderID:   req.OrderID,
		MyWUNumber:   req.Receiver.PartnerMemberID,
		RemitPartner: req.RemitPartner,
		TxnPurpose:   req.Receiver.TransactionPurpose,
		DestAmount:   core.MustMinor("0", req.DisburseCurrency),
		GrossTotal:   core.MustMinor("0", req.DisburseCurrency),
		Tax:          core.MustMinor("0", req.DisburseCurrency),
		Charge:       core.MustMinor("0", req.DisburseCurrency),
		Receiver: core.UserKYC{
			FName:           req.Receiver.ContactInfo.FirstName,
			MdName:          req.GetReceiver().GetContactInfo().GetMiddleName(),
			LName:           req.Receiver.ContactInfo.LastName,
			PartnerMemberID: req.Receiver.PartnerMemberID,
			Address: core.Address{
				Address1:   req.Receiver.ContactInfo.Address.Address1,
				Address2:   req.GetReceiver().GetContactInfo().GetAddress().GetAddress2(),
				City:       req.Receiver.ContactInfo.Address.City,
				Province:   req.Receiver.ContactInfo.Address.Province,
				PostalCode: req.Receiver.ContactInfo.Address.PostalCode,
				Country:    req.Receiver.ContactInfo.Address.Country,
				Zone:       req.Receiver.ContactInfo.Address.Zone,
			},
			Phone: core.PhoneNumber{
				Number: req.Receiver.ContactInfo.Phone.Number,
			},
			SourceFunds: req.Receiver.SourceFunds,
			Employment: core.Employment{
				Occupation: req.Receiver.Employment.Occupation,
			},
			ReceiverRelation: req.Receiver.ReceiverRelation,
			PrimaryID:        *ToID(req.Receiver.Identification, ""),
			BirthDate:        *core.ToDate(req.Receiver.Birthdate),
			BirthCountry:     bctry.Name.Common,
			BirthPlace:       req.Receiver.BirthPlace,
		},
		Agent: core.Agent{
			UserID:    int(req.Agent.UserID),
			IPAddress: req.Agent.IPAddress,
		},
	}

	if req.Transaction != nil {
		rm.TransactionDetails = core.TransactionDetails{
			SrcCtry:    ctryName,
			DestCtry:   req.Transaction.DestinationCountry,
			IsDomestic: perahub.IsDomestic(req.Transaction.SourceCountry, req.Transaction.DestinationCountry),
		}
	}
	return rm, nil
}
