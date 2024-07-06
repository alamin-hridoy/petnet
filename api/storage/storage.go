package storage

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/bojanz/currency"
	"github.com/jmoiron/sqlx/types"
	"github.com/lib/pq"
)

var (
	// ErrNotFound is returned when the requested resource does not exist.
	ErrNotFound = status.Error(codes.NotFound, "not found")
	// Conflict is returned when trying to create the same resource twice.
	Conflict   = status.Error(codes.AlreadyExists, "conflict")
	ErrInvalid = status.Error(codes.InvalidArgument, "invalid arguments")
)

type (
	TxnType   string
	TxnStatus string
	TxnStep   string
)

const (
	SendType     TxnType = "SEND"
	DisburseType TxnType = "DISBURSE"
)

const (
	SuccessStatus TxnStatus = "SUCCESS"
	FailStatus    TxnStatus = "FAIL"
)

const (
	StageStep   TxnStep = "STAGE"
	ConfirmStep TxnStep = "CONFIRM"
)

const (
	IGIDsLabel           = "Identification Type"
	IGRelationsLabel     = "Relationship"
	IGOccupsLabel        = "Occupation"
	IGPurposesLabel      = "Purpose"
	IGCountryLabel       = "Country"
	IGStateLabel         = "State"
	IGUsStateLabel       = "US State"
	IGCurrencyLabel      = "Currency"
	IGFundsLabel         = "Source Of Funds"
	IGPositionLabel      = "Position"
	IGOtherInfoLabel     = "Other Info"
	IGProvincesCityLabel = "Provinces City"
	IGBrgyLabel          = "Brgy"
	IGPartnerLabel       = "Partner"
	IGEmploymentLabel    = "Employment"
)

const (
	IGIDsGroup           = "ids"
	IGRelationsGroup     = "relations"
	IGOccupationsGroup   = "occupations"
	IGPurposesGroup      = "purposes"
	IGCountryGroup       = "country"
	IGStateGroup         = "state"
	IGUsStateGroup       = "usa-states"
	IGCurrencyGroup      = "currency"
	IGFundsGroup         = "funds"
	IGPositionGroup      = "positions"
	IGProvincesCityGroup = "provincescity"
	IGBrgyGroup          = "brgys"
	IGPartnerGroup       = "partners"
	IGEmploymentGroup    = "employments"
)

func (ig *InputGuide) ToGuide(label string) []*ppb.Input {
	if ig == nil {
		return nil
	}
	g := make([]*ppb.Input, len(ig.Data[label]))
	for i, c := range ig.Data[label] {
		g[i] = &ppb.Input{
			Value:         c.Value,
			Name:          c.Name,
			Description:   c.Description,
			HasIssueDate:  c.HasIssueDate,
			HasExpiration: c.HasExpiration,
			CountryName:   c.CountryName,
			StateName:     c.StateName,
			CountryCode:   c.CountryCode,
			CurrencyCode:  c.CurrencyCode,
		}
	}
	return g
}

type RemitCache struct {
	TxnID             string         `db:"transaction_id"`
	DsaID             string         `db:"dsa_id"`
	UserID            string         `db:"user_id"`
	RemcoID           string         `db:"remco_id"`
	RemType           string         `db:"remit_type"`
	PtnrRemType       string         `db:"partner_remit_type"`
	RemcoMemberID     string         `db:"remco_member_id"`
	RemcoControlNo    string         `db:"remco_control_number"`
	RemcoAltControlNo string         `db:"remco_alternate_control_number"`
	Step              TxnStep        `db:"status"`
	Remit             types.JSONText `db:"remit"`
	Updated           time.Time      `db:"updated"`
	Created           time.Time      `db:"created"`
}

type Taxes struct {
	Currency  string `json:"currency,omitempty"`
	State     string `json:"state,omitempty"`
	County    string `json:"county,omitempty"`
	Municipal string `json:"municipal,omitempty"`
	Total     string `json:"total,omitempty"`
}

