package intelexpress

import (
	"brank.as/petnet/api/core/static"
)

func (s *Svc) Kind() string {
	return static.IECode
}

type Svc struct{}

func New() *Svc {
	return &Svc{}
}
