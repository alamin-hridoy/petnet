package handler

import (
	"context"

	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/clientcredentials"
	"google.golang.org/grpc"

	bppb "brank.as/petnet/gunk/drp/v1/bills-payment"
	cicopb "brank.as/petnet/gunk/drp/v1/cashincashout"
	"brank.as/petnet/gunk/drp/v1/dsa"
	mipb "brank.as/petnet/gunk/drp/v1/microinsurance"
	rmpb "brank.as/petnet/gunk/drp/v1/remittance"
	revcom "brank.as/petnet/gunk/drp/v1/revenue-commission"
	tpb "brank.as/petnet/gunk/drp/v1/terminal"
	epb "brank.as/petnet/gunk/dsa/v1/email"
	rat "brank.as/petnet/gunk/dsa/v1/riskassesment"
	upb "brank.as/petnet/gunk/dsa/v1/user"
	bpb "brank.as/petnet/gunk/dsa/v2/branch"
	cspbl "brank.as/petnet/gunk/dsa/v2/cicopartnerlist"
	fpb "brank.as/petnet/gunk/dsa/v2/fees"
	fipb "brank.as/petnet/gunk/dsa/v2/file"
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	ptnrcom "brank.as/petnet/gunk/dsa/v2/partnercommission"
	spbl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	pfpb "brank.as/petnet/gunk/dsa/v2/profile"
	revsrng "brank.as/petnet/gunk/dsa/v2/revenuesharing"
	pfsvc "brank.as/petnet/gunk/dsa/v2/service"
	tmpb "brank.as/petnet/gunk/dsa/v2/temp"
	ttpb "brank.as/petnet/gunk/dsa/v2/transactiontype"
	mpb "brank.as/petnet/gunk/v1/mfa"
	sepb "brank.as/petnet/gunk/v1/session"
	"brank.as/petnet/svcutil/mw"
	rbipb "brank.as/rbac/gunk/v1/invite"
	rbmpb "brank.as/rbac/gunk/v1/mfa"
	rbaoa2pb "brank.as/rbac/gunk/v1/oauth2"
	rbppb "brank.as/rbac/gunk/v1/permissions"
	rbsapb "brank.as/rbac/gunk/v1/serviceaccount"
	rbupb "brank.as/rbac/gunk/v1/user"
)

type Conns struct {
	pfInt       *grpc.ClientConn
	pfIntFwd    *grpc.ClientConn
	idExtSys    *grpc.ClientConn
	idExtFwd    *grpc.ClientConn
	idInt       *grpc.ClientConn
	drpSBIntFwd *grpc.ClientConn
	drpLVIntFwd *grpc.ClientConn
}

func (cs *Conns) GetPfInt() *grpc.ClientConn {
	return cs.pfInt
}

func (cs *Conns) GetIdInt() *grpc.ClientConn {
	return cs.idInt
}

type Cl struct {
	rbacUserAuth partialIdentity
	rbac         identity
	pf           profile
	drpSB        drp
	drpLV        drp
}

// GetProfileCL gets profile client
func (c *Cl) GetProfileCL() profile {
	return c.pf
}

// GetDRPSandboxCL gets drp sandbox client
func (c *Cl) GetDRPSandboxCL() drp {
	return c.drpSB
}

// GetDRPLiveCL gets drp live client
func (c *Cl) GetDRPLiveCL() drp {
	return c.drpLV
}

// todo: making a partial identity for least effort so the UserServiceClient
// can be used with both system auth in identity and user auth in partialIdentity
// refactor later with a better solution
type partialIdentity interface {
	rbupb.UserServiceClient
}

type identity interface {
	rbupb.SignupClient
	rbupb.UserServiceClient
	rbsapb.SvcAccountServiceClient
	rbaoa2pb.AuthClientServiceClient
	rbmpb.MFAAuthServiceClient
	rbmpb.MFAServiceClient
	rbppb.ProductServiceClient
	rbipb.InviteServiceClient
	rbppb.RoleServiceClient
	rbppb.PermissionServiceClient
	rbppb.ValidationServiceClient
}

type profile interface {
	pfpb.OrgProfileServiceClient
	ttpb.TransactionTypeServiceClient
	epb.EmailServiceClient
	fpb.OrgFeesServiceClient
	bpb.BranchServiceClient
	fipb.FileServiceClient
	spb.PartnerServiceClient
	tmpb.EventServiceClient
	upb.SignupServiceClient
	mpb.MFAServiceClient
	upb.UserProfileServiceClient
	sepb.SessionServiceClient
	rat.RiskAssesmentServiceClient
	spbl.PartnerListServiceClient
	pfsvc.ServiceServiceClient
	ptnrcom.PartnerCommissionServiceClient
	cspbl.CICOPartnerListServiceClient
	revsrng.RevenueSharingServiceClient
}

type drp interface {
	tpb.TerminalServiceClient
	revcom.RevenueCommissionServiceClient
	dsa.DSAServiceClient
	rmpb.RemittanceServiceClient
	bppb.BillspaymentServiceClient
	cicopb.CashInCashOutServiceClient
	mipb.MicroInsuranceServiceClient
}

