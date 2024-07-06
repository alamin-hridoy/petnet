package revenue_commission

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
)

// DSA represents perahub DSA object
type DSA struct {
	DsaID          json.Number `json:"id"`
	DsaCode        string      `json:"dsa_code"`
	DsaName        string      `json:"dsa_name"`
	EmailAddress   string      `json:"email_address"`
	Status         json.Number `json:"status"`
	Vatable        json.Number `json:"vatable"`
	Address        string      `json:"address"`
	Tin            string      `json:"tin"`
	UpdatedBy      string      `json:"updated_by"`
	ContactPerson  string      `json:"contact_person"`
	City           string      `json:"city"`
	Province       string      `json:"province"`
	Zipcode        string      `json:"zipcode"`
	President      string      `json:"president"`
	GeneralManager string      `json:"general_manager"`
	CreatedAt      *time.Time  `json:"created_at"`
	UpdatedAt      *time.Time  `json:"updated_at"`
	DeletedAt      *time.Time  `json:"deleted_at"`
}

// SaveDSARequest ...
type SaveDSARequest struct {
	DsaID          uint32      `json:"-"`
	DsaCode        string      `json:"dsa_code"`
	DsaName        string      `json:"dsa_name"`
	EmailAddress   string      `json:"email_address"`
	Vatable        json.Number `json:"vatable"`
	Address        string      `json:"address"`
	Tin            string      `json:"tin"`
	UpdatedBy      string      `json:"updated_by"`
	ContactPerson  string      `json:"contact_person,omitempty"`
	City           string      `json:"city,omitempty"`
	Province       string      `json:"province,omitempty"`
	Zipcode        string      `json:"zipcode,omitempty"`
	President      string      `json:"president,omitempty"`
	GeneralManager string      `json:"general_manager,omitempty"`
}

// DSAResult ...
type DSAResult struct {
	Code    json.Number `json:"code"`
	Message string      `json:"message"`
	Result  *DSA        `json:"result"`
}

// ListDSAResult ...
type ListDSAResult struct {
	Code    json.Number `json:"code"`
	Message string      `json:"message"`
	Result  []DSA       `json:"result"`
}

// GetDSAByID returns DSA object found by id
func (c *Client) GetDSAByID(ctx context.Context, dsaID uint32) (*DSA, error) {
	if dsaID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "dsa id required")
	}

	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl(fmt.Sprintf("dsa/%d", dsaID)))

	return toDSAResponse(rawRes, err)
}

// ListDSA ...
func (c *Client) ListDSA(ctx context.Context) ([]DSA, error) {
	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl("dsa"))
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var resp ListDSAResult
	err = json.Unmarshal(rawRes, &resp)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}

// CreateDSA ...
func (c *Client) CreateDSA(ctx context.Context, req *SaveDSARequest) (*DSA, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	rawRes, err := c.phService.PostRevComm(ctx, c.getUrl("dsa"), *req)

	return toDSAResponse(rawRes, err)
}

// UpdateDSA ...
func (c *Client) UpdateDSA(ctx context.Context, req *SaveDSARequest) (*DSA, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	if req.DsaID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "dsa id required")
	}

	rawRes, err := c.phService.PutRevComm(ctx, c.getUrl(fmt.Sprintf("dsa/%d", req.DsaID)), *req)

	return toDSAResponse(rawRes, err)
}

// DeleteDSA ...
func (c *Client) DeleteDSA(ctx context.Context, dsaID uint32) (*DSA, error) {
	if dsaID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "dsa id required")
	}

	rawRes, err := c.phService.DeleteRevComm(ctx, c.getUrl(fmt.Sprintf("dsa/%d", dsaID)))
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return toDSAResponse(rawRes, err)
}

func toDSAResponse(jsonRaw json.RawMessage, err error) (*DSA, error) {
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var resp DSAResult
	err = json.Unmarshal(jsonRaw, &resp)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}