type Amount struct {
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}

type Charges map[string]Amount

type RemitHistory struct {
	TxnID                 string         `db:"remit_id"`
	DsaID                 string         `db:"dsa_id"`
	DsaOrderID            string         `db:"dsa_order_id"`
	UserID                string         `db:"user_id"`
	RemcoID               string         `db:"remco_id"`
	RemType               string         `db:"remit_type"`
	TxnStatus             string         `db:"txn_status"`
	TxnStep               string         `db:"txn_step"`
	SenderID              string         `db:"sender_member_id"`
	ReceiverID            string         `db:"receiver_member_id"`
	RemcoControlNo        string         `db:"remco_control_number"`
	ErrorCode             string         `db:"error_code"`
	ErrorMsg              string         `db:"error_message"`
	ErrorTime             sql.NullTime   `db:"error_time"`
	ErrorType             string         `db:"error_type"`
	Remittance            Remittance     `db:"remittance"`
	TxnStagedTime         sql.NullTime   `db:"txn_staged_time"`
	TxnCompletedTime      sql.NullTime   `db:"txn_completed_time"`
	Updated               time.Time      `db:"updated"`
	TransactionType       sql.NullString `db:"transaction_type"`
	RemitTransactionCount int            `db:"remit_transaction_count"`
	Total                 int
}

type Remittance struct {
	CustomerTxnID     string                    `json:"customer_txn_id"`
	RemcoAltControlNo string                    `json:"remco_alternate_control_number"`
	TxnType           string                    `json:"txn_type"`
	Remitter          Contact                   `json:"remitter"`
	Receiver          Contact                   `json:"receiver"`
	Business          Business                  `json:"business"`
	Account           Account                   `json:"account"`
	GrossTotal        GrossTotal                `json:"gross_total,omitempty"`
	SourceAmt         currency.Minor            `json:"source_amt,omitempty"`
	DestAmt           currency.Minor            `json:"dest_amt,omitempty"`
	Taxes             map[string]currency.Minor `json:"taxes,omitempty"`
	Tax               currency.Minor            `json:"tax,omitempty"`
	Charges           map[string]currency.Minor `json:"charges,omitempty"`
	Charge            currency.Minor            `json:"charge,omitempty"`
}

type GrossTotal struct {
	currency.Minor `json:"amount,omitempty"`
}

type Contact struct {
	FirstName     string `json:"first_name"`
	MiddleName    string `json:"middle_name"`
	LastName      string `json:"last_name"`
	RemcoMemberID string `json:"remco_member_id"`
	Email         string `json:"email"`
	Address1      string `json:"address1"`
	Address2      string `json:"address2"`
	City          string `json:"city"`
	State         string `json:"state"`
	PostalCode    string `json:"postal_code"`
	Country       string `json:"country"`
	Province      string `json:"province"`
	Zone          string `json:"zone"`
	PhoneCty      string `json:"phone_cty"`
	Phone         string `json:"phone"`
	MobileCty     string `json:"mobile_cty"`
	Mobile        string `json:"mobile"`
	Message       string `json:"message"`
}

type Business struct {
	Name      string `json:"name"`
	Account   string `json:"account"`
	ControlNo string `json:"control_no"`
	Country   string `json:"country"`
}

type Account struct {
	BIC     string `json:"bic"`
	AcctNo  string `json:"acct_no"`
	AcctSfx string `json:"acct_sfx"`
}

