package storage

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx/types"
)

type MicroInsuranceSort string

const (
	TraceNumberMICol MicroInsuranceSort = "trace_number"
	DsaIdMICol       MicroInsuranceSort = "dsa_id"
	TrxStatusMICol   MicroInsuranceSort = "trx_status"
	UserCodeMICol    MicroInsuranceSort = "user_code"
	LocationIDMICol  MicroInsuranceSort = "location_id"
	Fee              MicroInsuranceSort = "Fee"
	Commission       MicroInsuranceSort = "Commission"
	TotalAmount      MicroInsuranceSort = "TotalAmount"
	TranTime         MicroInsuranceSort = "TransactionCompletedTime"
)

type MicroInsuranceHistory struct {
	ID               string         `db:"id"`
	DsaID            string         `db:"dsa_id"`
	Coy              string         `db:"coy"`
	LocationID       string         `db:"location_id"`
	UserCode         string         `db:"user_code"`
	TrxDate          time.Time      `db:"trx_date"`
	PromoAmount      string         `db:"promo_amount"`
	PromoCode        string         `db:"promo_code"`
	Amount           string         `db:"amount"`
	CoverageCount    string         `db:"coverage_count"`
	ProductCode      string         `db:"product_code"`
	ProcessingBranch string         `db:"processing_branch"`
	ProcessedBy      string         `db:"processed_by"`
	UserEmail        string         `db:"user_email"`
	LastName         string         `db:"last_name"`
	FirstName        string         `db:"first_name"`
	MiddleName       string         `db:"middle_name"`
	Gender           string         `db:"gender"`
	Birthdate        time.Time      `db:"birthdate"`
	MobileNumber     string         `db:"mobile_number"`
	ProvinceCode     string         `db:"province_code"`
	CityCode         string         `db:"city_code"`
	Address          string         `db:"address"`
	MaritalStatus    string         `db:"marital_status"`
	Occupation       string         `db:"occupation"`
	CardNumber       string         `db:"card_number"`
	NumberUnits      string         `db:"number_units"`
	Beneficiaries    types.JSONText `db:"beneficiaries"`
	Dependents       types.JSONText `db:"dependents"`
	TrxStatus        string         `db:"trx_status"`
	TraceNumber      sql.NullString `db:"trace_number"`
	InsuranceDetails types.JSONText `db:"insurance_details"`
	ErrorCode        string         `db:"error_code"`
	ErrorMsg         string         `db:"error_message"`
	ErrorType        string         `db:"error_type"`
	ErrorTime        sql.NullTime   `db:"error_time"`
	Created          time.Time      `db:"created"`
	Updated          time.Time      `db:"updated"`
	OrgID            string         `db:"org_id"`
	Total            int
}

type MicroInsurancePerson struct {
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	MiddleName    string `json:"middle_name"`
	NoMiddleName  bool   `json:"no_middle_name"`
	ContactNumber string `json:"contact_number"`
	BirthDate     string `json:"birth_date"`
	Relationship  string `json:"relationship"`
}

type MicroInsuranceFilter struct {
	TraceNumber  string
	DsaID        string
	UserCode     string
	TrxStatus    string
	TrxDate      time.Time
	SortByColumn MicroInsuranceSort
	SortOrder    SortOrder
	Limit        int
	Offset       int
	From         time.Time
	Until        time.Time
	OrgID        string
}
