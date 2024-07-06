package partnerlist

import (
	"brank.as/petnet/profile/storage/postgres"
)

type Svc struct {
	st *postgres.Storage
}

func New(st *postgres.Storage) *Svc {
	return &Svc{st: st}
}
