package storage

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"github.com/lib/pq"
)

type InviteStatus = string

const (
	Invited    InviteStatus = "Invited"
	InviteSent InviteStatus = "Invite Sent"
	Expired    InviteStatus = "Expired"
	InProgress InviteStatus = "In-Progress"
	Revoked    InviteStatus = "Revoked"
	Approved   InviteStatus = "Approved"
)

var (
	// NotFound is returned when the requested resource does not exist.
	NotFound = errors.New("not found")
	// Conflict is returned when trying to create the same resource twice.
	Conflict = errors.New("conflict")
	// EmailExists is returned when signup email already exists in storage.
	EmailExists = errors.New("email already exists")
	// UsernameExists is returned when signup UserName already exists in storage.
	UsernameExists = errors.New("username already exists")
	// UsernameExists is returned when signup UserName already exists in storage.
	InvalidArgument = errors.New("Invalid Argument")
)

type Session struct {
	ID      string       `db:"id"`
	UserID  string       `db:"user_id"`
	Expiry  sql.NullTime `db:"session_expiry"`
	Created time.Time    `db:"created"`
	Updated time.Time    `db:"updated"`
	Deleted sql.NullTime `db:"deleted"`
}
type UserProfile struct {
	ID             string       `db:"id"`
	UserID         string       `db:"user_id"`
	OrgID          string       `db:"org_id"`
	Email          string       `db:"email"`
	ProfilePicture string       `db:"profile_picture"`
	Created        time.Time    `db:"created"`
	Updated        time.Time    `db:"updated"`
	Deleted        sql.NullTime `db:"deleted"`
}
type ApiKeyTransactionType struct {
	ID              string       `db:"id"`
	UserID          string       `db:"user_id"`
	OrgID           string       `db:"org_id"`
	ClientID        string       `db:"client_id"`
	Environment     string       `db:"environment"`
	TransactionType string       `db:"transaction_type"`
	Created         time.Time    `db:"created"`
	Deleted         sql.NullTime `db:"deleted"`
}
type OrgProfile struct {
	ID               string `db:"id"`
	UserID           string `db:"user_id"`
	OrgID            string `db:"org_id"`
	OrgType          int    `db:"org_type"`
	Status           int    `db:"status"`
	RiskScore        int    `db:"risk_score"`
	TransactionTypes string `db:"transaction_types"`
	BusinessInfo
	AccountInfo
	DateApplied       sql.NullTime `db:"date_applied"`
	ReminderSent      int          `db:"reminder_sent"`
	Created           time.Time    `db:"created"`
	Updated           time.Time    `db:"updated"`
	Deleted           sql.NullTime `db:"deleted"`
	DsaCode           string       `db:"dsa_code"`
	TerminalIdOtc     string       `db:"terminal_id_otc"`
	TerminalIdDigital string       `db:"terminal_id_digital"`
	IsProvider        bool         `db:"is_provider"`
	Partner           string       `db:"partner"`
	// Used for list query only
	Count int
}
type BusinessInfo struct {
	CompanyName   string `db:"bus_info_company_name"`
	StoreName     string `db:"bus_info_store_name"`
	PhoneNumber   string `db:"bus_info_phone_number"`
	FaxNumber     string `db:"bus_info_fax_number"`
	Website       string `db:"bus_info_website"`
	CompanyEmail  string `db:"bus_info_company_email"`
	ContactPerson string `db:"bus_info_contact_person"`
	Position      string `db:"bus_info_position"`
	Address
}
type AccountInfo struct {
	Bank                    string `db:"acc_info_bank"`
	BankAccountNumber       string `db:"acc_info_bank_account_number"`
	BankAccountHolder       string `db:"acc_info_bank_account_holder"`
	AgreeTermsConditions    int    `db:"acc_info_agree_terms_conditions"`
	AgreeOnlineSupplierForm int    `db:"acc_info_agree_online_supplier_form"`
	Currency                int    `db:"acc_info_currency"`
}
type Address struct {
	Address1   string `db:"bus_info_address1"`
	City       string `db:"bus_info_city"`
	State      string `db:"bus_info_state"`
	PostalCode string `db:"bus_info_postal_code"`
}
type BranchAddress struct {
	Address1   string `db:"address1"`
	City       string `db:"city"`
	State      string `db:"state"`
	PostalCode string `db:"postal_code"`
}
type Branch struct {
	ID           string `db:"id"`
	OrgID        string `db:"org_id"`
	OrgProfileID string `db:"org_profile_id"`
	Title        string `db:"title"`
	BranchAddress
	PhoneNumber   string       `db:"phone_number"`
	FaxNumber     string       `db:"fax_number"`
	ContactPerson string       `db:"contact_person"`
	Created       time.Time    `db:"created"`
	Updated       time.Time    `db:"updated"`
	Deleted       sql.NullTime `db:"deleted"`
	// Used for list query only
	Count int
}
type Amount struct {
	Amount   string `json:"amount,omitempty"`
	Currency string `json:"currency,omitempty"`
}
type FeeCommission struct {
	ID               string       `db:"id"`
	OrgID            string       `db:"org_id"`
	OrgProfileID     string       `db:"org_profile_id"`
	Type             int          `db:"fee_commision_type"`
	FeeAmount        string       `db:"fee_amount"`
	CommissionAmount string       `db:"commission_amount"`
	FeeStatus        int          `db:"-"`
	StartDate        sql.NullTime `db:"start_date"`
	EndDate          sql.NullTime `db:"end_date"`
	Created          time.Time    `db:"created"`
	Updated          time.Time    `db:"updated"`
	Deleted          sql.NullTime `db:"deleted"`
	Count            int
}
type Rate struct {
	ID              string `db:"id"`
	FeeCommissionID string `db:"fee_commission_id"`
	MinVolume       string `db:"min_volume"`
	MaxVolume       string `db:"max_volume"`
	TxnRate         string `db:"txn_rate"`
}
type Partner struct {
	ID        string       `db:"id"`
	OrgID     string       `db:"org_id"`
	Type      string       `db:"type"`
	Partner   string       `db:"service"`
	Created   time.Time    `db:"created"`
	Updated   time.Time    `db:"updated"`
	Deleted   sql.NullTime `db:"deleted"`
	UpdatedBy string       `db:"updated_by"`
	Status    string       `db:"status"`
}
type Question struct {
	ID             string    `db:"id"`
	OrgID          string    `db:"org_id"`
	UserID         string    `db:"user_id"`
	QID            string    `db:"qid"`
	ANS            string    `db:"ans"`
	QType          string    `db:"qtype"`
	CustomersTotal string    `db:"customers_total"`
	HrTotal        string    `db:"hr_total"`
	ImpactScore    string    `db:"impact_score"`
	Created        time.Time `db:"created"`
	Updated        time.Time `db:"updated"`
}
type WesternUnionPartner struct {
	Coy        string    `json:"coy"`
	TerminalID string    `json:"terminal_id"`
	UpdatedBy  string    `json:"updated_by"`
	Status     string    `json:"status"`
	Created    time.Time `json:"created"`
	Updated    time.Time `json:"updated"`
	StartDate  time.Time `json:"startdate"`
	EndDate    time.Time `json:"enddate"`
}
type IRemitPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type JapanRemitPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type InstantCashPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type TransfastPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type RemitlyPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type RiaPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type MetroBankPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type BPIPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type USSCPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type UnitellerPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type CebuanaPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type TransferWisePartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type CebuanaIntlPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type AyannahPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type IntelExpressPartner struct {
	Param1    string    `json:"param1"`
	Param2    string    `json:"param2"`
	UpdatedBy string    `json:"updated_by"`
	Status    string    `json:"status"`
	Created   time.Time `json:"created"`
	Updated   time.Time `json:"updated"`
	StartDate time.Time `json:"startdate"`
	EndDate   time.Time `json:"enddate"`
}
type FileUpload struct {
	FileID     string       `db:"file_id"`
	OrgID      string       `db:"org_id"`
	UserID     string       `db:"user_id"`
	UploadType string       `db:"upload_type"`
	FileNames  string       `db:"file_names"`
	BucketName string       `db:"bucket_name"`
	Submitted  int          `db:"submitted"`
	Checked    sql.NullTime `db:"checked"`
	Created    time.Time    `db:"created"`
	FileName   string       `db:"file_name"`
}
type FilterList struct {
	Limit             int32
	Offset            int32
	CompanyName       string
	SortBy            string
	SortByColumn      string
	DateApplied       string
	RiskScore         pq.Int32Array
	Status            pq.Int32Array
	SubmittedDocument string
	OrgType           string
	IsProvider        bool
}
type EventData struct {
	EventID  string       `db:"event_id"`
	Resource string       `db:"resource"`
	Action   string       `db:"action"`
	Data     string       `db:"data"`
	Created  time.Time    `db:"created"`
	Updated  time.Time    `db:"updated"`
	Deleted  sql.NullTime `db:"deleted"`
}
type Role struct {
	ID        string       `db:"id"`
	Data      RoleData     `db:"role_data"`
	CreatedBy string       `db:"created_by"`
	UpdatedBy string       `db:"updated_by"`
	Created   time.Time    `db:"created"`
	Updated   time.Time    `db:"updated"`
	Deleted   sql.NullTime `db:"deleted"`
}
type RoleData struct {
	Permissions []string
}
type FileUploadFilter struct {
	UploadTypes []string
}
type LimitOffsetFilter struct {
	Limit  int32
	Offset int32
	Type   int
}
type PartnerList struct {
	ID               string         `db:"id"`
	Stype            string         `db:"stype"`
	Name             string         `db:"name"`
	Created          time.Time      `db:"created"`
	Updated          time.Time      `db:"updated"`
	Deleted          sql.NullTime   `db:"deleted"`
	Status           string         `db:"status"`
	TransactionTypes string         `db:"transaction_types"`
	ServiceName      string         `db:"service_name"`
	UpdatedBy        string         `db:"updated_by"`
	DisableReason    sql.NullString `db:"disable_reason"`
	Platform         string         `db:"platform"`
	IsProvider       bool           `db:"is_provider"`
	PerahubPartnerID string         `db:"perahub_partner_id"`
	RemcoID          string         `db:"remco_id"`
}
type ServiceRequest struct {
	ID          string       `db:"id"`
	OrgID       string       `db:"org_id"`
	Partner     string       `db:"partner"`
	SvcName     string       `db:"service_name"`
	CompanyName string       `db:"company_name"`
	Remarks     string       `db:"remarks"`
	Status      string       `db:"status"`
	Enabled     bool         `db:"enabled"`
	UpdatedBy   string       `db:"updated_by"`
	Total       string       `db:"total"`
	Applied     sql.NullTime `db:"applied"`
	Updated     time.Time    `db:"updated"`
	Created     time.Time    `db:"created"`
	Partners    string       `db:"partners"`
	Pending     int32        `db:"pending"`
	Accepted    int32        `db:"accepted"`
	Rejected    int32        `db:"rejected"`
}
type SvcRequestFilter struct {
	OrgID        []string
	Status       []string
	SvcName      []string
	Partner      []string
	CompanyName  string
	SortByColumn string
	SortOrder    string
	Limit        int
	Offset       int
}
type ValidateSvcRequestFilter struct {
	OrgID               string
	Partner             string
	SvcName             string
	IsAnyPartnerEnabled bool
}
type ValidateSvcResponse struct {
	Enabled bool
}

