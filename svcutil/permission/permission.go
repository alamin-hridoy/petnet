package permission

const (
	AdminServiceName  = "Petnet Admin Permissions"
	SuperUserRoleName = "Petnet SuperUser Role"
)

const (
	// actions
	CreateAct    = "create"
	ReadAct      = "read"
	UpdateAct    = "update"
	DeleteAct    = "delete"
	SuperviseAct = "supervise"
)

const (
	// resources
	DSAListDetailRes = "dsaListDetail"
	TransactionRes   = "transaction"
	DocChecklistRes  = "documentChecklist"
	RiskAssRes       = "riskAssessment"
	CurrencyRes      = "currency"
	FeesRes          = "fees"
	BranchRes        = "branch"
	SvcCatalogRes    = "serviceCatalog"
)

type Permission struct {
	Description string
	Resource    string
	Actions     []string
}

type RBACPermission struct {
	Resource string
	Actions  []string
}

var Permissions = map[string]Permission{
	"dsaListAndDetails": {
		Description: "DSA List and Detail Management",
		Resource:    DSAListDetailRes,
		Actions:     []string{CreateAct, ReadAct, UpdateAct, DeleteAct, SuperviseAct},
	},
	"transactionListAndDetails": {
		Description: "Transactions Management",
		Resource:    TransactionRes,
		Actions:     []string{CreateAct, ReadAct, UpdateAct, DeleteAct, SuperviseAct},
	},
	"documentChecklist": {
		Description: "DocumentChecklist Management",
		Resource:    DocChecklistRes,
		Actions:     []string{CreateAct, ReadAct, UpdateAct, DeleteAct, SuperviseAct},
	},
	"riskAssessment": {
		Description: "RiskAssessment Management",
		Resource:    RiskAssRes,
		Actions:     []string{CreateAct, ReadAct, UpdateAct, DeleteAct, SuperviseAct},
	},
	"currency": {
		Description: "Currency Management",
		Resource:    CurrencyRes,
		Actions:     []string{CreateAct, ReadAct, UpdateAct, DeleteAct, SuperviseAct},
	},
	"feeAndCommission": {
		Description: "Fee and Comission Management",
		Resource:    FeesRes,
		Actions:     []string{CreateAct, ReadAct, UpdateAct, DeleteAct, SuperviseAct},
	},
	"locationAndBranches": {
		Description: "Location and Branches Management",
		Resource:    BranchRes,
		Actions:     []string{CreateAct, ReadAct, UpdateAct, DeleteAct, SuperviseAct},
	},
	"serviceCatalog": {
		Description: "Service Catalog Management",
		Resource:    SvcCatalogRes,
		Actions:     []string{CreateAct, ReadAct, UpdateAct, DeleteAct, SuperviseAct},
	},
}
