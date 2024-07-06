package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

const createRole = `
INSERT INTO roles (
    id,
    org_id,
    role_name,
    description,
    updateduid,
    create_user_id
) VALUES (
    :id,
    :org_id,
    :role_name,
    :description,
    :updateduid,
    :create_user_id
) RETURNING
    created,updated
`

func (s *Storage) CreateRole(ctx context.Context, p storage.Role) (*storage.Role, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, createRole)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing organization insert: %w", err)
	}
	return &p, nil
}

const getRole = `
SELECT *
FROM roles
WHERE id = $1
`

func (s *Storage) GetRole(ctx context.Context, id string) (*storage.Role, error) {
	var p storage.Role
	if err := s.db.Get(&p, getRole, id); err != nil {
		return nil, fmt.Errorf("executing role details: %w", err)
	}
	return &p, nil
}

const deleteRole = `
UPDATE roles
SET
(
    delete_user_id,
    deleted
) = (
    :delete_user_id,
    :deleted
) WHERE id = :id
RETURNING *
`

func (s *Storage) DeleteRole(ctx context.Context, p storage.Role) (*storage.Role, error) {
	p.Delete = sql.NullTime{
		Valid: true,
		Time:  time.Now().UTC(),
	}

	stmt, err := s.db.PrepareNamedContext(ctx, deleteRole)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, &p); err != nil {
		return nil, fmt.Errorf("executing role delete: %w", err)
	}
	return &p, nil
}

const updateRole = `
UPDATE roles SET
	role_name = :role_name,
	updateduid = :updateduid,
    description = :description
	WHERE id = :id
	RETURNING updated
`

func (s *Storage) UpdateRole(ctx context.Context, p storage.Role) (*storage.Role, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateRole)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&p, p); err != nil {
		return nil, fmt.Errorf("executing role update: %w", err)
	}
	return &p, nil
}

func (s *Storage) ListRole(ctx context.Context, f core.ListRoleFilter) ([]storage.Role, error) {
	var roles []storage.Role
	SortBy := "ASC"
	filterQL := []string{}
	filterQ := ""
	if f.SortBy != "" {
		SortBy = f.SortBy
	}
	if f.Name != "" {
		filterQL = append(filterQL, fmt.Sprintf("role_name = '%s'", f.Name))
	}
	if f.OrgID != "" {
		filterQL = append(filterQL, fmt.Sprintf("org_id = '%s'", f.OrgID))
	}
	if len(f.ID) > 0 {
		filterQL = append(filterQL, fmt.Sprintf("id IN('%s')", strings.Join(f.ID, ",")))
	}

	if len(filterQL) > 0 {
		filterQ = "WHERE " + strings.Join(filterQL, " AND ")
	}

	limit := ""
	if f.Limit > 0 {
		limit = fmt.Sprintf(" LIMIT NULLIF(%d, 0) OFFSET %d;", f.Limit, f.Offset)
	}
	listRole := fmt.Sprintf("WITH cnt AS (select count(*) as count FROM roles %s) SELECT p.*, cnt.count FROM roles as p left join cnt on true  %s ORDER BY role_name %s", filterQ, filterQ, SortBy)
	fullQuery := listRole + limit
	if err := s.db.Select(&roles, fullQuery); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return roles, nil
}