func (d RoleData) Value() (driver.Value, error) {
	return json.Marshal(d)
}

func (d *RoleData) Scan(value interface{}) error {
	b, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	return json.Unmarshal(b, &d)
}

type UploadServiceRequest struct {
	ID       string       `db:"id"`
	OrgID    string       `db:"org_id"`
	Partner  string       `db:"partner"`
	SvcName  string       `db:"service_name"`
	Status   string       `db:"status"`
	FileType string       `db:"file_type"`
	FileID   string       `db:"file_id"`
	CreateBy string       `db:"create_by"`
	VerifyBy string       `db:"verify_by"`
	Total    string       `db:"total"`
	Created  time.Time    `db:"created"`
	Verified sql.NullTime `db:"verified"`
}
type UploadSvcRequestFilter struct {
	OrgID   string
	Status  []string
	SvcName []string
	Partner []string
}

// partner commission management
type PartnerCommission struct {
	ID              string    `db:"id"`
	Partner         string    `db:"partner"`
	BoundType       string    `db:"bound_type"`
	RemitType       string    `db:"remit_type"`
	TransactionType string    `db:"transaction_type"`
	TierType        string    `db:"tier_type"`
	Amount          string    `db:"amount"`
	StartDate       time.Time `db:"start_date"`
	EndDate         time.Time `db:"end_date"`
	CreatedBy       string    `db:"created_by"`
	UpdatedBy       string    `db:"updated_by"`
	Created         time.Time `db:"created"`
	Updated         time.Time `db:"updated"`
	Count           int       `db:"count"`
}
type PartnerCommissionTier struct {
	ID                  string `db:"id"`
	PartnerCommissionID string `db:"partner_commission_config_id"`
	MinValue            string `db:"min_value"`
	MaxValue            string `db:"max_value"`
	Amount              string `db:"amount"`
}
type (
	TransactionType string
	TierType        string
	BoundType       string
)

