package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/petnet/api/storage"
	"brank.as/petnet/serviceutil/logging"
)

const createInputGuide = `
INSERT INTO inputguide (
	partner,
	inputguide
) VALUES (
	:partner,
	:inputguide
) RETURNING
created,updated
`

func (s *Storage) CreateInputGuide(ctx context.Context, r storage.InputGuide) (*storage.InputGuide, error) {
	log := logging.FromContext(ctx)
	log.WithField("inputguide", r).Trace("storing")

	stmt, err := s.db.PrepareNamedContext(ctx, createInputGuide)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing inputguide insert: %w", err)
	}
	return &r, nil
}

const updateInputGuide = `
	UPDATE inputguide
	SET
		inputguide= :inputguide
	WHERE partner= :partner
	RETURNING updated, created`

func (s *Storage) UpdateInputGuide(ctx context.Context, r storage.InputGuide) (*storage.InputGuide, error) {
	stmt, err := s.db.PrepareNamedContext(ctx, updateInputGuide)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	if err := stmt.Get(&r, r); err != nil {
		return nil, fmt.Errorf("executing inputguide update: %w", err)
	}
	return &r, nil
}

const getInputGuide = `
SELECT *
FROM inputguide
WHERE partner= $1
`

func (s *Storage) GetInputGuide(ctx context.Context, partner string) (*storage.InputGuide, error) {
	var r storage.InputGuide
	if err := s.db.Get(&r, getInputGuide, partner); err != nil {
		if err != nil {
			if err == sql.ErrNoRows {
				return nil, storage.ErrNotFound
			}
			return nil, err
		}
	}
	return &r, nil
}
