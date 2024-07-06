package bills_payment

import (
	"context"
	"encoding/json"
)

type BCBillerInfoRequest struct {
	Code       string `json:"code"`
	UserID     string `json:"user_id"`
	LocationID string `json:"location_id"`
}

type BCBillerInfoResponse struct {
	Code    int                `json:"code"`
	Message string             `json:"message"`
	Result  BCBillerInfoResult `json:"result"`
	RemcoID int                `json:"remco_id"`
}

type BCBillerInfoResult struct {
	Code            string     `json:"code"`
	IsCde           int        `json:"isCde"`
	IsAsync         int        `json:"isAsync"`
	Name            string     `json:"name"`
	Description     string     `json:"description"`
	Logo            string     `json:"logo"`
	Category        string     `json:"category"`
	Type            string     `json:"type"`
	IsMultipleBills int        `json:"isMultipleBills"`
	Parameters      Parameters `json:"parameters"`
}

type Parameters struct {
	Verify   []Verify   `json:"verify"`
	Transact []Transact `json:"transact"`
}

type Verify struct {
	ReferenceNumber ReferenceNumber `json:"referenceNumber"`
}

type ReferenceNumber struct {
	Description string          `json:"description"`
	Rules       BillerInfoRules `json:"rules"`
	Label       string          `json:"label"`
}

type BillerInfoRules struct {
	Digits8  BCCM `json:"digits:8"`
	Required BCCM `json:"required"`
}

type Transact struct {
	ClientReference ClientReference `json:"clientReference"`
}
type ClientReference struct {
	Description string `json:"description"`
	Rules       CRules `json:"rules"`
	Label       string `json:"label"`
}

type CRules struct {
	AlphaDash BCCM `json:"alpha_dash"`
	Required  BCCM `json:"required"`
	UniqueCrn BCCM `json:"unique_crn"`
}

type BCCM struct {
	Message string `json:"message"`
	Code    int    `json:"code"`
}

func (c *Client) BCBillerInfo(ctx context.Context, req BCBillerInfoRequest) (*BCBillerInfoResponse, error) {
	res, err := c.phService.BillsPost(ctx, c.getUrl("bayad/bayad-center/biller-info"), req)
	if err != nil {
		return nil, err
	}

	bcbi := &BCBillerInfoResponse{}
	if err := json.Unmarshal(res, bcbi); err != nil {
		return nil, err
	}
	return bcbi, nil
}