const (
	// TransactionType
	TransactionType_Digital TransactionType = "DIGITAL"
	TransactionType_Otc     TransactionType = "OTC"
	// BoundType
	BoundType_In     BoundType = "InBound"
	BoundType_Out    BoundType = "OutBound"
	BoundType_Others BoundType = "Others"
	// TierType
	TierType_Fixed_Amount          TierType = "Fixed"
	TierType_Fixed_Percentage      TierType = "Percentage"
	TierType_Fixed_Tier_Amount     TierType = "TierAmount"
	TierType_Fixed_Tier_Percentage TierType = "TierPercentage"
)

// revenue sharing report
type RevenueSharing struct {
	ID              string    `db:"id"`
	OrgID           string    `db:"org_id"`
	UserID          string    `db:"user_id"`
	Partner         string    `db:"partner"`
	BoundType       string    `db:"bound_type"`
	RemitType       string    `db:"remit_type"`
	TransactionType string    `db:"transaction_type"`
	TierType        string    `db:"tier_type"`
	Amount          string    `db:"amount"`
	CreatedBy       string    `db:"created_by"`
	UpdatedBy       string    `db:"updated_by"`
	Created         time.Time `db:"created"`
	Updated         time.Time `db:"updated"`
	Count           int       `db:"count"`
}
type RevenueSharingTier struct {
	ID               string `db:"id"`
	RevenueSharingID string `db:"revenue_sharing_id"`
	MinValue         string `db:"min_value"`
	MaxValue         string `db:"max_value"`
	Amount           string `db:"amount"`
}

