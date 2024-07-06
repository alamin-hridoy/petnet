package perahub

import (
	"context"
	"encoding/json"

	"brank.as/petnet/api/core/static"
	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type SDQsRequest struct {
	SDQType string `json:"sdq_type"`
	Remco   string `json:"remco"`
}

type SDQsResponseBody struct {
	ID           []SDQID
	Occupation   []SDQOccupation
	Position     []SDQPosition
	Purpose      []SDQPurpose
	Relationship []SDQRelationship
	SourceOfFund []SDQSourceOfFund
}

type SDQID struct {
	Index         string `json:"index"`
	TemplateValue string `json:"template_value"`
	DocumentType  string `json:"document_type"`
	DocDescEng    string `json:"document_desc_eng"`
	DocDescFil    string `json:"document_desc_fil"`
	HasIssueDate  string `json:"hasIssueDate"`
	HasExpiration string `json:"hasExpiration"`
}

type SDQOccupation struct {
	Occupation      string `json:"occupation"`
	OccupationValue string `json:"occupation_value"`
}

type SDQPosition struct {
	Position      string `json:"position"`
	PositionValue string `json:"position_value"`
}

type SDQPurpose struct {
	Purpose      string `json:"purpose"`
	PurposeValue string `json:"purpose_value"`
}

type SDQRelationship struct {
	Relationship      string `json:"relationship"`
	RelationshipValue string `json:"relationship_value"`
}

type SDQSourceOfFund struct {
	SourceOfFund      string `json:"source_of_funds"`
	SourceOfFundValue string `json:"source_of_funds_value"`
}

func (s *Svc) SDQs(ctx context.Context, r SDQsRequest) (*SDQsResponseBody, error) {
	if err := validation.ValidateStruct(&r,
		validation.Field(&r.SDQType, validation.Required, validation.In("id", "occupation", "position", "purpose", "relationship", "source_of_fund", "all")),
		validation.Field(&r.Remco, validation.Required, validation.In(static.WUCode)),
	); err != nil {
		return nil, err
	}

	ts := []string{r.SDQType}
	if r.SDQType == "all" {
		ts = []string{"id", "occupation", "position", "purpose", "relationship", "source_of_fund"}
	}
	sdqsRes := &SDQsResponseBody{}
	const mod, modReq = "sdq", "sdq"
	for _, t := range ts {
		r.SDQType = t
		req, err := s.newParahubRequest(ctx, mod, modReq, r)
		if err != nil {
			return nil, err
		}
		resp, err := s.post(ctx, s.moduleURL(mod, modReq), *req)
		if err != nil {
			return nil, err
		}

		switch t {
		case "id":
			if err := json.Unmarshal(resp, &sdqsRes.ID); err != nil {
				return nil, err
			}
		case "occupation":
			if err := json.Unmarshal(resp, &sdqsRes.Occupation); err != nil {
				return nil, err
			}
		case "position":
			if err := json.Unmarshal(resp, &sdqsRes.Position); err != nil {
				return nil, err
			}
		case "purpose":
			if err := json.Unmarshal(resp, &sdqsRes.Purpose); err != nil {
				return nil, err
			}
		case "relationship":
			if err := json.Unmarshal(resp, &sdqsRes.Relationship); err != nil {
				return nil, err
			}
		case "source_of_fund":
			if err := json.Unmarshal(resp, &sdqsRes.SourceOfFund); err != nil {
				return nil, err
			}
		}
	}
	return sdqsRes, nil
}
