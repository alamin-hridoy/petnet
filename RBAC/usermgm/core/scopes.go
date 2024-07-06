package core

import "time"

type Scope struct {
	ID      string
	Name    string
	Group   string
	Desc    string
	Updated time.Time
}

type OfferGrant struct {
	OrgID   string
	OrgName string
	Skip    bool
	Scopes  map[string]ScopeGroup
}

type ScopeGroup struct {
	Name    string
	Desc    string
	Scopes  []Scope
	Updated time.Time
}

type ConsentGrant struct {
	ID        string
	UserID    string
	ClientID  string
	OwnerID   string
	Scopes    []string
	Timestamp time.Time
}
