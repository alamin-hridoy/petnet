package metrobank

import (
	"brank.as/petnet/api/core/static"
)

func (s *Svc) Kind() string {
	return static.MBCode
}

type Svc struct{}

func New() *Svc {
	return &Svc{}
}
