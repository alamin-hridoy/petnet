package user

import (
	"context"
	"log"

	"brank.as/petnet/api/core"
	ceb "brank.as/petnet/api/core/user/cebuana"
	wum "brank.as/petnet/api/core/user/wise"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage/postgres"
	ppb "brank.as/petnet/gunk/drp/v1/profile"
)

type UserManager interface {
	RegisterUser(ctx context.Context, req core.RegisterUserReq) (*core.RegisterUserResp, error)
	CreateProfile(ctx context.Context, req core.CreateProfileReq) (*core.CreateProfileResp, error)
	GetProfile(ctx context.Context, req core.GetProfileReq) (*core.GetProfileResp, error)
	GetUser(ctx context.Context, req core.GetUserRequest) (*core.GetUserResponse, error)
	CreateRecipient(ctx context.Context, req *ppb.CreateRecipientRequest) (*ppb.CreateRecipientResponse, error)
	GetRecipients(ctx context.Context, req *ppb.GetRecipientsRequest) (*ppb.GetRecipientsResponse, error)
	RefreshRecipient(ctx context.Context, req *ppb.UpdateRecipientRequest) (*ppb.UpdateRecipientResponse, error)
	DeleteRecipient(ctx context.Context, req *ppb.DeleteRecipientRequest) (*ppb.DeleteRecipientResponse, error)
	Kind() string
}

type Svc struct {
	usermanagers map[string]UserManager
	st           *postgres.Storage
	ph           *perahub.Svc
}

func New(st *postgres.Storage, ph *perahub.Svc) *Svc {
	gs := []UserManager{wum.New(ph), ceb.New(ph)}
	s := &Svc{
		usermanagers: make(map[string]UserManager, len(gs)),
		st:           st,
		ph:           ph,
	}
	for i, r := range gs {
		switch {
		case r == nil:
			log.Fatalf("user manager %d nil", i)
		case r.Kind() == "":
			log.Fatalf("user manager %d missing partner type", i)
		}
		s.usermanagers[r.Kind()] = r
	}
	return s
}
