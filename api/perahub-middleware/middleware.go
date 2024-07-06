package phmw

import (
	"context"

	"brank.as/petnet/api/core"
	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/svcutil/mw/meta"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"

	authpb "brank.as/rbac/gunk/v1/authenticate"
)

var _ meta.RequestMeta = (*Svc)(nil)

type Svc struct{}

func New() *Svc { return &Svc{} }

const (
	FName           = "fname"
	LName           = "lname"
	BDay            = "bday"
	ValidID         = "id_valid"
	WUNo            = "wu_no"
	Ntly            = "nation"
	LoyaltyNo       = "loyalty_no"
	UsrName         = "username"
	UserID          = hydra.ClientIDKey
	Coy             = "coy"
	OperatorID      = "operator-id"
	TerminalID      = "terminal-id"
	DSAOrgID        = "dsa-org-id"
	OrgType         = "org-type"
	OrgInfo         = "org-info"
	OrgID           = "org-id"
	Partner         = "partner"
	TransactionType = "transaction-type"
	DsaCode         = "dsa_code"

	env    = "environment"
	apiEnv = "api-environment"

	// Indicates the owner of the auth client used for authentication.
	// Loaded by RBAC IDP.
	dsa = "owner"
)

// Metadata fulfills the metaloader interface
func (s *Svc) Metadata(ctx context.Context) (context.Context, error) {
	ex := hydra.GetExtra(ctx)
	md := metautils.ExtractIncoming(ctx)
	for _, k := range []string{FName, LName, BDay, ValidID, WUNo, Ntly, LoyaltyNo, dsa, DSAOrgID, OrgID} {
		v := ex[k]
		if v != "" {
			md.Set(k, v)
		}
	}
	return md.ToIncoming(ctx), nil
}

func UserSession(ctx context.Context, u core.User) (*authpb.Session, error) {
	return &authpb.Session{
		UserID: u.CustNo,
		Session: map[string]string{
			FName:     u.FirstName,
			LName:     u.LastName,
			BDay:      u.Birthdate,
			ValidID:   u.ValidIdnt,
			WUNo:      u.WUCardNo,
			LoyaltyNo: u.LoyaltyCardNo,
		},
		OpenID: map[string]string{
			FName: u.FirstName,
			LName: u.LastName,
			BDay:  u.Birthdate,
		},
	}, nil
}

// helper funcs
func GetUserID(ctx context.Context) string     { return metautils.ExtractIncoming(ctx).Get(UserID) }
func GetUsrName(ctx context.Context) string    { return metautils.ExtractIncoming(ctx).Get(UsrName) }
func GetWUNo(ctx context.Context) string       { return metautils.ExtractIncoming(ctx).Get(WUNo) }
func GetBDay(ctx context.Context) string       { return metautils.ExtractIncoming(ctx).Get(BDay) }
func GetDSA(ctx context.Context) string        { return metautils.ExtractIncoming(ctx).Get(dsa) }
func GetEnv(ctx context.Context) string        { return metautils.ExtractIncoming(ctx).Get(env) }
func GetAPIEnv(ctx context.Context) string     { return metautils.ExtractIncoming(ctx).Get(apiEnv) }
func GetCoy(ctx context.Context) string        { return metautils.ExtractIncoming(ctx).Get(Coy) }
func GetOperatorID(ctx context.Context) string { return metautils.ExtractIncoming(ctx).Get(OperatorID) }
func GetTerminalID(ctx context.Context) string { return metautils.ExtractIncoming(ctx).Get(TerminalID) }
func GetDSAOrgID(ctx context.Context) string   { return metautils.ExtractIncoming(ctx).Get(DSAOrgID) }
func GetOrgInfo(ctx context.Context) string    { return metautils.ExtractIncoming(ctx).Get(OrgInfo) }
func GetOrgType(ctx context.Context) string    { return metautils.ExtractIncoming(ctx).Get(OrgType) }
func GetPartner(ctx context.Context) string    { return metautils.ExtractIncoming(ctx).Get(Partner) }
func GetTransactionTypes(ctx context.Context) string {
	return metautils.ExtractIncoming(ctx).Get(TransactionType)
}
func GetDsaCode(ctx context.Context) string { return metautils.ExtractIncoming(ctx).Get(DsaCode) }
