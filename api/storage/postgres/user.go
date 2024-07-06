package postgres

import (
	"context"

	"brank.as/petnet/api/storage"
)

func (s *Storage) UpsertSession(ctx context.Context, sn storage.Session) (*storage.Session, error) {
	const q = `
INSERT INTO user_session (customer_code, session_data)
    VALUES (:customer_code, :session_data)
ON CONFLICT (customer_code) DO
UPDATE
    SET session_data = :session_data
	 WHERE user_session.customer_code = :customer_code
RETURNING created, updated
`
	stmt, err := s.db.PrepareNamedContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return &sn, stmt.Get(&sn, sn)
}

func (s *Storage) GetSession(ctx context.Context, custCd string) (*storage.Session, error) {
	const q = `select * from user_session where customer_code = $1`
	sn := &storage.Session{}
	return sn, s.db.Get(sn, q, custCd)
}
