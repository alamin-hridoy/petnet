package japanremit

import (
	"brank.as/petnet/api/core/static"
)

func (s *Svc) Kind() string {
	return static.JPRCode
}

type Svc struct{}

func New() *Svc {
	return &Svc{}
}
