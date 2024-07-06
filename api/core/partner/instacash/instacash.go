package instacash

import (
	"brank.as/petnet/api/core/static"
)

func (s *Svc) Kind() string {
	return static.ICCode
}

type Svc struct{}

func New() *Svc {
	return &Svc{}
}
