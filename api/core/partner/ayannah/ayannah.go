package ayannah

import (
	"brank.as/petnet/api/core/static"
)

func (s *Svc) Kind() string {
	return static.AYACode
}

type Svc struct{}

func New() *Svc {
	return &Svc{}
}