type RemittanceHistory struct {
	TxnID                   string `db:"transaction_id"`
	DsaID                   string `db:"dsa_id"`
	UserID                  string `db:"user_id"`
	RemcoID                 string `db:"remco_id"`
	RemcoMemberID           string `db:"remco_member_id"`
	RemcoControlNo          string `db:"remco_control_number"`
	RemcoAlternateControlNo string `db:"remco_alternate_control_number"`
	TXNType                 string `db:"txn_type"`

	SourceCurrency           string         `db:"source_currency"`
	DestinationCurrency      string         `db:"destination_currency"`
	ExchangeRate             string         `db:"exchange_rate"`
	SourceGrossAmount        string         `db:"source_gross_amount"`
	DestinationRemitAmount   bool           `db:"destination_remit_amount"`
	AdditionalChargeAmount   string         `db:"additional_charge_amount"`
	AdditionalChargeCurrency string         `db:"additional_charge_currency"`
	Taxes                    Taxes          `db:"taxes"`
	Charges                  Charges        `db:"charges"`
	PromoCode                string         `db:"promo_code"`
	PromoDescription         string         `db:"promo_description"`
	Message                  pq.StringArray `db:"message"`

	OriginCity       string `db:"origin_city"`
	OriginState      string `db:"origin_state"`
	DestinationCity  string `db:"destination_city"`
	DestinationState string `db:"destination_state"`

	BankBIC           string `db:"bank_bic"`
	BankLocation      string `db:"bank_location"`
	BankAccountNO     string `db:"bank_account_no"`
	BankAccountSuffix string `db:"bank_account_suffix"`

	SendingReason      string `db:"sending_reason"`
	SenderRelationship string `db:"sender_relationship"`

	RemitterFirstname     string `db:"remitter_fname"`
	RemitterMiddlename    string `db:"remitter_mname"`
	RemitterLastname      string `db:"remitter_lname"`
	RemitterGender        string `db:"remitter_gender"`
	RemitterAddress1      string `db:"remitter_address1"`
	RemitterAddress2      string `db:"remitter_address2"`
	RemitterCity          string `db:"remitter_city"`
	RemitterState         string `db:"remitter_state"`
	RemitterPostalCode    string `db:"remitter_postal_code"`
	RemitterCountry       string `db:"remitter_country"`
	RemitterPhoneCountry  string `db:"remitter_phone_country"`
	RemitterPhoneNumber   string `db:"remitter_phone_number"`
	RemitterMobileCountry string `db:"remitter_mobile_country"`
	RemitterMobileNumber  string `db:"remitter_mobile_number"`
	RemitterEmail         string `db:"remitter_email"`

	ReceiverFirstname     string `db:"receiver_fname"`
	ReceiverMiddlename    string `db:"receiver_mname"`
	ReceiverLastname      string `db:"receiver_lname"`
	ReceiverAddress1      string `db:"receiver_address1"`
	ReceiverAddress2      string `db:"receiver_address2"`
	ReceiverCity          string `db:"receiver_city"`
	ReceiverState         string `db:"receiver_state"`
	ReceiverPostalCode    string `db:"receiver_postal_code"`
	ReceiverCountry       string `db:"receiver_country"`
	ReceiverPhoneCountry  string `db:"receiver_phone_country"`
	ReceiverPhoneNumber   string `db:"receiver_phone_number"`
	ReceiverMobileCountry string `db:"receiver_mobile_country"`
	ReceiverMobileNumber  string `db:"receiver_mobile_number"`

	ErrorCode    string `db:"error_code"`
	ErrorMessage string `db:"error_message"`
	ErrorDetails string `db:"error_details"`

	Errored time.Time `db:"errored"`
	TxnTime time.Time `db:"txn_time"`
	Updated time.Time `db:"updated"`
	Created time.Time `db:"created"`
}

type WURef struct {
	Code string `db:"code"`
	Name string `db:"name"`
}

type ISOCty struct {
	Code        string `db:"code"`
	Country     string `db:"country"`
	Nationality string `db:"nationality"`
	NumCode     string `db:"iso"`
}

type InputGuideData map[string][]Input

type InputGuide struct {
	Partner string         `db:"partner"`
	Data    InputGuideData `db:"inputguide"`
	Created time.Time      `db:"created"`
	Updated time.Time      `db:"updated"`
}

