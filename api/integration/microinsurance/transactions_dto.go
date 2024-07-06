package microinsurance

import (
	"encoding/json"
	"time"
)

// TransactResult ...
type TransactResult struct {
	Code    string     `json:"code"`
	Message string     `json:"message"`
	Result  *Insurance `json:"result"`
}

// Insurance ...
type Insurance struct {
	SessionID         string      `json:"sessionID"`
	StatusCode        string      `json:"statusCode"`
	StatusDesc        string      `json:"statusDesc"`
	InsProductID      string      `json:"insProductID"`
	InsProductDesc    string      `json:"insProductDesc"`
	TrnDate           string      `json:"trnDate"`
	TrnAmount         json.Number `json:"trnAmount"`
	TraceNumber       string      `json:"traceNo"`
	ClientNo          string      `json:"clientNo"`
	NumUnits          json.Number `json:"numUnits"`
	BegPolicyNo       string      `json:"begPolicyNo"`
	EndPolicyNo       string      `json:"endPolicyNo"`
	EffectiveDate     string      `json:"effectiveDate"`
	ExpiryDate        string      `json:"expiryDate"`
	PocPDFLink        string      `json:"pocPDFLink"`
	CocPDFLink        string      `json:"cocPDFLink"`
	PartnerCommission json.Number `json:"partnerCommission"`
	TellerCommission  json.Number `json:"tellerCommission"`
	Timestamp         string      `json:"timestamp"`
}

// Dependent ...
type Dependent struct {
	LastName      string `json:"last_name"`
	FirstName     string `json:"first_name"`
	MiddleName    string `json:"middle_name"`
	NoMiddleName  bool   `json:"no_middle_name"`
	ContactNumber string `json:"contact_number"`
	BirthDate     string `json:"birth_date"`
	Relationship  string `json:"relationship"`
}

// TransactRequest ...
type TransactRequest struct {
	Coy               string      `json:"coy"`
	LocationID        string      `json:"location_id"`
	UserCode          string      `json:"user_code"`
	TrxDate           string      `json:"trx_date"`
	PromoAmount       json.Number `json:"promo_amount"`
	PromoCode         string      `json:"promo_code"`
	Amount            string      `json:"amount"`
	CoverageCount     string      `json:"coverage_count"`
	ProductCode       string      `json:"product_code"`
	ProcessingBranch  string      `json:"processing_branch"`
	ProcessedBy       string      `json:"processed_by"`
	UserEmail         string      `json:"user_email"`
	LastName          string      `json:"last_name"`
	FirstName         string      `json:"first_name"`
	MiddleName        string      `json:"middle_name"`
	Gender            string      `json:"gender"`
	Birthdate         string      `json:"birthdate"`
	MobileNumber      string      `json:"mobile_number"`
	ProvinceCode      string      `json:"province_code"`
	CityCode          string      `json:"city_code"`
	Address           string      `json:"address"`
	MaritalStatus     string      `json:"marital_status"`
	Occupation        string      `json:"occupation"`
	CardNumber        string      `json:"card_number"`
	Ben1LastName      string      `json:"ben1_last_name"`
	Ben1FirstName     string      `json:"ben1_first_name"`
	Ben1MiddleName    string      `json:"ben1_middle_name"`
	Ben1NoMiddleName  bool        `json:"ben1_no_middle_name"`
	Ben1ContactNumber string      `json:"ben1_contact_number"`
	Ben1Relationship  string      `json:"ben1_relationship"`
	Ben2LastName      string      `json:"ben2_last_name"`
	Ben2FirstName     string      `json:"ben2_first_name"`
	Ben2MiddleName    string      `json:"ben2_middle_name"`
	Ben2NoMiddleName  bool        `json:"ben2_no_middle_name"`
	Ben2ContactNumber string      `json:"ben2_contact_number"`
	Ben2Relationship  string      `json:"ben2_relationship"`
	Ben3LastName      string      `json:"ben3_last_name"`
	Ben3FirstName     string      `json:"ben3_first_name"`
	Ben3MiddleName    string      `json:"ben3_middle_name"`
	Ben3NoMiddleName  bool        `json:"ben3_no_middle_name"`
	Ben3ContactNumber string      `json:"ben3_contact_number"`
	Ben3Relationship  string      `json:"ben3_relationship"`
	Ben4LastName      string      `json:"ben4_last_name"`
	Ben4FirstName     string      `json:"ben4_first_name"`
	Ben4MiddleName    string      `json:"ben4_middle_name"`
	Ben4NoMiddleName  bool        `json:"ben4_no_middle_name"`
	Ben4ContactNumber string      `json:"ben4_contact_number"`
	Ben4Relationship  string      `json:"ben4_relationship"`
	NumberUnits       string      `json:"number_units"`
	Dependents        []Dependent `json:"dependents"`
}