type SortOrder string

const (
	Asc  SortOrder = "ASC"
	Desc SortOrder = "DESC"
)

type (
	RevenueSharingReportColumn string
)

const (
	OrgIDCol   RevenueSharingReportColumn = "OrgIDCol"
	UserColID  RevenueSharingReportColumn = "UserColID"
	StatusCol  RevenueSharingReportColumn = "StatusCol"
	PartnerCol RevenueSharingReportColumn = "PartnerCol"
)

type RevenueSharingReport struct {
	ID                string    `db:"id"`
	OrgID             string    `db:"org_id"`
	DsaCode           string    `db:"dsa_code"`
	YearMonth         string    `db:"year_month"`
	RemittanceCount   int       `db:"remittance_count"`
	CicoCount         int       `db:"cico_count"`
	BillsPaymentCount int       `db:"bills_payment_count"`
	InsuranceCount    int       `db:"insurance_count"`
	DsaCommission     string    `db:"dsa_commission"`
	CommissionType    string    `db:"dsa_commission_type"`
	Status            int       `db:"status"`
	Created           time.Time `db:"created"`
	Updated           time.Time `db:"updated"`
	Count             int       `db:"count"`
	SortByColumn      RevenueSharingReportColumn
	SortOrder         SortOrder
	Limit             int
	Offset            int
	Total             int
}

type GetAllPartnerListReq struct {
	ID              string       `db:"id"`
	Stype           string       `db:"stype"`
	Name            string       `db:"name"`
	TransactionType string       `db:"transaction_type"`
	Partner         string       `db:"partner"`
	Created         time.Time    `db:"created"`
	Updated         time.Time    `db:"updated"`
	Deleted         sql.NullTime `db:"deleted"`
	Status          string       `db:"status"`
}

type DSAPartnerList struct {
	Partner         string `db:"partner"`
	TransactionType string `db:"transaction_type"`
}

type GetDSAPartnerListRequest struct {
	TransactionTypes []string `db:"transaction_type"`
}

type CICOPartnerList struct {
	ID      string       `db:"id"`
	Stype   string       `db:"stype"`
	Name    string       `db:"name"`
	Created time.Time    `db:"created"`
	Updated time.Time    `db:"updated"`
	Deleted sql.NullTime `db:"deleted"`
	Status  string       `db:"status"`
}

type GetTransactionTypeByClientIdResponse struct {
	Environment     string `db:"environment"`
	TransactionType string `db:"transaction_type"`
}

type UpdateOrgProfileOrgIDUserID struct {
	OldOrgID string `db:"oldOrgID"`
	NewOrgID string `db:"newOrgID"`
	UserID   string `db:"user_id"`
}

type UpdateServiceRequestOrgID struct {
	OldOrgID string `db:"oldOrgID"`
	NewOrgID string `db:"newOrgID"`
	Status   string `db:"status"`
}
