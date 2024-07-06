package revenue_commission

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
)

// DSACommission ...
type DSACommission struct {
	ID                 json.Number    `json:"id"`
	DsaCode            string         `json:"dsa_code"`
	CommissionType     CommissionType `json:"commission_type"`
	TierID             json.Number    `json:"tier"`
	CommissionAmount   json.Number    `json:"commission_amount"`
	CommissionCurrency json.Number    `json:"commission_currency"`
	TrxType            string         `json:"dsa_trx_type"`
	RemitType          string         `json:"dsa_remit_type"`
	UpdatedBy          string         `json:"updated_by"`
	CreatedAt          *time.Time     `json:"created_at"`
	UpdatedAt          *time.Time     `json:"updated_at"`
	EffectiveDate      string         `json:"effective_date"`
}

// DSACommissionResult ...
type DSACommissionResult struct {
	Code    json.Number    `json:"code"`
	Message string         `json:"message"`
	Result  *DSACommission `json:"result"`
}

// ListDSACommissionResult ...
type ListDSACommissionResult struct {
	Code    json.Number     `json:"code"`
	Message string          `json:"message"`
	Result  []DSACommission `json:"result"`
}

// SaveDSACommissionRequest ...
type SaveDSACommissionRequest struct {
	ID                 uint32         `json:"-"`
	DsaCode            string         `json:"dsa_code"`
	CommissionType     CommissionType `json:"commission_type"`
	TierID             uint32         `json:"tier"`
	CommissionAmount   float64        `json:"commission_amount"`
	CommissionCurrency string         `json:"commission_currency"`
	UpdatedBy          string         `json:"updated_by"`
	EffectiveDate      string         `json:"effective_date"`
	TrxType            string         `json:"dsa_trx_type"`
	RemitType          string         `json:"dsa_remit_type"`
}

// GetDSACommissionByID returns DSA Commission object found by id
func (c *Client) GetDSACommissionByID(ctx context.Context, commissionID uint32) (*DSACommission, error) {
	if commissionID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "commission id required")
	}

	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl(fmt.Sprintf("dsa-comm/%d", commissionID)))

	return toDSACommissionResponse(rawRes, err)
}

// ListDSACommission ...
func (c *Client) ListDSACommission(ctx context.Context) ([]DSACommission, error) {
	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl("dsa-comm"))
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var resp ListDSACommissionResult
	err = json.Unmarshal(rawRes, &resp)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}

// CreateDSACommission ...
func (c *Client) CreateDSACommission(ctx context.Context, req *SaveDSACommissionRequest) (*DSACommission, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	rawRes, err := c.phService.PostRevComm(ctx, c.getUrl("dsa-comm"), *req)

	return toDSACommissionResponse(rawRes, err)
}

// UpdateDSACommission ...
func (c *Client) UpdateDSACommission(ctx context.Context, req *SaveDSACommissionRequest) (*DSACommission, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	if req.ID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "commission id required")
	}

	rawRes, err := c.phService.PutRevComm(ctx, c.getUrl(fmt.Sprintf("dsa-comm/%d", req.ID)), *req)

	return toDSACommissionResponse(rawRes, err)
}

// DeleteDSACommission ...
func (c *Client) DeleteDSACommission(ctx context.Context, commissionID uint32) (*DSACommission, error) {
	if commissionID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "commission id required")
	}

	rawRes, err := c.phService.DeleteRevComm(ctx, c.getUrl(fmt.Sprintf("dsa-comm/%d", commissionID)))
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return toDSACommissionResponse(rawRes, err)
}

func toDSACommissionResponse(jsonRaw json.RawMessage, err error) (*DSACommission, error) {
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var resp DSACommissionResult
	err = json.Unmarshal(jsonRaw, &resp)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}