// GetReprintRequest ...
type GetReprintRequest struct {
	TraceNumber string `json:"trace_number"`
}

// RetryTransactionRequest ...
type RetryTransactionRequest struct {
	ID string `json:"id"`
}

// GetTransactionListRequest ...
type GetTransactionListRequest struct {
	DateFrom time.Time
	DateTo   time.Time
}

func (r *GetTransactionListRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		DateFrom string `json:"date_from"`
		DateTo   string `json:"date_to"`
	}{
		DateFrom: r.DateFrom.Format("2006-01-02"),
		DateTo:   r.DateTo.Format("2006-01-02"),
	})
}

// GetTransactionListResult ...
type GetTransactionListResult struct {
	Code    string                 `json:"code"`
	Message string                 `json:"message"`
	Result  *TransactionListResult `json:"result"`
}

// TransactionListResult ...
type TransactionListResult struct {
	SessionID    string                 `json:"sessionID"`
	StatusCode   string                 `json:"statusCode"`
	StatusDesc   string                 `json:"statusDesc"`
	Transactions []InsuranceTransaction `json:"transactions"`
}

// InsuranceTransaction ...
type InsuranceTransaction struct {
	SysCode             string      `json:"sysCode"`
	ClientCode          string      `json:"clientCode"`
	BrhCode             string      `json:"brhCode"`
	TrnDate             string      `json:"trnDate"`
	TraceNo             string      `json:"traceNo"`
	ClientNo            string      `json:"clientNo"`
	LastName            string      `json:"lastName"`
	FirstName           string      `json:"firstName"`
	MiddleName          string      `json:"middleName"`
	Gender              string      `json:"gender"`
	BirthDate           string      `json:"birthDate"`
	CellNo              string      `json:"cellNo"`
	MaritalStatus       string      `json:"maritalStatus"`
	Occupation          string      `json:"occupation"`
	InsGroupID          string      `json:"insGroupID"`
	InsProductID        string      `json:"insProductID"`
	InsProductDesc      string      `json:"insProductDesc"`
	InsurerCode         string      `json:"insurerCode"`
	InsuranceType       string      `json:"insuranceType"`
	BegPolicyNo         string      `json:"begPolicyNo"`
	EndPolicyNo         string      `json:"endPolicyNo"`
	CoverageInMos       json.Number `json:"coverageInMos"`
	ContestAbilityInMos json.Number `json:"contestabilityInMos"`
	EffectiveDate       string      `json:"effectiveDate"`
	ExpiryDate          string      `json:"expiryDate"`
	InsCardNo           string      `json:"insCardNo"`
	Ben1LastName        string      `json:"ben1LastName"`
	Ben1FirstName       string      `json:"ben1FirstName"`
	Ben1MiddleName      string      `json:"ben1MiddleName"`
	Ben1Relationship    string      `json:"ben1Relationship"`
	Ben2LastName        string      `json:"ben2LastName"`
	Ben2FirstName       string      `json:"ben2FirstName"`
	Ben2MiddleName      string      `json:"ben2MiddleName"`
	Ben2Relationship    string      `json:"ben2Relationship"`
	Ben3LastName        string      `json:"ben3LastName"`
	Ben3FirstName       string      `json:"ben3FirstName"`
	Ben3MiddleName      string      `json:"ben3MiddleName"`
	Ben3Relationship    string      `json:"ben3Relationship"`
	Ben4LastName        string      `json:"ben4LastName"`
	Ben4FirstName       string      `json:"ben4FirstName"`
	Ben4MiddleName      string      `json:"ben4MiddleName"`
	Ben4Relationship    string      `json:"ben4Relationship"`
	LoanAmt             json.Number `json:"loanAmt"`
	LoanTerm            string      `json:"loanTerm"`
	NoOfMonths          json.Number `json:"noOfMonths"`
	NumUnits            json.Number `json:"numUnits"`
	PerUnitFee          json.Number `json:"perUnitFee"`
	TrnAmt              json.Number `json:"trnAmt"`
	TrnFee              json.Number `json:"trnFee"`
	TotAmt              json.Number `json:"totAmt"`
	StaffLoginName      string      `json:"staffLoginName"`
	StaffCommission     json.Number `json:"staffCommission"`
	InsTrnStatus        string      `json:"insTrnStatus"`
	ProvinceCode        string      `json:"provinceCode"`
	CityCode            string      `json:"cityCode"`
	Address             string      `json:"address"`
}
