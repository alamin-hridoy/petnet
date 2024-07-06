package user

import (
	"context"
	"log"

	"google.golang.org/grpc"

	"brank.as/petnet/api/core"
	"brank.as/petnet/api/core/static"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"

	ppb "brank.as/petnet/gunk/drp/v1/profile"
)

type (
	WISEVal struct{}

	CEBVal struct{}
)

func (*WISEVal) Kind() string {
	return static.WISECode
}

func (*CEBVal) Kind() string {
	return static.CEBCode
}

type Validator interface {
	RegisterUserValidate(context.Context, *ppb.RegisterUserRequest) error
	CreateProfileValidate(context.Context, *ppb.CreateProfileRequest) (*core.CreateProfileReq, error)
	GetProfileValidate(context.Context, *ppb.GetProfileRequest) (*core.GetProfileReq, error)
	GetUserValidate(context.Context, *ppb.GetUserRequest) (*core.GetUserRequest, error)
	CreateRecipientValidate(ctx context.Context, req *ppb.CreateRecipientRequest) error
	UpdateRecipientValidate(ctx context.Context, req *ppb.UpdateRecipientRequest) error
	GetRecipientsValidate(ctx context.Context, req *ppb.GetRecipientsRequest) error
	DeleteRecipientValidate(ctx context.Context, req *ppb.DeleteRecipientRequest) error
	Kind() string
}

func NewValidators() []Validator {
	return []Validator{&WISEVal{}, &CEBVal{}}
}

type UserStore interface {
	RegisterUser(ctx context.Context, req core.RegisterUserReq) (*core.RegisterUserResp, error)
	CreateProfile(ctx context.Context, req core.CreateProfileReq) (*core.CreateProfileResp, error)
	GetProfile(ctx context.Context, req core.GetProfileReq) (*core.GetProfileResp, error)
	GetUser(ctx context.Context, req core.GetUserRequest) (*core.GetUserResponse, error)
	CreateRecipient(ctx context.Context, req *ppb.CreateRecipientRequest) (*ppb.CreateRecipientResponse, error)
	RefreshRecipient(ctx context.Context, req *ppb.UpdateRecipientRequest) (*ppb.UpdateRecipientResponse, error)
	GetRecipients(ctx context.Context, req *ppb.GetRecipientsRequest) (*ppb.GetRecipientsResponse, error)
	DeleteRecipient(ctx context.Context, req *ppb.DeleteRecipientRequest) (*ppb.DeleteRecipientResponse, error)
}

type Svc struct {
	ppb.UnimplementedProfileServiceServer
	user       UserStore
	validators map[string]Validator
}

func New(st UserStore, vs []Validator) *Svc {
	s := &Svc{
		validators: make(map[string]Validator, len(vs)),
		user:       st,
	}
	for i, v := range vs {
		switch {
		case v == nil:
			log.Fatalf("validator %d nil", i)
		case v.Kind() == "":
			log.Fatalf("validator %d missing partner type", i)
		}
		s.validators[v.Kind()] = v
	}
	return s
}

// Register the remit service.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	ppb.RegisterProfileServiceServer(srv, s)
	return nil
}

// RegisterGateway parter endpoints.
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, endpoint string, opts []grpc.DialOption) (err error) {
	return ppb.RegisterProfileServiceHandlerFromEndpoint(ctx, mux, endpoint, opts)
}
