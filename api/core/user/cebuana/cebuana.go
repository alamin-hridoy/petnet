package cebuana

import (
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/integration/perahub"
)

func (s *Svc) Kind() string {
	return static.CEBCode
}

type Svc struct {
	ph *perahub.Svc
}

func New(ph *perahub.Svc) *Svc {
	return &Svc{
		ph: ph,
	}
}