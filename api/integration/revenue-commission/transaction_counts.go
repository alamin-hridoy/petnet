package revenue_commission

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"google.golang.org/grpc/codes"

	coreerror "brank.as/petnet/api/core/error"
)

// TransactionCount ...
type TransactionCount struct {
	ID                json.Number    `json:"id"`
	DsaCode           string         `json:"dsa_code"`
	YearMonth         string         `json:"year_month"`
	RemittanceCount   json.Number    `json:"remittance_count"`
	CiCoCount         json.Number    `json:"cico_count"`
	BillsPaymentCount json.Number    `json:"bills_payment_count"`
	InsuranceCount    json.Number    `json:"insurance_count"`
	UpdatedBy         string         `json:"updated_by"`
	DsaCommission     json.Number    `json:"dsa_commission"`
	DsaCommissionType CommissionType `json:"dsa_commission_type"`
	CreatedAt         *time.Time     `json:"created_at"`
	UpdatedAt         *time.Time     `json:"updated_at"`
}

// TransactionCountResult ...
type TransactionCountResult struct {
	Code    json.Number       `json:"code"`
	Message string            `json:"message"`
	Result  *TransactionCount `json:"result"`
}

// ListTransactionCountResult ...
type ListTransactionCountResult struct {
	Code    json.Number        `json:"code"`
	Message string             `json:"message"`
	Result  []TransactionCount `json:"result"`
}

// SaveTransactionCountRequest ...
type SaveTransactionCountRequest struct {
	ID                uint32         `json:"-"`
	DsaCode           string         `json:"dsa_code"`
	YearMonth         string         `json:"year_month"`
	RemittanceCount   json.Number    `json:"remittance_count"`
	CiCoCount         json.Number    `json:"cico_count"`
	BillsPaymentCount json.Number    `json:"bills_payment_count"`
	InsuranceCount    json.Number    `json:"insurance_count"`
	UpdatedBy         string         `json:"updated_by"`
	DsaCommission     json.Number    `json:"dsa_commission"`
	DsaCommissionType CommissionType `json:"dsa_commission_type"`
}

// GetTransactionCountByID returns transaction count object found by id
func (c *Client) GetTransactionCountByID(ctx context.Context, trxCntID uint32) (*TransactionCount, error) {
	if trxCntID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "transaction count id required")
	}

	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl(fmt.Sprintf("dsa-trx-count/%d", trxCntID)))

	return toTransactionCountResponse(rawRes, err)
}

// ListAllTransactionCount ...
func (c *Client) ListAllTransactionCount(ctx context.Context) ([]TransactionCount, error) {
	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl("dsa-trx-count"))
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	return toListTransactionCountResponse(rawRes, err)
}

// ListTransactionCountByYearMonth list transaction counts by yearMonth.
// yearMonth = [yyyy][mm]. Ex. Jan 2022 will be "202201"
func (c *Client) ListTransactionCountByYearMonth(ctx context.Context, yearMonth string) ([]TransactionCount, error) {
	if strings.TrimSpace(yearMonth) == "" {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "year month string required")
	}

	rawRes, err := c.phService.GetRevComm(ctx, c.getUrl("dsa-trx-count/year-month"+yearMonth))
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	return toListTransactionCountResponse(rawRes, err)
}

// CreateTransactionCount ...
func (c *Client) CreateTransactionCount(ctx context.Context, req *SaveTransactionCountRequest) (*TransactionCount, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	rawRes, err := c.phService.PostRevComm(ctx, c.getUrl("dsa-trx-count"), *req)

	return toTransactionCountResponse(rawRes, err)
}

// UpdateTransactionCount ...
func (c *Client) UpdateTransactionCount(ctx context.Context, req *SaveRemcoCommissionRequest) (*TransactionCount, error) {
	if req == nil {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "empty request")
	}

	if req.ID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "transaction count id required")
	}

	rawRes, err := c.phService.PutRevComm(ctx, c.getUrl(fmt.Sprintf("dsa-trx-count/%d", req.ID)), *req)

	return toTransactionCountResponse(rawRes, err)
}

// DeleteTransactionCount ...
func (c *Client) DeleteTransactionCount(ctx context.Context, trxCntID uint32) (*TransactionCount, error) {
	if trxCntID == 0 {
		return nil, coreerror.NewCoreError(codes.InvalidArgument, "transaction count id required")
	}

	rawRes, err := c.phService.DeleteRevComm(ctx, c.getUrl(fmt.Sprintf("dsa-trx-count/%d", trxCntID)))
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return toTransactionCountResponse(rawRes, err)
}

func toTransactionCountResponse(jsonRaw json.RawMessage, err error) (*TransactionCount, error) {
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var resp TransactionCountResult
	err = json.Unmarshal(jsonRaw, &resp)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}

func toListTransactionCountResponse(jsonRaw json.RawMessage, err error) ([]TransactionCount, error) {
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	var resp ListTransactionCountResult
	err = json.Unmarshal(jsonRaw, &resp)
	if err != nil {
		return nil, coreerror.NewCoreError(codes.Internal, coreerror.MsgDRPInternalError)
	}

	return resp.Result, nil
}