func NewConns(log *logrus.Entry, c *viper.Viper) *Conns {
	u := c.GetString("profile.internal")
	log.WithField("host", u).Info("dialing profile internal")
	pfInt, err := grpc.Dial(
		u,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatal(err)
	}

	log.WithField("host", u).Info("dialing profile internal auth forward")
	pfIntFwd, err := grpc.Dial(
		u,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
			mw.AuthForwarder(),
		),
	)
	if err != nil {
		log.Fatal("unable to connect to profile internal auth forward")
	}

	u = c.GetString("identity.external")
	log.WithField("host", u).Info("dialing identity external system auth")
	idExtSys, err := grpc.Dial(
		u,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
			SystemAuth(c),
		),
	)
	if err != nil {
		log.Fatal("unable to connect to identity external system auth")
	}

	log.WithField("host", u).Info("dialing identity external fwd auth")
	idExtFwd, err := grpc.Dial(
		u,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
			mw.AuthForwarder(),
		),
	)
	if err != nil {
		log.Fatal("unable to connect to identity external fwd auth")
	}

	u = c.GetString("identity.internal")
	log.WithField("host", u).Info("dialing identity internal")
	idInt, err := grpc.Dial(
		u,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	if err != nil {
		log.Fatal("unable to connect to identity internal")
	}

	drpSBIntFwd := &grpc.ClientConn{}
	if !c.GetBool("api.liveMock") {
		u = c.GetString("drp.sandboxInternal")
		log.WithField("host", u).Info("dialing drp sandbox internal")
		drpSBIntFwd, err = grpc.Dial(
			u,
			grpc.WithInsecure(),
			grpc.WithBlock(),
			grpc.WithChainUnaryInterceptor(
				otelgrpc.UnaryClientInterceptor(),
				mw.AuthForwarder(),
			),
		)
		if err != nil {
			log.Fatal("unable to connect to drp sandbox internal")
		}
	}

	u = c.GetString("drp.liveInternal")
	log.WithField("host", u).Info("dialing drp live internal")
	drpLVIntFwd, err := grpc.Dial(
		u,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithChainUnaryInterceptor(
			otelgrpc.UnaryClientInterceptor(),
			mw.AuthForwarder(),
		),
	)
	if err != nil {
		log.Fatal("unable to connect to drp live internal")
	}

	return &Conns{
		pfInt:       pfInt,
		pfIntFwd:    pfIntFwd,
		idExtSys:    idExtSys,
		idExtFwd:    idExtFwd,
		idInt:       idInt,
		drpSBIntFwd: drpSBIntFwd,
		drpLVIntFwd: drpLVIntFwd,
	}
}

