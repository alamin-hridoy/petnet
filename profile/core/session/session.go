package session

import (
	"context"
	"database/sql"
	"time"

	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
)

type Svc struct {
	st *postgres.Storage
}

func New(st *postgres.Storage) *Svc {
	return &Svc{st: st}
}

type UpsertSessionReq struct {
	Type    int
	ID      string
	Expiry  sql.NullTime
	Deleted sql.NullTime
}

type GetSessionReq struct {
	Type int
	ID   string
}

type GetSessionRes struct {
	Expired bool
	Expiry  time.Time
}

// UpsertSession ...
func (s *Svc) UpsertSession(ctx context.Context, req *UpsertSessionReq) (string, error) {
	uid := req.ID
	if req.Type == 0 {
		pf, err := s.st.GetUserProfileByEmail(ctx, req.ID)
		if err != nil {
			return "", err
		}
		uid = pf.UserID
	}

	id, err := s.st.UpsertSession(ctx, &storage.Session{
		UserID: uid,
		Expiry: req.Expiry,
	})
	if err != nil {
		return "", err
	}
	return id, nil
}

// UpsertSession ...
func (s *Svc) GetSession(ctx context.Context, req *GetSessionReq) (*GetSessionRes, error) {
	uid := req.ID
	if req.Type == 0 {
		pf, err := s.st.GetUserProfileByEmail(ctx, req.ID)
		if err != nil {
			return nil, err
		}
		uid = pf.UserID
	}

	sess, err := s.st.GetSession(ctx, uid)
	if err != nil {
		return nil, err
	}

	return &GetSessionRes{
		Expired: time.Until(sess.Expiry.Time) <= 0,
		Expiry:  sess.Expiry.Time,
	}, nil
}