type Input struct {
	Value         string `json:"code"`
	Name          string `json:"name"`
	CountryName   string `json:"country_name"`
	CountryCode   string `json:"country_code"`
	CurrencyCode  string `json:"currency_code"`
	StateName     string `json:"state_name"`
	Description   string `json:"description"`
	HasIssueDate  string `json:"has_issue_date"`
	HasExpiration string `json:"has_expiration"`
}

func (t Taxes) Value() (driver.Value, error) {
	return json.Marshal(t)
}

func (t *Taxes) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &t)
}

func (c Charges) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Charges) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}

func (c Remittance) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *Remittance) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(b, &c)
}

func (c *InputGuideData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

func (c InputGuideData) Value() (driver.Value, error) {
	return json.Marshal(c)
}

func (c *GrossTotal) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &c)
}

func (c GrossTotal) Value() (driver.Value, error) {
	return json.Marshal(c)
}

type SortOrder string

const (
	Asc  SortOrder = "ASC"
	Desc SortOrder = "DESC"
)

type Column string

const (
	ControlNumberCol         Column = "ControlNumber"
	RemittedToCol            Column = "RemittedTo"
	TotalRemittedAmountCol   Column = "TotalRemittedAmount"
	TransactionTimeCol       Column = "TransactionTime"
	TransactionCompletedTime Column = "TransactionCompletedTime"
	UserIDCol                Column = "UserID"
	PartnerCol               Column = "Partner"
)

type LRHFilter struct {
	ControlNo       []string
	DsaOrgID        string
	Partner         string
	ExcludePartner  string
	ExcludeType     string
	TxnStep         string
	TxnStatus       string
	RemType         string
	SortByColumn    Column
	SortOrder       SortOrder
	Limit           int
	Offset          int
	From            time.Time
	Until           time.Time
	Transactiontype string
	DsaID           string
}

type BillDetails struct {
	Message         string         `json:"message"`
	Timestamp       string         `json:"timestamp"`
	ReferenceNumber string         `json:"reference_number"`
	Status          string         `json:"status"`
	ServiceCharge   currency.Minor `json:"service_charge"`
	TransactionID   string         `json:"transaction_id"`
	ClientReference string         `json:"client_reference"`
	BillerReference string         `json:"biller_reference"`
	PaymentMethod   string         `json:"payment_method"`
	Amount          currency.Minor `json:"amount"`
	OtherCharges    currency.Minor `json:"other_charges"`
	CreatedAt       string         `json:"created_at"`
	URL             string         `json:"url"`
}

type Bills struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Result  BillDetails `json:"result"`
	RemcoID int         `json:"remco_id"`
}

type BillPayment struct {
	BillPaymentID           string         `db:"bill_payment_id"`
	BillID                  int32          `db:"bill_id"`
	UserID                  string         `db:"user_id"`
	SenderMemberID          string         `db:"sender_member_id"`
	BillPaymentStatus       string         `db:"bill_payment_status"`
	ErrorCode               string         `db:"error_code"`
	ErrorMsg                string         `db:"error_message"`
	ErrorType               string         `db:"error_type"`
	Bills                   types.JSONText `db:"bills"`
	BillerTag               string         `db:"biller_tag"`
	LocationID              string         `db:"location_id"`
	CurrencyID              string         `db:"currency_id"`
	AccountNumber           string         `db:"account_number"`
	Amount                  types.JSONText `db:"amount"`
	Identifier              string         `db:"identifier"`
	Coy                     string         `db:"coy"`
	ServiceCharge           types.JSONText `db:"service_charge"`
	TotalAmount             types.JSONText `db:"total_amount"`
	BillPaymentDate         time.Time      `db:"bill_payment_date"`
	PartnerID               string         `db:"partner_id"`
	BillerName              string         `db:"biller_name"`
	RemoteUserID            string         `db:"remote_user_id"`
	CustomerID              string         `db:"customer_id"`
	RemoteLocationID        string         `db:"remote_location_id"`
	LocationName            string         `db:"location_name"`
	FormType                string         `db:"form_type"`
	FormNumber              string         `db:"form_number"`
	PaymentMethod           string         `db:"payment_method"`
	OtherInfo               types.JSONText `db:"other_info"`
	TrxDate                 time.Time      `db:"trx_date"`
	Total                   int            `db:"total"`
	Created                 time.Time      `db:"created"`
	Updated                 time.Time      `db:"updated"`
	ClientRefNumber         string         `db:"client_reference_number"`
	BillPartnerID           int32          `db:"bill_partner_id"`
	PartnerCharge           string         `db:"partner_charge"`
	ReferenceNumber         string         `db:"reference_number"`
	ValidationNumber        string         `db:"validation_number"`
	ReceiptValidationNumber string         `db:"receipt_validation_number"`
	TpaID                   string         `db:"tpa_id"`
	Type                    string         `db:"type"`
	TxnID                   string         `db:"txnid"`
	OrgID                   string         `db:"org_id"`
}

