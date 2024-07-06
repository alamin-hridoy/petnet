package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
)

const insertTransactionTypeApi = `
INSERT INTO api_key_transaction_type (
    org_id,
    user_id,
    client_id,
	environment,
	transaction_type
) VALUES (
	 :org_id,
	 :user_id,
	 :client_id,
	 :environment,
	 :transaction_type
) RETURNING
    id
`

func (s *Storage) InsertApiKeyTransactionType(ctx context.Context, pf *storage.ApiKeyTransactionType) (string, error) {
	log := logging.FromContext(ctx)
	switch {
	case pf.OrgID == "":
		return "", fmt.Errorf("org id cannot be empty")
	case pf.UserID == "":
		return "", fmt.Errorf("user id cannot be empty")
	case pf.ClientID == "":
		return "", fmt.Errorf("client id cannot be empty")
	case pf.Environment == "":
		return "", fmt.Errorf("environment cannot be empty")
	case pf.TransactionType == "":
		return "", fmt.Errorf("transaction type cannot be empty")
	}
	pstmt, err := s.db.PrepareNamedContext(ctx, insertTransactionTypeApi)
	if err != nil {
		logging.WithError(err, log).Error("insert transaction type api")
		return "", err
	}
	defer pstmt.Close()
	if err := pstmt.Get(pf, pf); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return "", storage.Conflict
		}
		return "", fmt.Errorf("executing transaction type api insert: %w", err)
	}
	return pf.UserID, nil
}

func (s *Storage) GetAPITransactionType(ctx context.Context, pf *storage.ApiKeyTransactionType) (*storage.ApiKeyTransactionType, error) {
	const getTransactionTypeApi = `SELECT * FROM api_key_transaction_type WHERE org_id = $1 AND user_id = $2 AND environment = $3 AND transaction_type = $4 AND deleted IS NULL`
	switch {
	case pf.OrgID == "":
		return nil, fmt.Errorf("org id cannot be empty")
	case pf.UserID == "":
		return nil, fmt.Errorf("user id cannot be empty")
	case pf.Environment == "":
		return nil, fmt.Errorf("environment cannot be empty")
	case pf.TransactionType == "":
		return nil, fmt.Errorf("transaction type cannot be empty")
	}
	var pfs storage.ApiKeyTransactionType
	if err := s.db.Get(&pfs, getTransactionTypeApi, pf.OrgID, pf.UserID, pf.Environment, pf.TransactionType); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &pfs, nil
}

func (s *Storage) ListUserAPIKeyTransactionType(ctx context.Context, oid, userid string) ([]storage.ApiKeyTransactionType, error) {
	const transactionTypeApiID = `SELECT * FROM api_key_transaction_type WHERE org_id = $1 and user_id = $2 AND deleted IS NULL`
	switch {
	case oid == "":
		return nil, fmt.Errorf("org id cannot be empty")
	case userid == "":
		return nil, fmt.Errorf("user id cannot be empty")
	}
	var pf []storage.ApiKeyTransactionType
	if err := s.db.Select(&pf, transactionTypeApiID, oid, userid); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return pf, nil
}

func (s *Storage) GetTransactionTypeByClientId(ctx context.Context, clintid string) (*storage.GetTransactionTypeByClientIdResponse, error) {
	switch {
	case clintid == "":
		return nil, fmt.Errorf("client id cannot be empty")
	}
	const getTransactionTypeByClientId = `SELECT environment, transaction_type FROM api_key_transaction_type WHERE client_id = $1 AND deleted IS NULL`
	var pfs storage.GetTransactionTypeByClientIdResponse
	if err := s.db.Get(&pfs, getTransactionTypeByClientId, clintid); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &pfs, nil
}
