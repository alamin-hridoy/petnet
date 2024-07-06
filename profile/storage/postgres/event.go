package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"brank.as/petnet/profile/storage"
	"github.com/lib/pq"
)

// CreateEventData ...
func (s *Storage) CreateEventData(ctx context.Context, e *storage.EventData) error {
	const insertEvent = `
INSERT INTO temp_event_data (
    event_id,
	 resource,
	 action, 
	 data
) VALUES (
	 :event_id,
	 :resource,
	 :action,
	 :data
)
`
	pstmt, err := s.db.PrepareNamedContext(ctx, insertEvent)
	if err != nil {
		return err
	}
	defer pstmt.Close()
	if _, err := pstmt.Exec(e); err != nil {
		pErr, ok := err.(*pq.Error)
		if ok && pErr.Code == pqUnique {
			return storage.Conflict
		}
		return fmt.Errorf("executing event data insert: %w", err)
	}
	return nil
}

// GetEventData return event data matched against eventID
func (s *Storage) GetEventData(ctx context.Context, eid string) (*storage.EventData, error) {
	const eventSelect = `
SELECT *
FROM temp_event_data 
WHERE
	event_id = :event_id AND deleted IS NULL
`
	e := &storage.EventData{}
	stmt, err := s.db.PrepareNamed(eventSelect)
	if err != nil {
		return nil, fmt.Errorf("preparing named query GetEventData: %w", err)
	}
	arg := map[string]interface{}{
		"event_id": eid,
	}
	if err := stmt.Get(e, arg); err != nil {
		if err == sql.ErrNoRows {
			return nil, storage.NotFound
		}
		return nil, err
	}
	return e, nil
}

// DeleteEventData deletes given event data by eventID
func (s *Storage) DeleteEventData(ctx context.Context, eid string) error {
	const serviceUpdate = `
DELETE
FROM temp_event_data
WHERE
	event_id = :event_id
`
	stmt, err := s.db.PrepareNamedContext(ctx, serviceUpdate)
	if err != nil {
		return err
	}
	defer stmt.Close()
	arg := map[string]interface{}{
		"event_id": eid,
	}
	if _, err := stmt.Exec(arg); err != nil {
		if err == sql.ErrNoRows {
			return storage.NotFound
		}
		return fmt.Errorf("executing event data delete: %w", err)
	}
	return nil
}
