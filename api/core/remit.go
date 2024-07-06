package core

import (
	"fmt"
	"time"

	"github.com/bojanz/currency"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"

	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	"brank.as/petnet/serviceutil/logging"
)

type (
	ContactInfo tpb.Contact
	Phone       ppb.PhoneNumber
	SendAmount  tpb.SendAmount
	Taxes       tpb.Taxes
)

const (
	CreateRemType   = "CREATE"
	DisburseRemType = "DISBURSE"
)

type Contact struct {
	FirstName  string
	MiddleName string
	LastName   string
	Email      string
	Address    Address
	Phone      PhoneNumber
	Mobile     PhoneNumber
}

type Date struct {
	Year  string
	Month string
	Day   string
}

type Remittance struct {
	DsaID             string
	DsaOrderID        string
	UserID            string
	TransactionID     string
	CustomerTxnID     string
	GeneratedRefNo    string
	RemcoAltControlNo string

	ControlNo         string
	MyWUNumber        string
	RemitPartner      string
	SendRemitType     SendRemitType
	DisburseRemitType DisburseRemitType
	SendReason        string
	TxnPurpose        string

	Remitter UserKYC
	Receiver UserKYC
	Business Business
	Agent    Agent

	GrossTotal   currency.Minor
	SourceAmount currency.Minor
	DestAmount   currency.Minor
	Tax          currency.Minor
	Charge       currency.Minor
	TargetDest   bool

	DestState   string
	DestCity    string
	DestAccount Account

	TransactionDetails TransactionDetails

	Promo string

	Message string

	FixedAmountFlag string

	IntlPrtCode string
}

type TransactionDetails struct {
	SrcCtry    string
	DestCtry   string
	IsDomestic int
}

type UserKYC struct {
	PartnerMemberID     string
	FName               string
	MdName              string
	LName               string
	Gender              string
	KYCVerified         bool
	Address             Address
	Phone               PhoneNumber
	Mobile              PhoneNumber
	SourceFunds         string
	Employment          Employment
	ReceiverRelation    string
	SendingReasonID     string
	ProofOfAddress      string
	PrimaryID           Identification
	AlternateID         []Identification
	Email               string
	BirthDate           Date
	BirthCountry        string
	BirthPlace          string
	Nationality         string
	CurrentAddress      Address
	PermanentAddress    Address
	RecipientID         string
	AccountHolderName   string
	SourceAccountNumber string
}

type Agent struct {
	UserID       int
	IPAddress    string
	DeviceID     string
	LocationID   string
	LocationName string
	LocationCode string
	AgentID      string
	AgentCode    string
	OutletCode   string
}

type Employment struct {
	Employer      string
	OccupationID  string
	Occupation    string
	PositionLevel string
}

type Account struct {
	BIC     string
	AcctNo  string
	AcctSfx string
}

type PhoneNumber struct {
	CtyCode string
	Number  string
}

type Business struct {
	Name      string
	Account   string
	ControlNo string
	Country   string
}

type RemitResponse struct {
	PrincipalAmount  currency.Minor
	RemitAmount      currency.Minor
	Taxes            map[string]currency.Minor
	Tax              currency.Minor
	Charges          map[string]currency.Minor
	Charge           currency.Minor
	GrossTotal       currency.Minor
	PromoDescription string
	PromoMessage     string
	TransactionID    string
}

type ProcessRemit struct {
	TransactionID string
	AuthSource    string
	AuthCode      string
	ControlNumber string
	Processed     time.Time
	RemitCache    storage.RemitCache
}

type Identification struct {
	IDType  string
	Number  string
	Country string
	State   string
	City    string
	Istate  string
	Issued  Date
	Expiry  *Date
}

type AddrSource interface {
	GetAddress1() string
	GetAddress2() string
	GetCity() string
	GetState() string
	GetPostalCode() string
	GetCountry() string
	GetZone() string
}

type ProvinceSrc interface {
	GetProvince() string
	GetRegion() string
}

type Address struct {
	Address1   string
	Address2   string
	City       string
	State      string
	Province   string
	Region     string
	PostalCode string
	Country    string
	Zone       string
}

func ToAddr(s AddrSource) Address {
	a := Address{
		Address1:   s.GetAddress1(),
		Address2:   s.GetAddress2(),
		City:       s.GetCity(),
		State:      s.GetState(),
		PostalCode: s.GetPostalCode(),
		Country:    s.GetCountry(),
		Zone:       s.GetZone(),
	}
	if p, ok := s.(ProvinceSrc); ok {
		a.Province = p.GetProvince()
		a.Region = p.GetRegion()
	}
	return a
}

