package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
	"github.com/sirupsen/logrus"

	"brank.as/rbac/serviceutil/logging"
	"brank.as/rbac/usermgm/storage"
)

const upsertScopeQry = `
INSERT INTO scopes (
	id,
	name,
	group_name,
	description
) VALUES (
	:id,
	:name,
	:group_name,
	:description
) ON CONFLICT (id) DO
UPDATE SET (
	name,
	group_name,
	description
) = (
	:name,
	:group_name,
	:description
) RETURNING *;
`

const upsertScopeGroupQry = `INSERT INTO scope_group (name) VALUES ($1)
ON CONFLICT (name) DO NOTHING`

// UpsertScope storage.
func (s *Storage) UpsertScope(ctx context.Context, sc storage.Scope) (*storage.Scope, error) {
	q, err := s.db.PrepareNamedContext(ctx, upsertScopeQry)
	if err != nil {
		return nil, err
	}
	defer q.Close()

	if err := q.Get(&sc, sc); err != nil {
		return nil, err
	}
	if _, err := s.db.ExecContext(ctx, upsertScopeGroupQry, sc.Group); err != nil {
		logging.FromContext(ctx).WithFields(logrus.Fields{
			"method": "storage.upsertscope",
			"scope":  sc.ID,
			"group":  sc.Group,
		}).WithError(err).Error("failed to store group")
	}
	return &sc, nil
}

func (s *Storage) GetScopes(ctx context.Context, sc []string) ([]storage.Scope, error) {
	const getQry = `SELECT * FROM scopes WHERE id IN (?)`
	q, a, err := sqlx.In(getQry, sc)
	if err != nil {
		return nil, err
	}
	sp := []storage.Scope{}
	if err := s.db.SelectContext(ctx, &sp, s.db.Rebind(q), a...); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return sp, nil
}

func (s *Storage) UpdateGroup(ctx context.Context, sg storage.ScopeGroup) (*storage.ScopeGroup, error) {
	const updateScopeGroupQry = `UPDATE scope_group
SET (description, updated) = ($1, now()) WHERE name=$2 RETURNING *`
	return &sg, s.db.GetContext(ctx, &sg, updateScopeGroupQry, sg.Desc, sg.Name)
}

func (s *Storage) GetScopeGroups(ctx context.Context, gr []string) ([]storage.ScopeGroup, error) {
	const getQry = `SELECT * FROM scope_group WHERE name IN (?)`
	q, a, err := sqlx.In(getQry, gr)
	if err != nil {
		return nil, err
	}
	sg := []storage.ScopeGroup{}
	if err := s.db.SelectContext(ctx, &sg, s.db.Rebind(q), a...); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return sg, nil
}

const recordQry = `INSERT INTO consent_grant (
	user_id,
	client_id,
	owner_id,
	scopes
) VALUES (
	:user_id,
	:client_id,
	:owner_id,
	:scopes
)
RETURNING *`

type consent struct {
	ID        string         `db:"grant_id"`
	UserID    string         `db:"user_id"`
	ClientID  string         `db:"client_id"`
	OwnerID   string         `db:"owner_id"`
	Scopes    pq.StringArray `db:"scopes"`
	Timestamp time.Time      `db:"timestamp"`
}

func (s *Storage) RecordGrant(ctx context.Context, g storage.ConsentGrant) (*storage.ConsentGrant, error) {
	gr := consent{
		UserID:   g.UserID,
		ClientID: g.ClientID,
		OwnerID:  g.OwnerID,
		Scopes:   g.Scopes,
	}
	q, err := s.db.PrepareNamedContext(ctx, recordQry)
	if err != nil {
		return nil, err
	}
	defer q.Close()
	if err := q.Get(&gr, gr); err != nil {
		return nil, err
	}
	g.ID = gr.ID
	g.Timestamp = gr.Timestamp
	return &g, nil
}
