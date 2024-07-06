package core

type UserSession struct {
	UserID        string
	OrgID         string
	Session       map[string]string
	OpenIDSession map[string]string
}
