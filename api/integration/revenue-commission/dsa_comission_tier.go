package revenue_commission

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
)

// CommissionTier ...
type CommissionTier struct {
	TierID    json.Number `json:"id"`
	TierNo    string      `json:"tier_no"`
	Minimum   json.Number `json:"minimum"`
	Maximum   json.Number `json:"maximum"`
	UpdatedBy string      `json:"updated_by"`
	CreatedAt *time.Time  `json:"created_at"`
	UpdatedAt *time.Time  `json:"updated_at"`
}

// CommissionTierResult ...
type CommissionTierResult struct {
	Code    json.Number     `json:"code"`
	Message string          `json:"message"`
	Result  *CommissionTier `json:"result"`
}

// ListCommissionTierResult ...
type ListCommissionTierResult struct {
	Code    json.Number      `json:"code"`
	Message string           `json:"message"`
	Result  []CommissionTier `json:"result"`
}

// SaveCommissionTierRequest ...
type SaveCommissionTierRequest struct {
	TierID    uint32      `json:"-"`
	TierNo    string      `json:"tier_no"`
	Minimum   json.Number `json:"minimum"`
	Maximum   json.Number `json:"maximum"`
	UpdatedBy string      `json:"updated_by"`
}

// GetCommissionTierByID returns Commission Tier object found by id
func (c *Client) GetCommissionTierByID(ctx context.Context, tierID uint32) (*CommissionTier, error) {
	if tierID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "tier id required")
	}

	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl(fmt.Sprintf("dsa-comm-tier/%d", tierID)))

	return toCommissionTierResponse(rawRes, err)
}

// ListCommissionTier ...
func (c *Client) ListCommissionTier(ctx context.Context) ([]CommissionTier, error) {
	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl("dsa-comm-tier"))
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var resp ListCommissionTierResult
	err = json.Unmarshal(rawRes, &resp)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}

// CreateCommissionTier ...
func (c *Client) CreateCommissionTier(ctx context.Context, req *SaveCommissionTierRequest) (*CommissionTier, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	rawRes, err := c.phService.PostRevComm(ctx, c.getUrl("dsa-comm-tier"), *req)

	return toCommissionTierResponse(rawRes, err)
}

// UpdateCommissionTier ...
func (c *Client) UpdateCommissionTier(ctx context.Context, req *SaveCommissionTierRequest) (*CommissionTier, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	if req.TierID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "tier id required")
	}

	rawRes, err := c.phService.PutRevComm(ctx, c.getUrl(fmt.Sprintf("dsa-comm-tier/%d", req.TierID)), *req)

	return toCommissionTierResponse(rawRes, err)
}

// DeleteCommissionTier ...
func (c *Client) DeleteCommissionTier(ctx context.Context, tierID uint32) (*CommissionTier, error) {
	if tierID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "tier id required")
	}

	rawRes, err := c.phService.DeleteRevComm(ctx, c.getUrl(fmt.Sprintf("dsa-comm-tier/%d", tierID)))
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return toCommissionTierResponse(rawRes, err)
}

func toCommissionTierResponse(jsonRaw json.RawMessage, err error) (*CommissionTier, error) {
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var resp CommissionTierResult
	err = json.Unmarshal(jsonRaw, &resp)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}