type BillsOtherInfo struct {
	LastName        string `db:"LastName"`
	FirstName       string `db:"FirstName"`
	MiddleName      string `db:"MiddleName"`
	PaymentType     string `db:"PaymentType"`
	Course          string `db:"Course"`
	TotalAssessment string `db:"TotalAssessment"`
	SchoolYear      string `db:"SchoolYear"`
	Term            string `db:"Term"`
}

type BillPayColumn string

const (
	BillPaymentIDCol  BillPayColumn = "BillPaymentIDCol"
	BillIDCol         BillPayColumn = "BillIDCol"
	DsaIDCol          BillPayColumn = "DsaIDCol"
	SenderMemberIDCol BillPayColumn = "SenderMemberIDCol"
	UserCol           BillPayColumn = "UserCol"
	FeesCol           BillPayColumn = "Fee"
	CommissionsCol    BillPayColumn = "Commission"
	TransactTimeCol   BillPayColumn = "TransactionCompletedTime"
	AmountCol         BillPayColumn = "GrossAmount"
)

type BillPaymentFilter struct {
	BillPaymentID     string
	BillID            string
	UserID            string
	SenderMemberID    string
	BillPaymentStatus string
	TrxDate           time.Time
	SortByColumn      BillPayColumn
	SortOrder         SortOrder
	Limit             int
	Offset            int
	From              time.Time
	Until             time.Time
	OrgID             string
	ReferenceNumber   string
	ExcludePartners   []string
}

type (
	RemittanceHistoryColumn string
	TransactionStatus       string
)

const (
	RemittanceHistoryIDCol RemittanceHistoryColumn = "RemittanceHistoryIDCol"
	DsaColID               RemittanceHistoryColumn = "DsaColID"
	PhrnIDCol              RemittanceHistoryColumn = "PhrnIDCol"
	UserColID              RemittanceHistoryColumn = "UserColID"
)

const (
	VALIDATE_SEND    TransactionStatus = "VALIDATE_SEND"
	CONFIRM_SEND     TransactionStatus = "CONFIRM_SEND"
	CANCEL_SEND      TransactionStatus = "CANCEL_SEND"
	VALIDATE_RECEIVE TransactionStatus = "VALIDATE_RECEIVE"
	CONFIRM_RECEIVE  TransactionStatus = "CONFIRM_RECEIVE"
	TRANSACTION_FAIL TransactionStatus = "FAIL"
)

