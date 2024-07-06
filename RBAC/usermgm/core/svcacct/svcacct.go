package svcacct

import (
	"context"
	"time"

	"github.com/spf13/viper"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	client "brank.as/rbac/svcutil/hydraclient"

	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

type Svc struct {
	cl            *client.AdminClient
	store         SvcAccountStore
	bsHydraClient string
	bsHydraSecret string
	asn           RoleAssigner
}

type SvcAccountStore interface {
	CreateSvcAccount(context.Context, storage.SvcAccount) (string, error)
	DisableSvcAccount(context.Context, storage.SvcAccount) (*time.Time, error)
	ValidateSvcAccount(ctx context.Context, id, key string) (*storage.SvcAccount, error)
	GetSvcAccountByID(context.Context, string) (*storage.SvcAccount, error)
	GetOrgByID(context.Context, string) (*storage.Organization, error)
	GetRole(context.Context, string) (*storage.Role, error)
}

type RoleAssigner interface {
	AssignRole(ctx context.Context, g core.Grant) (*core.Role, error)
}

// New live hydra integration.
func New(conf *viper.Viper, cl *client.AdminClient, store SvcAccountStore, asn RoleAssigner) *Svc {
	hs := conf.GetString("bootstrap.hydraSecret")
	hc := conf.GetString("bootstrap.hydraClient")
	return &Svc{cl: cl, store: store, bsHydraClient: hc, bsHydraSecret: hs, asn: asn}
}

// DisableSvcAccount removes client from hydra and records disable record in storage.
func (h *Svc) DisableSvcAccount(ctx context.Context, sa storage.SvcAccount) (*time.Time, error) {
	log := logging.FromContext(ctx).WithField("method", "hydra.DisableSvcAccount")
	if err := h.cl.DeleteClient(ctx, sa.ClientID); err != nil {
		logging.WithError(err, log).Error("failed to delete in hydra")
		return nil, status.Error(codes.Internal, "failed to disable account")
	}
	tm, err := h.store.DisableSvcAccount(ctx, sa)
	if err != nil {
		// Disable in hydra will stop all authenticated calls.
		// Will need to fix the storage separately.
		logging.WithError(err, log).Error("failed to disable in storage")
		t := time.Now()
		return &t, nil
	}
	return tm, nil
}
