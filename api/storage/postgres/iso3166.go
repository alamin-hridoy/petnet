package postgres

import (
	"context"

	"brank.as/petnet/api/storage"
)

func (s *Storage) GetISO(ctx context.Context, c string) (*storage.ISOCty, error) {
	const getISO = `select * from iso3166 where Code = $1`
	i := storage.ISOCty{}
	if err := s.db.Get(&i, getISO, c); err != nil {
		return nil, err
	}
	return &i, nil
}