type PerahubRemittanceHistory struct {
	RemittanceHistoryID           string            `db:"remittance_history_id"`
	DsaID                         string            `db:"dsa_id"`
	UserID                        string            `db:"user_id"`
	Phrn                          string            `db:"phrn"`
	SendValidateReferenceNumber   string            `db:"send_validate_reference_number"`
	CancelSendReferenceNumber     string            `db:"cancel_send_reference_number"`
	PayoutValidateReferenceNumber string            `db:"payout_validate_reference_number"`
	TxnStatus                     TransactionStatus `db:"txn_status"`
	ErrorCode                     string            `db:"error_code"`
	ErrorMessage                  string            `db:"error_message"`
	ErrorTime                     string            `db:"error_time"`
	ErrorType                     string            `db:"error_type"`
	Details                       types.JSONText    `db:"details"`
	Remarks                       string            `db:"remarks"`
	TxnCreatedTime                time.Time         `db:"txn_created_time"`
	TxnUpdatedTime                time.Time         `db:"txn_updated_time"`
	TxnConfirmTime                time.Time         `db:"txn_confirm_time"`
	Total                         int               `db:"total"`
	PayHisErr                     error
}

type PerahubRemittanceHistoryDetails struct {
	PartnerReferenceNumber string `db:"partner_reference_number"`
	PrincipalAmount        string `db:"principal_amount"`
	ServiceFee             string `db:"service_fee"`
	IsoCurrency            string `db:"iso_currency"`
	ConversionRate         string `db:"conversion_rate"`
	IsoOriginatingCountry  string `db:"iso_originating_country"`
	IsoDestinationCountry  string `db:"iso_destination_country"`
	SenderLastName         string `db:"sender_last_name"`
	SenderFirstName        string `db:"sender_first_name"`
	SenderMiddleName       string `db:"sender_middle_name"`
	ReceiverLastName       string `db:"receiver_last_name"`
	ReceiverFirstName      string `db:"receiver_first_name"`
	ReceiverMiddleName     string `db:"receiver_middle_name"`
	SenderBirthDate        string `db:"sender_birth_date"`
	SenderBirthPlace       string `db:"sender_birth_place"`
	SenderBirthCountry     string `db:"sender_birth_country"`
	SenderGender           string `db:"sender_gender"`
	SenderRelationship     string `db:"sender_relationship"`
	SenderPurpose          string `db:"sender_purpose"`
	SenderOfFund           string `db:"sender_of_fund"`
	SenderOccupation       string `db:"sender_occupation"`
	SenderEmploymentNature string `db:"sender_employment_nature"`
	SendPartnerCode        string `db:"send_partner_code"`
	PayoutPartnerCode      string `db:"payout_partner_code"`
	PartnerCode            string `db:"partner_code"`
}

type RemittanceHistoryFilter struct {
	RemittanceHistoryID string
	DsaID               string
	UserID              string
	TxnStatus           string
	TxnCreatedTime      time.Time
	TxnUpdatedTime      time.Time
	TxnConfirmTime      time.Time
	Phrn                string
	SortByColumn        RemittanceHistoryColumn
	SortOrder           SortOrder
	Limit               int
	Offset              int
	From                time.Time
	Until               time.Time
}

type (
	CICOHistoryColumn string
)

const (
	CICOHistoryIDCol CICOHistoryColumn = "CICOHistoryIDCol"
	OrgIDCol         CICOHistoryColumn = "OrgIDCol"
	PartnerCodeCol   CICOHistoryColumn = "PartnerCodeCol"
	ProviderCol      CICOHistoryColumn = "ProviderCol"
	TrxTypeCol       CICOHistoryColumn = "TrxTypeCol"
	FeeCol           CICOHistoryColumn = "Fee"
	CommissionCol    CICOHistoryColumn = "Commission"
	TotalAmountCol   CICOHistoryColumn = "TotalAmount"
	TranTimeCol      CICOHistoryColumn = "TransactionCompletedTime"
)

