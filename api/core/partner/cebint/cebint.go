package cebuanaint

import (
	"brank.as/petnet/api/core/static"
)

func (s *Svc) Kind() string {
	return static.CEBINTCode
}

type Svc struct{}

func New() *Svc {
	return &Svc{}
}
