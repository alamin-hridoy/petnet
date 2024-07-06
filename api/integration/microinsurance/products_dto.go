package microinsurance

import "encoding/json"

// GetProductResult ...
type GetProductResult struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Result  *ProductResult `json:"result"`
}

// ProductResult ...
type ProductResult struct {
	SessionID  string     `json:"sessionID"`
	StatusCode string     `json:"statusCode"`
	StatusDesc string     `json:"statusDesc"`
	Product    *Product   `json:"product"`
	Coverages  []Coverage `json:"coverages"`
}

// Product ...
type Product struct {
	InsGroupID          string      `json:"insGroupID"`
	InsProductID        string      `json:"insProductID"`
	InsProductDesc      string      `json:"insProductDesc"`
	InsurerCode         string      `json:"insurerCode"`
	InsuranceType       string      `json:"insuranceType"`
	InsuranceCategory   string      `json:"insuranceCategory"`
	PolicyNo            string      `json:"policyNo"`
	MinAge              json.Number `json:"minAge"`
	MaxAge              json.Number `json:"maxAge"`
	CoverageInMos       json.Number `json:"coverageInMos"`
	ContestAbilityInMos json.Number `json:"contestabilityInMos"`
	ActivationDelay     json.Number `json:"activationDelay"`
	MaxUnits            json.Number `json:"maxUnits"`
	PerUnitFee          json.Number `json:"perUnitFee"`
	ProductName         string      `json:"product_name"`
}

// Coverage ...
type Coverage struct {
	InsCoverageID     string `json:"insCoverageID"`
	InsCoverageDesc   string `json:"insCoverageDesc"`
	InsCoverageIconID string `json:"insCoverageIconID"`
	InsCoverageType1  string `json:"insCoverageType1"`
	InsCoverageAmt1   string `json:"insCoverageAmt1"`
	InsCoverageType2  string `json:"insCoverageType2"`
	InsCoverageAmt2   string `json:"insCoverageAmt2"`
	InsCoverageType3  string `json:"insCoverageType3"`
	InsCoverageAmt3   string `json:"insCoverageAmt3"`
	InsCoverageType4  string `json:"insCoverageType4"`
	InsCoverageAmt4   string `json:"insCoverageAmt4"`
	InsCoverageType5  string `json:"insCoverageType5"`
	InsCoverageAmt5   string `json:"insCoverageAmt5"`
}

// MinMax ...
type MinMax struct {
	Max json.Number `json:"max"`
	Min json.Number `json:"min"`
}

// MinMaxAge ...
type MinMaxAge struct {
	MaxAge json.Number `json:"max_age"`
	MinAge json.Number `json:"min_age"`
}

// OfferProduct ...
type OfferProduct struct {
	ProductName          string      `json:"product_name"`
	ProductCode          string      `json:"product_code"`
	ProductType          string      `json:"product_type"`
	Dependents           json.Number `json:"dependents"`
	Beneficiary          json.Number `json:"beneficiary"`
	BeneficiaryPolicy    *MinMax     `json:"beneficiary_policy"`
	AgePolicy            *AgePolicy  `json:"age_policy"`
	EndSpielsTitle       string      `json:"end_spiels_title"`
	EndSpielsDescription string      `json:"end_spiels_description"`
	SalesPitch           string      `json:"sales_pitch"`
	TermsAndCondition    string      `json:"terms_and_condition"`
	DataPrivacy          string      `json:"data_privacy"`
}

// ActiveProduct ...
type ActiveProduct struct {
	OfferProduct
	DependentsPolicy *MinMax `json:"dependents_policy"`
}

// AgePolicy ...
type AgePolicy struct {
	Insurer    *MinMaxAge        `json:"insurer"`
	Dependents *DependentsPolicy `json:"dependents,omitempty"`
}

// DependentsPolicy ...
type DependentsPolicy struct {
	Children *MinMaxAge `json:"CHI"`
	Parents  *MinMaxAge `json:"PAR"`
	Siblings *MinMaxAge `json:"SIB"`
	Spouse   *MinMaxAge `json:"SPS"`
}

// GetProductRequest ...
type GetProductRequest struct {
	ProductCode string `json:"product_code"`
}

// GetOfferProductRequest ...
type GetOfferProductRequest struct {
	LastName   string      `json:"last_name"`
	FirstName  string      `json:"first_name"`
	MiddleName string      `json:"middle_name"`
	Birthdate  string      `json:"birthdate"`
	Gender     string      `json:"gender"`
	TrxType    json.Number `json:"trx_type"`
	Amount     json.Number `json:"amount"`
}

// CheckActiveProductRequest ...
type CheckActiveProductRequest struct {
	LastName    string `json:"last_name"`
	FirstName   string `json:"first_name"`
	MiddleName  string `json:"middle_name"`
	Birthdate   string `json:"birthdate"`
	Gender      string `json:"gender"`
	ProductCode string `json:"product_code"`
}

// GetOfferProductResult ...
type GetOfferProductResult struct {
	Code    string        `json:"code"`
	Message string        `json:"message"`
	Result  *OfferProduct `json:"result"`
}

// CheckActiveProductResult ...
type CheckActiveProductResult struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Result  *ActiveProduct `json:"result"`
}

// GetProductListResult ...
type GetProductListResult struct {
	Code    string          `json:"code"`
	Message string          `json:"message"`
	Result  []ActiveProduct `json:"result"`
}
