package iremit

import (
	"brank.as/petnet/api/core/static"
)

func (s *Svc) Kind() string {
	return static.IRCode
}

type Svc struct{}

func New() *Svc {
	return &Svc{}
}