type CashInCashOutHistory struct {
	ID                 string         `db:"id"`
	OrgID              string         `db:"org_id"`
	PartnerCode        string         `db:"partner_code"`
	SvcProvider        string         `db:"svc_provider"`
	Provider           string         `db:"trx_provider"`
	TrxType            string         `db:"trx_type"`
	ReferenceNumber    string         `db:"reference_number"`
	PetnetTrackingNo   string         `db:"petnet_trackingno"`
	ProviderTrackingNo string         `db:"provider_trackingno"`
	PrincipalAmount    int            `db:"principal_amount"`
	Charges            int            `db:"charges"`
	TotalAmount        int            `db:"total_amount"`
	TrxDate            time.Time      `db:"trx_date"`
	Details            types.JSONText `db:"details"`
	TxnStatus          string         `db:"txn_status"`
	ErrorCode          string         `db:"error_code"`
	ErrorMessage       string         `db:"error_message"`
	ErrorTime          string         `db:"error_time"`
	ErrorType          string         `db:"error_type"`
	CreatedBy          string         `db:"created_by"`
	UpdatedBy          string         `db:"updated_by"`
	Created            time.Time      `db:"created"`
	Updated            time.Time      `db:"updated"`
	Count              int
	SortByColumn       CICOHistoryColumn
	SortOrder          SortOrder
	Limit              int
	Offset             int
	From               time.Time
	Until              time.Time
	Total              int
}

type CashInCashOutHistoryDetails struct {
	ID                 string    `json:"id"`
	PartnerCode        string    `json:"partner_code"`
	Provider           string    `json:"trx_provider"`
	PetnetTrackingno   string    `json:"petnet_trackingno"`
	TrxDate            time.Time `json:"trx_date"`
	TrxType            string    `json:"trx_type"`
	ProviderTrackingNo string    `json:"provider_trackingno"`
	ReferenceNumber    string    `json:"reference_number"`
	PrincipalAmount    int       `json:"principal_amount"`
	Charges            int       `json:"charges"`
	TotalAmount        int       `json:"total_amount"`
}
type CashInCashOutHistoryRes struct {
	Code    int                         `json:"code"`
	Message string                      `json:"message"`
	Result  CashInCashOutHistoryDetails `json:"result"`
}

type CashInCashOutFilter struct {
	OrgID              string
	PartnerCode        string
	SvcProvider        string
	Provider           string
	TrxType            string
	TxnStatus          string
	ReferenceNumber    string
	PetnetTrackingNo   string
	ProviderTrackingNo string
	TrxDate            time.Time
	SortByColumn       CICOHistoryColumn
	SortOrder          SortOrder
	Limit              int
	Offset             int
	From               time.Time
	Until              time.Time
}

type CashInCashOutTrxListFilter struct {
	ReferenceNumber    string
	ExcludeProviders   []string
	OrgID              string
	PartnerCode        string
	SvcProvider        string
	Provider           string
	TrxType            string
	TxnStatus          string
	PetnetTrackingNo   string
	ProviderTrackingNo string
	SortByColumn       CICOHistoryColumn
	SortOrder          SortOrder
	Limit              int
	Offset             int
	From               time.Time
	Until              time.Time
	Total              int
}

