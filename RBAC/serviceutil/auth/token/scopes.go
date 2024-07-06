package token

// Scope defines an authorization scope for use in the token service.
type Scope = string

const (
	Admin Scope = "brankas"

	OrgSvc     Scope = "brankas.v2.org.OrgService"
	OrgUserSvc Scope = "brankas.v2.org.UserService"

	DataOrgUserSvc       Scope = "brankas.v1.data.orguser.OrgUserService"
	DataOrgCredentialSvc Scope = "brankas.v1.data.orgcredential.OrgCredentialService"

	AccountSvc      Scope = "brankas.v2.account.AccountService"
	BankSvc         Scope = "brankas.v2.bank.BankService"
	MatchSvc        Scope = "brankas.v2.match.MatchService"
	TransactionSvc  Scope = "brankas.v2.transaction.TransactionService"
	DisbursementSvc Scope = "brankas.v2.disbursement.DisbursementService"

	APIServer             Scope = "brankas.v2.server.api"
	APIServerFastCheckout Scope = "brankas.v2.server.api.TransferService"
	AuthServer            Scope = "brankas.v2.server.auth"
	DisbursementServer    Scope = "brankas.v2.server.disbursement"
)

// MetaTag defines known entities in the context metadata.
type MetaTag = string

const (
	MetaOrgID  MetaTag = "org-id"
	MetaUserID MetaTag = "user-id"
)
