package oauthclient

import (
	"context"
	"time"

	client "brank.as/rbac/svcutil/hydraclient"
	"brank.as/rbac/usermgm/storage"
)

type Svc struct {
	hy *client.AdminClient
	st Store
}

type Store interface {
	CreateOauthClient(context.Context, storage.OAuthClient) (*storage.OAuthClient, error)
	GetOauthClientByID(context.Context, string) (*storage.OAuthClient, error)
	GetOauthClientByOrgID(context.Context, string, bool) ([]storage.OAuthClient, error)
	UpdateOauthClient(context.Context, storage.OAuthClient) (*storage.OAuthClient, error)
	DeleteOauthClient(context.Context, storage.OAuthClient) (*time.Time, error)
}

func New(st Store, cl *client.AdminClient) *Svc {
	return &Svc{
		hy: cl,
		st: st,
	}
}