type RemitToAccountHistoryDetails struct {
	Code            string    `json:"code"`
	SenderRefId     string    `json:"senderRefId"`
	State           string    `json:"state"`
	Uuid            string    `json:"uuid"`
	Description     string    `json:"description"`
	Type            string    `json:"type"`
	Amount          string    `json:"amount"`
	UbpTranId       string    `json:"ubpTranId"`
	TranRequestDate time.Time `json:"tranRequestDate"`
	TranFinacleDate string    `json:"tranFinacleDate"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

type (
	RTAHistoryColumn string
)

const (
	RTAHistoryIDCol RTAHistoryColumn = "RTAHistoryIDCol"
	OrgIDRTACol     RTAHistoryColumn = "OrgIDRTACol"
	TrxTypeRTACol   RTAHistoryColumn = "TrxTypeRTACol"
)

type RemitToAccountHistoryRes struct {
	Code     int                          `json:"code"`
	Message  string                       `json:"message"`
	Result   RemitToAccountHistoryDetails `json:"result"`
	RemcoID  int                          `json:"remco_id"`
	BankCode string                       `json:"bank_code"`
}

type RemitToAccountHistory struct {
	ID                          string         `db:"id"`
	OrgID                       string         `db:"org_id"`
	Partner                     string         `db:"partner"`
	ReferenceNumber             string         `db:"reference_number"`
	TrxDate                     time.Time      `db:"trx_date"`
	AccountNumber               string         `db:"account_number"`
	Currency                    string         `db:"currency"`
	ServiceCharge               string         `db:"service_charge"`
	Remarks                     string         `db:"remarks"`
	Particulars                 string         `db:"particulars"`
	MerchantName                string         `db:"merchant_name"`
	BankID                      int            `db:"bank_id"`
	LocationID                  int            `db:"location_id"`
	UserID                      int            `db:"user_id"`
	CurrencyID                  string         `db:"currency_id"`
	CustomerID                  string         `db:"customer_id"`
	FormType                    string         `db:"form_type"`
	FormNumber                  string         `db:"form_number"`
	TrxType                     string         `db:"trx_type"`
	RemoteLocationID            int            `db:"remote_location_id"`
	RemoteUserID                int            `db:"remote_user_id"`
	BillerName                  string         `db:"biller_name"`
	TrxTime                     string         `db:"trx_time"`
	TotalAmount                 string         `db:"total_amount"`
	AccountName                 string         `db:"account_name"`
	BeneficiaryAddress          string         `db:"beneficiary_address"`
	BeneficiaryBirthDate        string         `db:"beneficiary_birthdate"`
	BeneficiaryCity             string         `db:"beneficiary_city"`
	BeneficiaryCivil            string         `db:"beneficiary_civil"`
	BeneficiaryCountry          string         `db:"beneficiary_country"`
	BeneficiaryCustomerType     string         `db:"beneficiary_customertype"`
	BeneficiaryFirstName        string         `db:"beneficiary_firstname"`
	BeneficiaryLastName         string         `db:"beneficiary_lastname"`
	BeneficiaryMiddleName       string         `db:"beneficiary_middlename"`
	BeneficiaryTin              string         `db:"beneficiary_tin"`
	BeneficiarySex              string         `db:"beneficiary_sex"`
	BeneficiaryState            string         `db:"beneficiary_state"`
	CurrencyCodePrincipalAmount string         `db:"currency_code_principal_amount"`
	PrincipalAmount             string         `db:"principal_amount"`
	RecordType                  string         `db:"record_type"`
	RemitterAddress             string         `db:"remitter_address"`
	RemitterBirthDate           string         `db:"remitter_birthdate"`
	RemitterCity                string         `db:"remitter_city"`
	RemitterCivil               string         `db:"remitter_civil"`
	RemitterCountry             string         `db:"remitter_country"`
	RemitterCustomerType        string         `db:"remitter_customer_type"`
	RemitterFirstName           string         `db:"remitter_firstname"`
	RemitterGender              string         `db:"remitter_gender"`
	RemitterID                  int            `db:"remitter_id"`
	RemitterLastName            string         `db:"remitter_lastname"`
	RemitterMiddleName          string         `db:"remitter_middlename"`
	RemitterState               string         `db:"remitter_state"`
	SettlementMode              string         `db:"settlement_mode"`
	Notification                bool           `db:"notification"`
	BeneZipCode                 string         `db:"bene_zip_code"`
	Info                        types.JSONText `db:"info"`
	Details                     types.JSONText `db:"details"`
	TxnStatus                   string         `db:"txn_status"`
	ErrorCode                   string         `db:"error_code"`
	ErrorMessage                string         `db:"error_message"`
	ErrorTime                   string         `db:"error_time"`
	ErrorType                   string         `db:"error_type"`
	CreatedBy                   string         `db:"created_by"`
	UpdatedBy                   string         `db:"updated_by"`
	Created                     time.Time      `db:"created"`
	Updated                     time.Time      `db:"updated"`
	SortByColumn                RTAHistoryColumn
	SortOrder                   SortOrder
	Limit                       int
	Offset                      int
	Total                       int
}
