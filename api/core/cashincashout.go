package core

import (
	"time"
)

type CICOTransactListResponse struct {
	Total             int
	CICOTransactLists []CICOTransactList
}

type CICOTransactList struct {
	// DsaID      string
	// UserID     string
	OrgID       string
	PartnerCode string
	TrxProvider string

	TrxType          string
	ReferenceNumber  string
	PetnetTrackingno string
	PrincipalAmount  int
	Charges          int
	TotalAmount      int

	TrxStatus string

	TrxDate          time.Time
	TxnStagedTime    time.Time
	TxnCompletedTime time.Time
}

type CICOFilterList struct {
	From             time.Time
	Until            time.Time
	Limit            int
	Offset           int
	SortOrder        string
	SortByColumn     string
	ReferenceNumber  string
	ExcludeProviders []string
	OrgID            string
}
