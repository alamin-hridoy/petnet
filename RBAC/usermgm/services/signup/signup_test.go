package signup

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus/hooks/test"

	"brank.as/rbac/serviceutil/logging"
	pgtest "brank.as/rbac/serviceutil/storage/postgres"
	"brank.as/rbac/usermgm/core/org"
	"brank.as/rbac/usermgm/core/permissions"
	"brank.as/rbac/usermgm/core/user"
	"brank.as/rbac/usermgm/integrations/email"
	"brank.as/rbac/usermgm/integrations/keto"
	"brank.as/rbac/usermgm/storage"
	"brank.as/rbac/usermgm/storage/postgres"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	upb "brank.as/rbac/gunk/v1/user"
)

const (
	confirmationCode = "12345678-aaaa-bbbb-cccc-1234567890abcd"
	password         = "a-valid-password"
)

var testSignup = storage.User{
	Username:   "user-full-name",
	FirstName:  "user-first",
	LastName:   "user-last",
	Email:      "valid-email@example.com",
	InviteCode: "valid-code",
}

func TestSignup(t *testing.T) {
	t.Skip("TODO")
	t.Parallel()

	tests := []struct {
		desc string
		in   *upb.SignupRequest
		want *upb.SignupResponse
		// check only for the error code as the message might change from time to time
		wantErrCode codes.Code
		withCode    string
	}{
		{
			desc: "ValidSignup",
			in: &upb.SignupRequest{
				Username:  testSignup.Username,
				FirstName: testSignup.FirstName,
				LastName:  testSignup.LastName,
				Email:     testSignup.Email,
				Password:  password,
			},
			want: &upb.SignupResponse{
				UserID: uuid.New().String(),
				OrgID:  uuid.New().String(),
			},
		},
		{
			desc: "MissingUsername",
			in: &upb.SignupRequest{
				Email:     testSignup.Email,
				FirstName: testSignup.FirstName,
				LastName:  testSignup.LastName,
				Password:  password,
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			desc: "MissingEmail",
			in: &upb.SignupRequest{
				Username:  testSignup.Username,
				FirstName: testSignup.FirstName,
				LastName:  testSignup.LastName,
				Password:  password,
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			desc: "MissingPassword",
			in: &upb.SignupRequest{
				Username:  testSignup.Username,
				FirstName: testSignup.FirstName,
				LastName:  testSignup.LastName,
				Email:     testSignup.Email,
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			desc: "EmailExists",
			in: &upb.SignupRequest{
				Username:  testSignup.Username,
				FirstName: testSignup.FirstName,
				LastName:  testSignup.LastName,
				Email:     testSignup.Email,
				Password:  password,
			},
			wantErrCode: codes.InvalidArgument,
		},
		{
			desc: "InvCodeExists",
			in: &upb.SignupRequest{
				Username:  testSignup.Username,
				FirstName: testSignup.FirstName,
				LastName:  testSignup.LastName,
				Email:     testSignup.Email,
				Password:  password,
			},
			want: &upb.SignupResponse{},
		},
	}

	db, cleanup := pgtest.MustNewDevelopmentDB(os.Getenv("DATABASE_CONNECTION"),
		filepath.Join("..", "..", "migrations", "sql"))
	t.Cleanup(cleanup)
	st := postgres.NewStorageDB(db)
	mail := FakeMailer{
		WantErr: false,
		Code:    confirmationCode,
	}
	k := keto.New(os.Getenv("KETO_URL"))
	pm := permissions.New(st, k)
	core := user.New(user.Config{
		PublicSignup:  true,
		AutoApprove:   true,
		ResetDuration: 259200 * time.Second,
	}, st, st, mail, org.New(st, pm, st), nil)
	for _, tst := range tests {
		t.Run(tst.desc, func(t *testing.T) {
			// t.Parallel()
			logr, _ := test.NewNullLogger()
			log := logr.WithField("test", tst.desc)
			ctx := logging.WithLogger(context.Background(), log)

			handler := New(core, nil)
			got, err := handler.Signup(ctx, tst.in)
			if err != nil {
				// We're expecting a status.Code if we wanted an error via wantErrCode,
				// the error is unexpected if its 0.
				st, ok := status.FromError(err)
				if !ok || tst.wantErrCode == 0 {
					t.Fatalf("%s: want nil error, got %v", tst.desc, err)
				}

				if wantCode := tst.wantErrCode; st.Code() != wantCode {
					t.Fatalf("%s: want error code %s, got %s", tst.desc, wantCode, st.Code())
				}
			}
			if tst.want == nil && got == nil {
				return
			}
			if tst.want == nil && got != nil {
				t.Fatal("unexpected response", got)
			}
			if tst.want != nil && got == nil {
				t.Fatal("response missing", got)
			}
			w, werr := uuid.Parse(tst.want.OrgID)
			g, gerr := uuid.Parse(got.OrgID)
			if (werr != gerr) || (werr == nil && gerr == nil && w.Version() != g.Version()) {
				t.Error("invalid OrgID", tst.want.OrgID, got.OrgID)
			}
			w, werr = uuid.Parse(tst.want.UserID)
			g, gerr = uuid.Parse(got.UserID)
			if (werr != gerr) || (werr == nil && gerr == nil && w.Version() != g.Version()) {
				t.Error("invalid UserID", tst.want.UserID, got.UserID)
			}
		})
	}
}

// FakeMailer for local dev and testing
type FakeMailer struct {
	WantErr bool
	Code    string
}

func (f FakeMailer) ForgotPassword(email, name, code string) error {
	if f.WantErr {
		return fmt.Errorf("an error")
	}
	if f.Code != code {
		return fmt.Errorf("different token")
	}
	return nil
}

func (f FakeMailer) Approved(email string, fm email.Approved) error {
	if f.WantErr {
		return fmt.Errorf("an error")
	}
	return nil
}

func (f FakeMailer) InviteUser(email, inviteCode string, inv email.Invite) error {
	if f.WantErr {
		return fmt.Errorf("an error")
	}
	if f.Code != inviteCode {
		return fmt.Errorf("different token")
	}
	return nil
}

func (f FakeMailer) EmailMFA(email, code string) error {
	if f.WantErr {
		return fmt.Errorf("an error")
	}
	if f.Code != code {
		return fmt.Errorf("different token")
	}
	return nil
}

func (f FakeMailer) DisableUser(email string, dis email.User) error {
	if f.WantErr {
		return fmt.Errorf("an error")
	}
	return nil
}

func (f FakeMailer) EnableUser(email string, en email.User) error {
	if f.WantErr {
		return fmt.Errorf("an error")
	}
	return nil
}
