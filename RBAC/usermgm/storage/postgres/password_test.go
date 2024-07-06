package postgres

import (
	"context"
	"testing"
	"time"

	"github.com/google/uuid"

	"brank.as/rbac/usermgm/storage"
)

func TestPasswordReset(t *testing.T) {
	t.Parallel()
	ts := newTestStorage(t)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		in := storage.User{
			OrgID:      uuid.New().String(),
			Username:   "TestPasswordReset-name-reset",
			Email:      "TestPasswordReset-email-valid",
			InviteCode: "TestPasswordReset",
		}
		cred := storage.Credential{}
		user, err := ts.CreateUser(context.TODO(), in, cred)
		if err != nil {
			t.Fatalf("CreateSignup() = got error %v, want nil", err)
		}

		code, err := ts.CreatePasswordReset(context.TODO(), user.ID, time.Hour)
		if len(code) != 36 {
			t.Fatalf("CreatePasswordReset() = bad confirmation code UUIDv4: %s", code)
		}
		if err != nil {
			t.Fatalf("CreatePasswordReset() = got error %v, want nil", err)
		}

		newPassword := "new-password"
		if err := ts.PasswordReset(context.TODO(), code, newPassword); err != nil {
			t.Fatalf("PasswordReset() = got error %v, want nil", err)
		}
		if err := ts.PasswordReset(context.TODO(), code, newPassword); err != storage.NotFound {
			t.Fatalf("PasswordReset() = update with expired code expecting storage.ErrnotFound got %v", err)
		}

		if _, err := ts.GetUser(context.TODO(), in.Username, newPassword); err != nil {
			t.Fatalf("GetUserByID() = new password hash doesn't match: %+v \n pw-db=%s pw-new=%s",
				err, cred.Password, newPassword)
		}
	})

	t.Run("Invalid code", func(t *testing.T) {
		t.Parallel()
		if err := ts.PasswordReset(context.TODO(), "12345678-aaaa-bbbb-cccc-1234567890ab", "password"); err != storage.NotFound {
			t.Fatalf("PasswordReset() = update with expired code expecting storage.ErrnotFound got %v", err)
		}
	})

	t.Run("Expired", func(t *testing.T) {
		t.Parallel()
		ctx := context.TODO()
		in := storage.User{
			OrgID:      uuid.New().String(),
			Username:   "TestPasswordReset-name",
			Email:      "TestPasswordReset-email-expired",
			InviteCode: randomString(10),
		}
		cred := storage.Credential{
			Password: "TestPasswordReset-password",
		}
		user, err := ts.CreateUser(ctx, in, cred)
		if err != nil {
			t.Fatalf("CreateSignup() = got error %v, want nil", err)
		}

		code, err := ts.CreatePasswordReset(ctx, user.ID, -time.Second)
		if len(code) != 36 {
			t.Fatalf("CreatePasswordReset() = bad confirmation code UUIDv4: %s", code)
		}
		if err != nil {
			t.Fatalf("CreatePasswordReset() = got error %v, want nil", err)
		}

		newPassword := "new-password"
		if err := ts.PasswordReset(ctx, code, newPassword); err != storage.NotFound {
			t.Fatalf("PasswordReset() = update with expired code expecting storage.ErrnotFound got %v", err)
		}
	})
}

func TestChangePassword(t *testing.T) {
	t.Parallel()
	ts := newTestStorage(t)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		evID := uuid.New().String()
		oldPass := "TestChangePassword-old-password"
		newPass := "TestChangePassword-new-password"

		in := storage.User{
			OrgID:    uuid.New().String(),
			Username: "TestChangePassword-name-reset",
			Email:    "TestChangePassword-email-valid",
		}
		cred := storage.Credential{
			Password: oldPass,
		}
		user, err := ts.CreateUser(context.TODO(), in, cred)
		if err != nil {
			t.Fatalf("CreateUser() = got error %v, want nil", err)
		}
		if _, err := ts.GetUser(context.TODO(), in.Username, oldPass); err != nil {
			t.Fatalf("GetUser() = password hash doesn't match: %+v \n pw-db=%s pw-new=%s",
				err, cred.Password, oldPass)
		}

		if _, err := ts.CreateChangePassword(context.TODO(), user.ID, evID, newPass); err != nil {
			t.Fatalf("CreateChangePassword() = got error %v, want nil", err)
		}

		if err := ts.ChangePassword(context.TODO(), user.ID, evID); err != nil {
			t.Fatalf("ChangePassword() = got error %v, want nil", err)
		}
		if _, err := ts.GetUser(context.TODO(), in.Username, newPass); err != nil {
			t.Fatalf("GetUser() = new password hash doesn't match: %+v \n pw-db=%s pw-new=%s",
				err, cred.Password, newPass)
		}
	})
}
