package rbsignup

import (
	"context"
	"unicode"

	"brank.as/petnet/serviceutil/logging"
	apb "brank.as/rbac/gunk/v1/authenticate"
	upb "brank.as/rbac/gunk/v1/user"
	"brank.as/rbac/usermgm/errors/session"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Svc struct {
	upb.UnsafeSignupServer
	scl upb.SignupClient
}

// New User Service.
func New(cl upb.SignupClient) *Svc {
	return &Svc{scl: cl}
}

// RegisterService with grpc server.
func (s *Svc) RegisterSvc(srv *grpc.Server) error {
	upb.RegisterSignupServer(srv, s)
	return nil
}

// RegisterService with grpc server.
func (s *Svc) Register(srv *grpc.Server) { upb.RegisterSignupServer(srv, s) }

// RegisterGateway grpcgw
func (s *Svc) RegisterGateway(ctx context.Context, mux *runtime.ServeMux, address string, options []grpc.DialOption) error {
	return upb.RegisterSignupHandlerFromEndpoint(ctx, mux, address, options)
}

func (s *Svc) ResetPassword(ctx context.Context, req *upb.ResetPasswordRequest) (*upb.ResetPasswordResponse, error) {
	if !isPasswordValid(req.GetPassword()) {
		err := status.Error(codes.InvalidArgument, "The password entered does not meet policy.")
		return s.ErrorHandler(ctx, err)
	}

	ctx = metautils.ExtractIncoming(ctx).ToOutgoing(ctx)
	return s.scl.ResetPassword(ctx, req)
}

func (s *Svc) Signup(ctx context.Context, req *upb.SignupRequest) (*upb.SignupResponse, error) {
	ctx = metautils.ExtractIncoming(ctx).ToOutgoing(ctx)
	return s.scl.Signup(ctx, req)
}

func (s *Svc) ResendConfirmEmail(ctx context.Context, req *upb.ResendConfirmEmailRequest) (*upb.ResendConfirmEmailResponse, error) {
	ctx = metautils.ExtractIncoming(ctx).ToOutgoing(ctx)
	return s.scl.ResendConfirmEmail(ctx, req)
}

func (s *Svc) EmailConfirmation(ctx context.Context, req *upb.EmailConfirmationRequest) (*upb.EmailConfirmationResponse, error) {
	ctx = metautils.ExtractIncoming(ctx).ToOutgoing(ctx)
	return s.scl.EmailConfirmation(ctx, req)
}

func (s *Svc) ForgotPassword(ctx context.Context, req *upb.ForgotPasswordRequest) (*upb.ForgotPasswordResponse, error) {
	ctx = metautils.ExtractIncoming(ctx).ToOutgoing(ctx)
	return s.scl.ForgotPassword(ctx, req)
}

func isPasswordValid(s string) bool {
	var (
		hasMinLen  = false
		hasUpper   = false
		hasLower   = false
		hasNumber  = false
		hasSpecial = false
	)
	if len(s) >= 8 && len(s) <= 64 {
		hasMinLen = true
	}
	for _, char := range s {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}
	return hasMinLen && hasUpper && hasLower && hasNumber && hasSpecial
}

func (s *Svc) ErrorHandler(ctx context.Context, err error) (*upb.ResetPasswordResponse, error) {
	log := logging.FromContext(ctx)
	log.Info("received")
	switch status.Code(err) {
	case codes.NotFound:
		return nil, session.Error(codes.NotFound, "NotFound", &apb.SessionError{
			Message:           "User Account Not Found",
			RemainingAttempts: 1,
			ErrorDetails: map[string]string{
				"NotFoundError": err.Error(),
			},
		})
	case codes.PermissionDenied:
		return nil, session.Error(codes.PermissionDenied, "PermissionDenied", &apb.SessionError{
			Message:           "Permission Denied",
			RemainingAttempts: 1,
			ErrorDetails: map[string]string{
				"PermissionError": err.Error(),
			},
		})
	case codes.ResourceExhausted:
		return nil, session.Error(codes.ResourceExhausted, "ResourceExhausted", &apb.SessionError{
			Message:           "Resource Exhausted",
			RemainingAttempts: 1,
			ErrorDetails: map[string]string{
				"ResourceError": err.Error(),
			},
		})
	case codes.InvalidArgument:
		return nil, session.Error(codes.InvalidArgument, "PasswordRulesInvalid", &apb.SessionError{
			Message:           "The password entered does not meet policy.",
			RemainingAttempts: 1,
			ErrorDetails: map[string]string{
				"ResourceError": err.Error(),
			},
		})
	default:
		return nil, session.Error(codes.PermissionDenied, "default", &apb.SessionError{
			Message:           "The password entered does not meet policy.",
			RemainingAttempts: 1,
			ErrorDetails: map[string]string{
				"DefaultError": err.Error(),
			},
		})
	}
}
