package svcacct

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/rbac/serviceutil/logging"
	client "brank.as/rbac/svcutil/hydraclient"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/core"
	"brank.as/rbac/usermgm/storage"
)

const APIKeyPrefix = 8

// CreateSvcAccount register hydra client and record in storage.
func (h *Svc) CreateSvcAccount(ctx context.Context, sa storage.SvcAccount) (id, secret string, err error) {
	log := logging.FromContext(ctx).WithField("method", "hydra.CreateSvcAccount")
	log.WithField("svcacct", sa).Debug("processing")
	switch sa.AuthType {
	case storage.OAuth2:
		id, secret, err = h.createOAuth(ctx, sa)
	case storage.APIKey:
		id, secret, err = h.createApiKey(ctx, sa)
	default:
		return "", "", status.Error(codes.InvalidArgument, "invalid auth type")
	}
	if err != nil {
		return "", "", err
	}
	if sa.Role != "" {
		if _, err := h.asn.AssignRole(ctx, core.Grant{
			RoleID:  sa.Role,
			GrantID: id,
		}); err != nil {
			logging.WithError(err, log).Error("assigning role")
		}
	}
	return id, secret, nil
}

func (h *Svc) createOAuth(ctx context.Context, sa storage.SvcAccount) (id, secret string, err error) {
	log := logging.FromContext(ctx).WithField("method", "hydra.CreateOAuth2SvcClient")

	bufID, err := random.String(24)
	if err != nil {
		logging.WithError(err, log).Error("random id")
		return "", "", status.Error(codes.Internal, "failed to generate service account")
	}
	ac := client.AuthClient{
		OwnerID:    sa.OrgID,
		ClientID:   bufID,
		GrantTypes: []string{"client_credentials"},
		AuthMethod: "client_secret_basic",
	}
	const bsClientName = "BootstrapClient"
	if sa.ClientName == bsClientName && h.bsHydraClient != "" && h.bsHydraSecret != "" {
		// this is to have a static bootstrap client for local development
		ac.ClientID = h.bsHydraClient
		ac.Secret = h.bsHydraSecret
	} else {
		bufSec, err := random.String(48)
		if err != nil {
			logging.WithError(err, log).Error("random secret")
			return "", "", status.Error(codes.Internal, "failed to generate service account")
		}
		ac.Secret = bufSec
	}
	cl, err := h.cl.CreateClient(ctx, ac)
	if err != nil {
		logging.WithError(err, log).Error("create in hydra")
		return "", "", status.Error(codes.Internal, "failed to generate service account")
	}
	sa.ClientID = cl.ClientID
	clID, err := h.store.CreateSvcAccount(ctx, sa)
	if err != nil {
		logging.WithError(err, log).Error("db store new svcacct")
		return "", "", err
	}
	if clID != cl.ClientID {
		logging.WithError(err, log).Error("store returned incorrect id")
		return "", "", status.Error(codes.DataLoss, "incorrect client id from store")
	}
	return cl.ClientID, cl.Secret, nil
}

func (h *Svc) createApiKey(ctx context.Context, sa storage.SvcAccount) (id, secret string, err error) {
	log := logging.FromContext(ctx).WithField("method", "hydra.CreateAPIKey")
	key, err := random.String(48)
	if err != nil {
		logging.WithError(err, log).Error("random id")
		return "", "", status.Error(codes.Internal, "failed to generate service account")
	}

	sa.ClientID = key[:APIKeyPrefix]
	sa.Challenge = key
	clID, err := h.store.CreateSvcAccount(ctx, sa)
	if err != nil {
		logging.WithError(err, log).Error("db store new svcacct")
		return "", "", err
	}
	if clID != sa.ClientID {
		logging.WithError(err, log).Error("store returned incorrect id")
		return "", "", status.Error(codes.DataLoss, "incorrect client id from store")
	}
	return sa.ClientID, key[:APIKeyPrefix] + "." + key[APIKeyPrefix:], nil
}
