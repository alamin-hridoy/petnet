package auth

import (
	"context"

	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage"
)

type SessionStore interface {
	UpsertSession(context.Context, storage.Session) (*storage.Session, error)
	GetSession(context.Context, string) (*storage.Session, error)
}

type Svc struct {
	p  *perahub.Svc
	st SessionStore
}

func New(p *perahub.Svc, st SessionStore) *Svc {
	return &Svc{p: p, st: st}
}
