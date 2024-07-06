package mfauth

import (
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/spf13/viper"

	"brank.as/rbac/usermgm/storage/postgres"
)

type Svc struct {
	svcName string
	dur     time.Duration
	st      *postgres.Storage
	em      MFAEmailer
}

type MFAEmailer interface {
	EmailMFA(email, code string) error
}

func New(config *viper.Viper, st *postgres.Storage, em MFAEmailer) (*Svc, error) {
	s := &Svc{
		svcName: config.GetString("project.mfaissuer"),
		dur:     config.GetDuration("project.mfatimeout"),
		st:      st,
		em:      em,
	}
	if err := validation.ValidateStruct(s,
		validation.Field(&s.svcName, validation.Required, validation.Length(3, 0)),
		validation.Field(&s.dur, validation.Min(time.Minute), validation.Max(4*time.Hour)),
	); err != nil {
		return nil, err
	}
	return s, nil
}
