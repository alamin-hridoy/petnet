package phmw

import (
	"context"

	"brank.as/petnet/serviceutil/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/svcutil/mw/meta"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	rev "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	trxtp "brank.as/petnet/gunk/dsa/v2/transactiontype"
	sapb "brank.as/rbac/gunk/v1/serviceaccount"
)

const (
	Sandbox    = "sandbox"
	Production = "production"
	Live       = "live"
	Provider   = "provider"
	Petnet     = "petnet"
	DSA        = "dsa"
)

type SvcAcct struct {
	cl  sapb.ValidationServiceClient
	trx trxtp.TransactionTypeServiceClient
	pf  ppb.OrgProfileServiceClient
	log *logrus.Entry
}

func NewServiceAccount(cl sapb.ValidationServiceClient, trx trxtp.TransactionTypeServiceClient, dpf ppb.OrgProfileServiceClient, log *logrus.Entry) *SvcAcct {
	return &SvcAcct{
		cl:  cl,
		trx: trx,
		pf:  dpf,
		log: log,
	}
}

func Reset() meta.MetaFunc {
	return func(ctx context.Context) (context.Context, error) {
		return metautils.ExtractIncoming(ctx).
			Del(UsrName).
			Del(WUNo).
			Del(BDay).
			Del(dsa).
			Del(DSAOrgID).
			Del(hydra.ClientIDKey).
			Del(hydra.OrgIDKey).
			ToIncoming(ctx), nil
	}
}

func ConfirmDSA(env string) meta.MetaFunc {
	return func(ctx context.Context) (context.Context, error) {
		if GetDSA(ctx) == "" {
			return ctx, status.Error(codes.Unauthenticated, "DSA not authenticated")
		}
		if GetEnv(ctx) != env {
			return ctx, status.Error(codes.PermissionDenied, "token not permitted for environment: "+env)
		}
		if hydra.OrgID(ctx) == "" {
			ctx = metautils.ExtractIncoming(ctx).Set(hydra.OrgIDKey, GetDSA(ctx)).ToIncoming(ctx)
		}
		return ctx, nil
	}
}

// Metadata for service account validation
func (s *SvcAcct) Metadata(ctx context.Context) (context.Context, error) {
	trxTp := ""

	if GetEnv(ctx) != "" {
		return ctx, nil
	}
	clientID := GetUserID(ctx)
	v, err := s.cl.ValidateAccount(ctx, &sapb.ValidateAccountRequest{
		ClientID: clientID,
	})
	if err != nil {
		return ctx, err
	}

	environment := Sandbox
	if clientID != "" {
		trxdetls, _ := s.trx.GetTransactionTypeByClientId(ctx, &trxtp.GetTransactionTypeByClientIdRequest{
			ClientID: clientID,
		})
		if trxdetls != nil {
			trxTp = trxdetls.GetTransactionType()
			environment = trxdetls.GetEnvironment()
		}
	}

	// added this check to convert environment from production to live
	// hydra service account has either live or sandbox environment
	// which doesn't match with production environment set on the api_key_transaction_type table
	if environment == Production {
		environment = Live
	}

	org := s.getOrgProfile(ctx, trxTp, v.GetOrgID())

	return metautils.ExtractIncoming(ctx).
		Add(dsa, v.GetOrgID()).
		Add(env, v.GetEnvironment()).
		Add(TransactionType, trxTp).
		Add(TerminalID, org.TerminalID).
		Add(OrgInfo, org.OrgType).
		Add(DsaCode, org.DSACode).
		Add(apiEnv, environment).
		Add(DSAOrgID, v.GetOrgID()).
		Add(UsrName, v.GetClientName()).ToIncoming(ctx), nil
}

type Org struct {
	DSACode    string
	TerminalID string
	OrgType    string
	Partner    string
}

func (s *SvcAcct) getOrgProfile(ctx context.Context, trxIp, orgID string) *Org {
	if orgID == "" {
		s.log.Error("empty orgID for user " + GetUserID(ctx))
		return nil
	}

	dsp, err := s.pf.GetProfile(ctx, &ppb.GetProfileRequest{
		OrgID: orgID,
	})
	if err != nil {
		logging.WithError(err, s.log).Error("get profile for service account error")
		return nil
	}

	if dsp.GetProfile() == nil {
		s.log.Error("empty profile for orgID " + orgID)
		return nil
	}

	terminalID := dsp.GetProfile().TerminalIdOtc
	if trxIp == rev.TransactionType_DIGITAL.String() {
		terminalID = dsp.GetProfile().TerminalIdDigital
	}

	return &Org{
		DSACode:    dsp.GetProfile().DsaCode,
		TerminalID: terminalID,
		OrgType: func() string {
			if ppb.OrgType_PetNet == dsp.GetProfile().OrgType {
				return Petnet
			}

			if ppb.OrgType_DSA == dsp.GetProfile().OrgType {
				if dsp.GetProfile().IsProvider {
					return Provider
				}
				return DSA
			}
			return ""
		}(),
		Partner: dsp.GetProfile().Partner,
	}
}
