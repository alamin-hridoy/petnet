package md

import (
	"context"

	"brank.as/petnet/serviceutil/auth/hydra"
	"brank.as/petnet/svcutil/mw"
	"brank.as/petnet/svcutil/mw/meta"
	mta "github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc"
)

func OrgForwarder() grpc.UnaryClientInterceptor {
	return func(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn,
		invoker grpc.UnaryInvoker, opts ...grpc.CallOption,
	) error {
		if oid := mw.GetOrgID(ctx); oid != "" {
			ctx = mta.ExtractOutgoing(ctx).Add(OrgIDKey, oid).ToOutgoing(ctx)
		}
		if typ := mw.GetOrgType(ctx); typ != "" {
			ctx = mta.ExtractOutgoing(ctx).Add(OrgTypeKey, typ).ToOutgoing(ctx)
		}
		return invoker(ctx, method, req, reply, cc, opts...)
	}
}

const (
	OrgTypeKey   = "org-type"
	ProfileIDKey = "profileid"
)

func LoadOrgProfile(st interface {
	GetProfileID(context.Context, string) (string, error)
},
) meta.MetaFunc {
	return func(ctx context.Context) (context.Context, error) {
		md := mta.ExtractIncoming(ctx)
		if o := md.Get(hydra.OrgIDKey); o != "" {
			p, err := st.GetProfileID(ctx, o)
			if err == nil {
				ctx = md.Set(ProfileIDKey, p).ToIncoming(ctx)
			}
		}
		return ctx, nil
	}
}

func GetProfileID(ctx context.Context) string { return mta.ExtractIncoming(ctx).Get(ProfileIDKey) }
func GetOrgType(ctx context.Context) string   { return mta.ExtractIncoming(ctx).Get(OrgTypeKey) }
