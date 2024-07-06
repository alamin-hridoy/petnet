package core

import (
	"time"

	"github.com/bojanz/currency"
)

type SearchRemitResponse struct {
	Total        int
	SearchRemits []SearchRemit
}

type SearchRemit struct {
	DsaID      string
	UserID     string
	OrgID      string
	PtnrUserID string
	DeviceID   string

	RemitPartner  string
	RemitType     string
	PtnrRemitType string
	DestCurrency  string
	ControlNo     string

	Remitter    Contact
	Receiver    Contact
	SentCountry string
	DestCity    string
	DestState   string

	Status  string
	Message string

	RemitAmount    currency.Minor
	DisburseAmount currency.Minor
	Taxes          map[string]currency.Minor
	Tax            currency.Minor
	Charges        map[string]currency.Minor
	Charge         currency.Minor

	OtherInfo []byte

	TxnStagedTime    time.Time
	TxnCompletedTime time.Time
}
