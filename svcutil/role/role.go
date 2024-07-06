package role

type UserRole string

const (
	PFAdmin       UserRole = "Org Admin"
	PFTeamMember  UserRole = "Team Member"
	PFApplication UserRole = "Application"
	PXTAdmin      UserRole = "Proxtera Admin"
	None          UserRole = "None"
)

var (
	ValidUserRoles   = []UserRole{PFTeamMember}
	ValidAdminRoles  = []UserRole{PXTAdmin, PFAdmin, PFApplication}
	ValidMemberRoles = []UserRole{PFTeamMember, PFAdmin, PFApplication}
)
