package perahub

import (
	"context"
	"encoding/json"
)

type WISERefreshRecipientReq struct {
	Email             string            `json:"email"`
	Currency          string            `json:"currency"`
	Type              string            `json:"type"`
	OwnedByCustomer   bool              `json:"ownedByCustomer"`
	AccountHolderName string            `json:"accountHolderName"`
	Requirements      []WISERequirement `json:"requirements"`
}

type WISERefreshRecipientResp struct {
	Requirements []WISERequirementsResp `json:"requirements"`
}

type WISERequirementsResp struct {
	Type      string      `json:"type"`
	Title     string      `json:"title"`
	UsageInfo string      `json:"usageInfo"`
	Fields    []WISEField `json:"fields"`
}

type WISEField struct {
	Name  string      `json:"name"`
	Group []WISEGroup `json:"group"`
}

type WISEGroup struct {
	Key                string              `json:"key"`
	Name               string              `json:"name"`
	Type               string              `json:"type"`
	RefreshReqOnChange bool                `json:"refreshRequirementsOnChange"`
	Required           bool                `json:"required"`
	DisplayFormat      string              `json:"displayFormat"`
	Example            string              `json:"example"`
	MinLength          json.Number         `json:"minLength"`
	MaxLength          json.Number         `json:"maxLength"`
	ValidationRegexp   string              `json:"validationRegexp"`
	ValidationAsync    WISEValidationAsync `json:"validationAsync"`
	ValuesAllowed      []WISEValueAllowed  `json:"valuesAllowed"`
}

type WISEValidationAsync struct {
	URL    string      `json:"url"`
	Params []WISEParam `json:"params"`
}

type WISEParam struct {
	Key       string `json:"key"`
	ParamName string `json:"parameterName"`
	Required  bool   `json:"required"`
}

type WISEValueAllowed struct {
	Key  string `json:"key"`
	Name string `json:"name"`
}

func (s *Svc) WISERefreshRecipient(ctx context.Context, req WISERefreshRecipientReq) (*WISERefreshRecipientResp, error) {
	res, err := s.postNonex(ctx, s.nonexURL("transferwise/recipients/refresh-onchange"), req)
	if err != nil {
		return nil, err
	}

	rb := &WISERefreshRecipientResp{}
	if err := json.Unmarshal(res, rb); err != nil {
		return nil, err
	}
	return rb, nil
}