func NewSvcClients(cs *Conns) Cl {
	return Cl{
		pf: struct { // All required profile clients
			pfpb.OrgProfileServiceClient
			ttpb.TransactionTypeServiceClient
			epb.EmailServiceClient
			fpb.OrgFeesServiceClient
			bpb.BranchServiceClient
			fipb.FileServiceClient
			spb.PartnerServiceClient
			tmpb.EventServiceClient
			upb.SignupServiceClient
			mpb.MFAServiceClient
			upb.UserProfileServiceClient
			sepb.SessionServiceClient
			rat.RiskAssesmentServiceClient
			spbl.PartnerListServiceClient
			pfsvc.ServiceServiceClient
			ptnrcom.PartnerCommissionServiceClient
			cspbl.CICOPartnerListServiceClient
			revsrng.RevenueSharingServiceClient
		}{
			OrgProfileServiceClient:        pfpb.NewOrgProfileServiceClient(cs.pfInt),
			TransactionTypeServiceClient:   ttpb.NewTransactionTypeServiceClient(cs.pfInt),
			EmailServiceClient:             epb.NewEmailServiceClient(cs.pfInt),
			OrgFeesServiceClient:           fpb.NewOrgFeesServiceClient(cs.pfInt),
			BranchServiceClient:            bpb.NewBranchServiceClient(cs.pfInt),
			FileServiceClient:              fipb.NewFileServiceClient(cs.pfInt),
			PartnerServiceClient:           spb.NewPartnerServiceClient(cs.pfInt),
			EventServiceClient:             tmpb.NewEventServiceClient(cs.pfInt),
			SignupServiceClient:            upb.NewSignupServiceClient(cs.pfInt),
			MFAServiceClient:               mpb.NewMFAServiceClient(cs.pfInt),
			UserProfileServiceClient:       upb.NewUserProfileServiceClient(cs.pfInt),
			SessionServiceClient:           sepb.NewSessionServiceClient(cs.pfInt),
			RiskAssesmentServiceClient:     rat.NewRiskAssesmentServiceClient(cs.pfInt),
			PartnerListServiceClient:       spbl.NewPartnerListServiceClient(cs.pfInt),
			ServiceServiceClient:           pfsvc.NewServiceServiceClient(cs.pfInt),
			PartnerCommissionServiceClient: ptnrcom.NewPartnerCommissionServiceClient(cs.pfInt),
			CICOPartnerListServiceClient:   cspbl.NewCICOPartnerListServiceClient(cs.pfInt),
			RevenueSharingServiceClient:    revsrng.NewRevenueSharingServiceClient(cs.pfInt),
		},
		rbacUserAuth: struct { // All required RBAC clients
			rbupb.UserServiceClient
		}{
			UserServiceClient: rbupb.NewUserServiceClient(cs.idExtFwd),
		},
		rbac: struct { // All required RBAC clients
			rbupb.SignupClient
			rbupb.UserServiceClient
			rbsapb.SvcAccountServiceClient
			rbaoa2pb.AuthClientServiceClient
			rbmpb.MFAAuthServiceClient
			rbmpb.MFAServiceClient
			rbppb.ProductServiceClient
			rbipb.InviteServiceClient
			rbppb.RoleServiceClient
			rbppb.PermissionServiceClient
			rbppb.ValidationServiceClient
		}{
			SignupClient:            rbupb.NewSignupClient(cs.idExtSys),
			UserServiceClient:       rbupb.NewUserServiceClient(cs.idExtSys),
			SvcAccountServiceClient: rbsapb.NewSvcAccountServiceClient(cs.idExtFwd),
			AuthClientServiceClient: rbaoa2pb.NewAuthClientServiceClient(cs.idExtFwd),
			MFAAuthServiceClient:    rbmpb.NewMFAAuthServiceClient(cs.idInt),
			MFAServiceClient:        rbmpb.NewMFAServiceClient(cs.idExtFwd),
			ProductServiceClient:    rbppb.NewProductServiceClient(cs.idExtFwd),
			InviteServiceClient:     rbipb.NewInviteServiceClient(cs.idExtFwd),
			RoleServiceClient:       rbppb.NewRoleServiceClient(cs.idExtFwd),
			PermissionServiceClient: rbppb.NewPermissionServiceClient(cs.idExtFwd),
			ValidationServiceClient: rbppb.NewValidationServiceClient(cs.idInt),
		},
		drpSB: struct { // All required drp clients
			tpb.TerminalServiceClient
			revcom.RevenueCommissionServiceClient
			dsa.DSAServiceClient
			rmpb.RemittanceServiceClient
			bppb.BillspaymentServiceClient
			cicopb.CashInCashOutServiceClient
			mipb.MicroInsuranceServiceClient
		}{
			TerminalServiceClient:          tpb.NewTerminalServiceClient(cs.drpSBIntFwd),
			RevenueCommissionServiceClient: revcom.NewRevenueCommissionServiceClient(cs.drpSBIntFwd),
			DSAServiceClient:               dsa.NewDSAServiceClient(cs.drpSBIntFwd),
			RemittanceServiceClient:        rmpb.NewRemittanceServiceClient(cs.drpSBIntFwd),
			BillspaymentServiceClient:      bppb.NewBillspaymentServiceClient(cs.drpSBIntFwd),
			CashInCashOutServiceClient:     cicopb.NewCashInCashOutServiceClient(cs.drpSBIntFwd),
			MicroInsuranceServiceClient:    mipb.NewMicroInsuranceServiceClient(cs.drpSBIntFwd),
		},
		drpLV: struct {
			// All required drp clients
			tpb.TerminalServiceClient
			revcom.RevenueCommissionServiceClient
			dsa.DSAServiceClient
			rmpb.RemittanceServiceClient
			bppb.BillspaymentServiceClient
			cicopb.CashInCashOutServiceClient
			mipb.MicroInsuranceServiceClient
		}{
			TerminalServiceClient:          tpb.NewTerminalServiceClient(cs.drpLVIntFwd),
			RevenueCommissionServiceClient: revcom.NewRevenueCommissionServiceClient(cs.drpLVIntFwd),
			DSAServiceClient:               dsa.NewDSAServiceClient(cs.drpLVIntFwd),
			RemittanceServiceClient:        rmpb.NewRemittanceServiceClient(cs.drpLVIntFwd),
			BillspaymentServiceClient:      bppb.NewBillspaymentServiceClient(cs.drpLVIntFwd),
			CashInCashOutServiceClient:     cicopb.NewCashInCashOutServiceClient(cs.drpLVIntFwd),
			MicroInsuranceServiceClient:    mipb.NewMicroInsuranceServiceClient(cs.drpLVIntFwd),
		},
	}
}

func SystemAuth(c *viper.Viper) grpc.UnaryClientInterceptor {
	cred := clientcredentials.Config{
		ClientID:     c.GetString("auth.systemClientID"),
		ClientSecret: c.GetString("auth.systemClientSecret"),
		TokenURL:     c.GetString("auth.url") + "/oauth2/token",
		AuthStyle:    oauth2.AuthStyleInHeader,
	}
	ts := cred.TokenSource(context.Background())
	return func(c context.Context, m string, rq, rp interface{}, cc *grpc.ClientConn, inv grpc.UnaryInvoker, o ...grpc.CallOption) error {
		tok, err := ts.Token()
		if err != nil {
			return err
		}
		c = metautils.ExtractOutgoing(c).
			Set("authorization", "Bearer "+tok.AccessToken).ToOutgoing(c)
		return inv(c, m, rq, rp, cc, o...)
	}
}
