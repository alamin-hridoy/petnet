package microinsurance

import (
	"context"
	"encoding/json"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
)

// Relationship ...
type Relationship struct {
	Relationship      string `json:"relationship"`
	RelationshipValue string `json:"relationship_value"`
}

// GetRelationshipsResult ...
type GetRelationshipsResult struct {
	Code    string         `json:"code"`
	Message string         `json:"message"`
	Result  []Relationship `json:"result"`
}

// GetRelationships ...
func (c *Client) GetRelationships(ctx context.Context) ([]Relationship, error) {
	rawRes, err := c.phService.GetMicroInsurance(ctx, c.getUrl("relationships"))
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var res GetRelationshipsResult
	err = json.Unmarshal(rawRes, &res)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, "Invalid perahub response")
	}

	return res.Result, nil
}
