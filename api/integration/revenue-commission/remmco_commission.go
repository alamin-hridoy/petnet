package revenue_commission

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
)

// RemcoCommission ...
type RemcoCommission struct {
	ID                 json.Number    `json:"id"`
	RemcoID            json.Number    `json:"remco_id"`
	CommissionType     CommissionType `json:"commission_type"`
	CommissionValue    json.Number    `json:"commission_value"`
	CommissionValueOtc json.Number    `json:"commission_value_otc"`
	CurrencyID         string         `json:"currency_id"`
	TrxType            TrxType        `json:"trx_type"`
	UpdatedBy          string         `json:"updated_by"`
	Corridor           string         `json:"corridor"`
	CreatedAt          *time.Time     `json:"created_at"`
	UpdatedAt          *time.Time     `json:"updated_at"`
}

// RemcoCommissionResult ...
type RemcoCommissionResult struct {
	Code    json.Number      `json:"code"`
	Message string           `json:"message"`
	Result  *RemcoCommission `json:"result"`
}

// ListRemcoCommissionResult ...
type ListRemcoCommissionResult struct {
	Code    json.Number       `json:"code"`
	Message string            `json:"message"`
	Result  []RemcoCommission `json:"result"`
}

// SaveRemcoCommissionRequest ...
type SaveRemcoCommissionRequest struct {
	ID                 uint32         `json:"-"`
	RemcoID            json.Number    `json:"remco_id"`
	CommissionType     CommissionType `json:"commission_type"`
	CommissionValue    json.Number    `json:"commission_value"`
	CommissionValueOtc json.Number    `json:"commission_value_otc"`
	CurrencyID         string         `json:"currency_id"`
	TrxType            TrxType        `json:"trx_type"`
	UpdatedBy          string         `json:"updated_by"`
	Corridor           string         `json:"corridor"`
}

// GetRemcoCommissionByID returns remco commission object found by id
func (c *Client) GetRemcoCommissionByID(ctx context.Context, commissionID uint32) (*RemcoCommission, error) {
	if commissionID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "remco commission id required")
	}

	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl(fmt.Sprintf("remco-comm/%d", commissionID)))

	return toRemcoCommissionResponse(rawRes, err)
}

// ListRemcoCommission ...
func (c *Client) ListRemcoCommission(ctx context.Context) ([]RemcoCommission, error) {
	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl("remco-comm"))
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var resp ListRemcoCommissionResult
	err = json.Unmarshal(rawRes, &resp)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}

// CreateRemcoCommission ...
func (c *Client) CreateRemcoCommission(ctx context.Context, req *SaveRemcoCommissionRequest) (*RemcoCommission, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	rawRes, err := c.phService.PostRevComm(ctx, c.getUrl("remco-comm"), *req)

	return toRemcoCommissionResponse(rawRes, err)
}

// UpdateRemcoCommission ...
func (c *Client) UpdateRemcoCommission(ctx context.Context, req *SaveRemcoCommissionRequest) (*RemcoCommission, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	if req.ID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "remco commission id required")
	}

	rawRes, err := c.phService.PutRevComm(ctx, c.getUrl(fmt.Sprintf("remco-comm/%d", req.ID)), *req)

	return toRemcoCommissionResponse(rawRes, err)
}

// DeleteRemcoCommission ...
func (c *Client) DeleteRemcoCommission(ctx context.Context, commissionID uint32) (*RemcoCommission, error) {
	if commissionID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "remco commission id required")
	}

	rawRes, err := c.phService.DeleteRevComm(ctx, c.getUrl(fmt.Sprintf("remco-comm/%d", commissionID)))
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return toRemcoCommissionResponse(rawRes, err)
}

func toRemcoCommissionResponse(jsonRaw json.RawMessage, err error) (*RemcoCommission, error) {
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var resp RemcoCommissionResult
	err = json.Unmarshal(jsonRaw, &resp)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}
