package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/serviceutil/logging"
	"github.com/lib/pq"
)

// UpsertSession creates or update a session and returns the ID.
func (s *Storage) UpsertSession(ctx context.Context, sess *storage.Session) (string, error) {
	const upsertSession = `
INSERT INTO session as sess (
    user_id,
	 session_expiry
) VALUES (
	 :user_id,
	 :session_expiry
)
ON CONFLICT (user_id) DO UPDATE
SET
	session_expiry= :session_expiry,
	deleted= COALESCE(:deleted, sess.deleted)
WHERE
	sess.user_id= :user_id
RETURNING
    id,created,updated
`

	log := logging.FromContext(ctx)
	pstmt, err := s.db.PrepareNamedContext(ctx, upsertSession)
	if err != nil {
		logging.WithError(err, log).Error("upsert session")
		return "", err
	}
	defer pstmt.Close()
	if err := pstmt.Get(sess, sess); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return "", storage.Conflict
		}
		return "", fmt.Errorf("executing session upsert: %w", err)
	}
	return sess.ID, nil
}

// GetSession return profile matched against user ID
func (s *Storage) GetSession(ctx context.Context, uid string) (*storage.Session, error) {
	const getSession = `SELECT * FROM session WHERE user_id = $1`
	var sess storage.Session
	if err := s.db.Get(&sess, getSession, uid); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return &sess, nil
}
