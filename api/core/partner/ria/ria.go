package ria

import (
	"brank.as/petnet/api/core/static"
)

func (s *Svc) Kind() string {
	return static.RIACode
}

type Svc struct{}

func New() *Svc {
	return &Svc{}
}
