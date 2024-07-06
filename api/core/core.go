package core

import (
	"time"

	rbupb "brank.as/rbac/gunk/v1/user"
	"github.com/bojanz/currency"
)

const LocationID = "371"

type User struct {
	FrgnRefNo     string
	CustNo        string
	LastName      string
	FirstName     string
	Birthdate     string
	Nationality   string
	Address       string
	Occupation    string
	Employer      string
	ValidIdnt     string
	WUCardNo      string
	DebitCardNo   string
	LoyaltyCardNo string
}

type Remco struct {
	Country string
	Code    string
	Name    string
	// Types is a map of Type Name to SendRemitType Details
	SendTypes     map[string]SendRemitType
	DisburseTypes map[string]DisburseRemitType
}

type SendRemitType struct {
	Code        string
	Description string
	Receiver    bool
	BankAccount bool
	Business    bool
}

type DisburseRemitType struct {
	Code        string
	Description string
}

type FeeInquiryReq struct {
	RemitPartner      string
	RemitType         SendRemitType
	PrincipalAmount   currency.Minor
	DestinationAmount bool
	DestCountry       string
	DestCurrency      string
	Promo             string
	Message           string
}

type FilterList struct {
	From           time.Time
	Until          time.Time
	Limit          int
	Offset         int
	SortOrder      string
	SortByColumn   string
	ControlNumbers []string
	ExcludePartner string
	ExcludeType    string
}

type UserFilterList struct {
	Name         string
	SortBy       int
	SortByColumn int
	Status       []rbupb.Status
}

type RegisterUserReq struct {
	Email     string
	FirstName string
	LastName  string
	BirthDate string
	CpCtryID  string
	ContactNo string
	TpCtryID  string
	TpArCode  string
	CrtyAdID  string
	PAdd      string
	CAdd      string
	UserID    string
	SOFID     string
	Tin       string
	TpNo      string
	AgentCode string
}

type RegisterUserResp struct {
	Code    int
	Message string
	Result  RUResult
	RemcoID int
}

type RUResult struct {
	ResultStatus string
	MessageID    int
	LogID        int
	ClientID     int
	ClientNo     string
}

type CreateProfileReq struct {
	Email      string
	Type       string
	FirstName  string
	LastName   string
	BirthDate  string
	Phone      PhoneNumber
	Address    Address
	Occupation string
}

type CreateProfileResp struct{}

type GetProfileReq struct {
	Email string
}

type GetProfileResp struct {
	ID         string
	Type       string
	FirstName  string
	LastName   string
	BirthDate  string
	Phone      string
	Occupation string
	Address    Address
}

type GetUserRequest struct {
	FirstName    string
	LastName     string
	BirthDate    string
	ClientNumber string
}

type GetUserResponse struct {
	Code    int
	Message string
	Result  GUResult
	RemcoID int
}

type GUResult struct {
	Client Client
}

type Client struct {
	ClientID     int
	ClientNumber string
	FirstName    string
	MiddleName   string
	LastName     string
	BirthDate    string
	CPCountry    CrtyID
	TPCountry    CrtyID
	CtryAddress  CrtyID
	CSOfFund     CSOfFund
}

type CrtyID struct {
	CountryID int
}

type CSOfFund struct {
	SourceOfFundID int
}

type CebAddBftReq struct {
	FirstName          string
	MiddleName         string
	LastName           string
	SenderClientID     int
	BirthDate          string
	CellphoneCountryID string
	ContactNumber      string
	TelephoneCountryID string
	TelephoneAreaCode  string
	TelephoneNumber    string
	CountryAddressID   string
	BirthCountryID     string
	ProvinceAddress    string
	Address            string
	UserID             int
	Occupation         string
	ZipCode            string
	StateIDAddress     string
	Tin                string
}

type CebAddBfResp struct {
	Code    int      `json:"code"`
	Message string   `json:"message"`
	Result  ABResult `json:"result"`
	RemcoID int      `json:"remco_id"`
}

type ABResult struct {
	ResultStatus  string `json:"ResultStatus"`
	MessageID     int    `json:"MessageID"`
	LogID         int    `json:"LogID"`
	BeneficiaryID int    `json:"BeneficiaryID"`
}

type InputGuideRequest struct {
	Ptnr      string
	SrcCtry   string
	SrcCncy   string
	AgentCode string
	CtryCode  string
	City      string
	ID        int
}
