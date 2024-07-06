package ussc

import (
	"brank.as/petnet/api/core/static"
)

func (s *Svc) Kind() string {
	return static.USSCCode
}

type Svc struct{}

func New() *Svc {
	return &Svc{}
}
