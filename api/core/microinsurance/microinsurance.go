package microinsurance

import (
	"brank.as/petnet/api/integration/microinsurance"
	"brank.as/petnet/api/storage/postgres"
)

// MICoreSvc micro insurance core service
type MICoreSvc struct {
	cl      *microinsurance.Client
	storage *postgres.Storage
}

// NewMicroInsuranceCoreSvc ...
func NewMicroInsuranceCoreSvc(st *postgres.Storage, cl *microinsurance.Client) *MICoreSvc {
	return &MICoreSvc{
		cl:      cl,
		storage: st,
	}
}