func ToPhone(p interface {
	GetCountryCode() string
	GetNumber() string
},
) PhoneNumber {
	return PhoneNumber{
		CtyCode: p.GetCountryCode(),
		Number:  p.GetNumber(),
	}
}

func ToDate(d *ppb.Date) *Date {
	if d == nil {
		return nil
	}
	return &Date{
		Year:  d.GetYear(),
		Month: d.GetMonth(),
		Day:   d.GetDay(),
	}
}

func MustAmount(amt, cur string) currency.Amount {
	c, _ := currency.NewAmount(amt, cur)
	return c
}

func MustMinor(amt, cur string) currency.Minor {
	c, _ := currency.NewMinor(amt, cur)
	return c
}

func (d *Date) Validate() error {
	return validation.Validate(
		fmt.Sprintf("%s/%s/%s", d.Year, d.Month, d.Day), validation.Date("2006/1/2"),
		validation.By(func(val interface{}) error {
			t, err := time.Parse("2006/1/2", val.(string))
			if err != nil {
				return err
			}
			if t.Before(time.Now().AddDate(-120, 0, 0)) {
				return fmt.Errorf("invalid date")
			}
			return nil
		}),
	)
}

func (d *Date) String() string {
	if d == nil {
		return ""
	}
	t, err := time.Parse("2006/1/2", fmt.Sprintf("%s/%s/%s", d.Year, d.Month, d.Day))
	if err != nil {
		logging.WithError(err, logging.NewLogger(nil)).Error("context")
		return ""
	}
	return t.Format("02012006")
}

func (i *Identification) Validate() error {
	return validation.ValidateStruct(&i,
		validation.Field(&i.Country, validation.Required, is.CountryCode2),
		validation.Field(&i.IDType, validation.Required),
		validation.Field(&i.Number, validation.Required),
		validation.Field(&i.Issued, validation.Required, validation.By(func(interface{}) error {
			return (&Date{
				Year:  i.Issued.Year,
				Month: i.Issued.Month,
				Day:   i.Issued.Day,
			}).Validate()
		})),
		validation.Field(&i.Expiry, validation.By(func(interface{}) error {
			return (&Date{
				Year:  i.Expiry.Year,
				Month: i.Expiry.Month,
				Day:   i.Expiry.Day,
			}).Validate()
		})),
	)
}

func (p *Phone) Validate() error {
	return validation.ValidateStruct(&p,
		validation.Field(&p.CountryCode, validation.Required, is.Int),
		validation.Field(&p.Number, validation.Required, is.Int),
	)
}

func (a *Address) Validate() error {
	return validation.ValidateStruct(&a,
		validation.Field(&a.Address1, validation.Required),
		validation.Field(&a.Address2),
		validation.Field(&a.City, validation.Required),
		validation.Field(&a.Country, validation.Required, is.CountryCode2),
		validation.Field(&a.PostalCode, validation.Required, is.Int),
		validation.Field(&a.State, validation.Required),
	)
}

func (c *ContactInfo) Validate() error {
	return validation.ValidateStruct(&c,
		validation.Field(&c.FirstName, validation.Required),
		validation.Field(&c.MiddleName, validation.Required),
		validation.Field(&c.LastName, validation.Required),
		validation.Field(&c.Address, validation.Required),
		validation.Field(&c.Phone, validation.Required, is.Int),
		validation.Field(&c.Mobile, validation.Required, is.Int),
	)
}

func (r *UserKYC) Validate() error {
	return validation.ValidateStruct(&r,
		validation.Field(&r.Email, validation.Required, is.Email),
		validation.Field(&r.Phone, validation.Required, is.Int),
		validation.Field(&r.ReceiverRelation, validation.Required),
		validation.Field(&r.SourceFunds, validation.Required),
	)
}

func (s *SendAmount) Validate() error {
	return validation.ValidateStruct(&s,
		validation.Field(&s.Amount, validation.Required, is.Int),
		validation.Field(&s.DestinationAmount, is.Int),
		validation.Field(&s.DestinationCurrency, validation.Required, is.CurrencyCode),
		validation.Field(&s.SourceCurrency, validation.Required, is.CurrencyCode),
	)
}
