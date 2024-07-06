package terminal

import (
	"fmt"
	"strings"

	"brank.as/petnet/api/core"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/jbub/banking/swift"
)

type ContactInclude int

const (
	withProvince ContactInclude = iota
	withState
	withZone
	withMobile
	withPhoneCtryCode
	withNameOnly
	withAddrCtryOnly
	withoutAddress
)

type IDInclude int

const (
	withCtry IDInclude = iota
	withExp
	withIss
	withCity
	withStates
	withIstate
)

type EmployInclude int

const (
	withOccID EmployInclude = iota
	withPosLvl
	withEmpl
)

type TxnInclude int

const (
	withSrcCtry TxnInclude = iota
	withDestCtry
)

type AgentInclude int

const (
	withDevID AgentInclude = iota
)

func valAddr(r *tpb.Address, withProvince, withState, withZone, countryOnly bool) validation.Rule {
	return validation.By(func(interface{}) error {
		if r == nil {
			return nil
		}
		return validation.ValidateStruct(r,
			validation.Field(&r.Address1, validation.When(!countryOnly, validation.Required)),
			validation.Field(&r.Address2),
			validation.Field(&r.City, validation.When(!countryOnly, validation.Required)),
			validation.Field(&r.Country, validation.Required, is.CountryCode2),
			validation.Field(&r.PostalCode, validation.When(!countryOnly, validation.Required, is.Alphanumeric)),
			validation.Field(&r.Province, validation.When(!countryOnly && withProvince, validation.Required)),
			validation.Field(&r.State, validation.When(!countryOnly && withState, validation.Required)),
		)
	})
}

func valPhone(r *ppb.PhoneNumber, withCtryCode bool) validation.Rule {
	return validation.By(func(interface{}) error {
		if r == nil {
			return nil
		}
		return validation.ValidateStruct(r,
			validation.Field(&r.CountryCode, validation.When(withCtryCode, validation.Required, is.Int)),
			validation.Field(&r.Number, validation.Required, validation.By(func(interface{}) error {
				// Allow leading zeros
				return validation.Validate(strings.TrimLeft(r.GetNumber(), "0"), is.Int)
			})),
		)
	})
}

func validateAccount(r *tpb.BankAccount) validation.Rule {
	return validation.By(func(interface{}) error {
		if r == nil {
			return nil
		}
		return validation.ValidateStruct(r,
			validation.Field(&r.BIC, validation.Required,
				validation.By(func(interface{}) error {
					_, err := swift.Parse(r.GetBIC())
					if err != nil {
						return fmt.Errorf("invalid BIC. Use assigned SWIFT BIC for the bank")
					}
					return nil
				}),
			),
			validation.Field(&r.AccountNumber, validation.Required, is.Int),
			validation.Field(&r.AccountSuffix, validation.Required),
		)
	})
}

func valContact(r *tpb.Contact, incl ...ContactInclude) validation.Rule {
	pr, st, zo, pc, mo, no, co := false, false, false, false, false, false, false
	for _, v := range incl {
		switch ContactInclude(v) {
		case withProvince:
			pr = true
		case withState:
			st = true
		case withZone:
			zo = true
		case withPhoneCtryCode:
			pc = true
		case withMobile:
			mo = true
		case withNameOnly:
			no = true
		case withAddrCtryOnly:
			co = true
		}
	}
	return validation.By(func(interface{}) error {
		if r == nil {
			return nil
		}
		return validation.ValidateStruct(r,
			validation.Field(&r.FirstName, validation.Required),
			validation.Field(&r.LastName, validation.Required),
			validation.Field(&r.MiddleName),
			validation.Field(&r.Address, validation.When(!no || co, validation.Required, valAddr(r.Address, pr, st, zo, co))),
			validation.Field(&r.Phone, validation.When(!no, validation.Required, valPhone(r.Phone, pc))),
			validation.Field(&r.Mobile, validation.When(mo && !no, validation.Required, valPhone(r.Mobile, pc))),
		)
	})
}

func valID(r *ppb.Identification, incl ...IDInclude) validation.Rule {
	ctry, exp, iss, city := false, false, false, false
	for _, v := range incl {
		switch IDInclude(v) {
		case withCtry:
			ctry = true
		case withExp:
			exp = true
		case withIss:
			iss = true
		case withCity:
			city = true
		}
	}
	return validation.By(func(interface{}) error {
		if r == nil {
			return nil
		}
		return validation.ValidateStruct(r,
			validation.Field(&r.Type, validation.Required),
			validation.Field(&r.Number, validation.Required),
			validation.Field(&r.Country, is.CountryCode2, validation.Required.When(ctry)),
			validation.Field(&r.City, validation.Required.When(city)),
			validation.Field(&r.Expiration, validation.When(exp, validation.Required, validateDate(r.GetExpiration()))),
			validation.Field(&r.Issued, validation.When(iss, validation.Required, validateDate(r.GetIssued()))),
		)
	})
}

func valAgent(r *tpb.Agent, incl ...AgentInclude) validation.Rule {
	di := false
	for _, v := range incl {
		switch AgentInclude(v) {
		case withDevID:
			di = true
		}
	}
	return validation.By(func(interface{}) error {
		if r == nil {
			return nil
		}
		return validation.ValidateStruct(r,
			validation.Field(&r.UserID, validation.Required),
			validation.Field(&r.IPAddress, is.IP, validation.Required),
			validation.Field(&r.DeviceID, validation.Required.When(di)),
		)
	})
}

func valTxn(r *tpb.Transaction, incl ...TxnInclude) validation.Rule {
	oc, dc := false, false
	for _, v := range incl {
		switch TxnInclude(v) {
		case withSrcCtry:
			oc = true
		case withDestCtry:
			dc = true
		}
	}
	return validation.By(func(interface{}) error {
		if r == nil {
			return nil
		}
		return validation.ValidateStruct(r,
			validation.Field(&r.SourceCountry, validation.When(oc, is.CountryCode2)),
			validation.Field(&r.DestinationCountry, validation.When(dc, is.CountryCode2)),
		)
	})
}

func validateDate(r *ppb.Date) validation.Rule {
	return validation.By(func(interface{}) error {
		return (core.ToDate(r)).Validate()
	})
}

func valEmploy(r *tpb.Employment, incl ...EmployInclude) validation.Rule {
	oi, pl, ep := false, false, false
	for _, v := range incl {
		switch EmployInclude(v) {
		case withOccID:
			oi = true
		case withPosLvl:
			pl = true
		case withEmpl:
			ep = true
		}
	}
	return validation.By(func(interface{}) error {
		if r == nil {
			return nil
		}
		required := validation.Required
		return validation.ValidateStruct(r,
			validation.Field(&r.OccupationID, validation.When(oi, is.Digit, required)),
			validation.Field(&r.Occupation, is.ASCII, validation.Required),
			validation.Field(&r.PositionLevel, validation.When(pl, is.ASCII, required)),
			validation.Field(&r.Employer, validation.When(ep, is.ASCII, required)),
		)
	})
}

func ToID(d *ppb.Identification, ctry string) *core.Identification {
	if d == nil {
		return nil
	}
	if ctry == "" {
		ctry = d.GetCountry()
	}
	id := &core.Identification{
		IDType:  d.Type,
		Number:  d.Number,
		Country: ctry,
		City:    d.City,
		Issued:  core.Date{},
		Expiry:  core.ToDate(d.GetExpiration()),
	}
	if d.GetIssued() != nil {
		id.Issued = *core.ToDate(d.Issued)
	}
	return id
}
