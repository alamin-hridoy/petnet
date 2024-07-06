package bpi

import (
	"brank.as/petnet/api/core/static"
)

func (s *Svc) Kind() string {
	return static.BPICode
}

type Svc struct{}

func New() *Svc {
	return &Svc{}
}
