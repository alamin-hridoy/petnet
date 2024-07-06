package postgres

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"brank.as/rbac/usermgm/storage"
)

func TestCreateSignup(t *testing.T) {
	t.Parallel()
	ts := newTestStorage(t)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		in := &storage.User{
			OrgID:      uuid.New().String(),
			Username:   "TestCreateSignup-name-valid",
			Email:      "TestCreateSignup-email",
			InviteCode: "TestCreateSignup",
		}
		cred := storage.Credential{
			Password: "TestCreateSignup-password",
		}
		user, err := ts.CreateUser(context.TODO(), *in, cred)
		if err != nil {
			t.Fatalf("CreateSignup() = got error %v, want nil", err)
		}
		if user.ID == "" {
			t.Fatal("CreateSignup() = returned empty ID")
		}
	})

	t.Run("ExistingEmail", func(t *testing.T) {
		t.Parallel()
		// Create the entry first.
		in := &storage.User{
			OrgID:      uuid.New().String(),
			Username:   "TestCreateSignup-name",
			Email:      "TestCreateSignup-existing-email",
			InviteCode: "TestCreateSignup2",
		}
		cred := storage.Credential{
			Password: "TestCreateSignup-password",
		}

		if _, err := ts.CreateUser(context.TODO(), *in, cred); err != nil {
			t.Fatalf("CreateSignup() = got error %v, want nil", err)
		}

		// Attempt to create the same entry.
		wantErr := storage.UsernameExists
		if _, err := ts.CreateUser(context.TODO(), *in, cred); err != wantErr {
			t.Fatalf("wrong error for duplicated email: expected=%+v, actual=%+v", wantErr, err)
		}
	})

	t.Run("ExistingInvCode", func(t *testing.T) {
		t.Parallel()
		// Create the entry first.
		in := &storage.User{
			OrgID:      uuid.New().String(),
			Username:   "TestCreateSignup-name-existinv",
			Email:      "TestCreateSignup-existing-email2",
			InviteCode: "TestCreateSignup3",
		}
		cred := storage.Credential{
			Password: "TestCreateSignup-password",
		}

		if _, err := ts.CreateUser(context.TODO(), *in, cred); err != nil {
			t.Fatalf("CreateSignup() = got error %v, want nil", err)
		}

		in.Email += "new"
		// Attempt to create the same entry.
		in.Username += "again"
		wantErr := storage.InvCodeExists
		if _, err := ts.CreateUser(context.TODO(), *in, cred); err != wantErr {
			t.Fatalf("wrong error for duplicated inv code: expected=%+v, actual=%+v", wantErr, err)
		}
	})

	t.Run("SetUsername", func(t *testing.T) {
		t.Parallel()
		// Create the entry first.
		in := &storage.User{
			OrgID:      uuid.New().String(),
			Username:   "TestCreateSignup-existing-email3",
			Email:      "TestCreateSignup-existing-email3",
			InviteCode: "TestCreateSignup4",
		}
		cred := storage.Credential{
			Password: "TestCreateSignup-password",
		}
		usr, err := ts.CreateUser(context.TODO(), *in, cred)
		if err != nil {
			t.Fatalf("CreateSignup() = got error %v, want nil", err)
		}

		newUsername := "TestCreateSignup-name-setuser"
		if err := ts.SetUsername(context.TODO(), in.InviteCode, newUsername); err != nil {
			t.Error(err)
		}
		u, err := ts.GetUser(context.TODO(), newUsername, cred.Password)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(usr.ID, u.ID) {
			t.Error(cmp.Diff(usr.ID, u.ID))
		}
	})
}

func TestGetUser(t *testing.T) {
	t.Parallel()
	ts := newTestStorage(t)

	t.Run("Valid-ByID", func(t *testing.T) {
		t.Parallel()
		in := &storage.User{
			OrgID:      uuid.New().String(),
			Username:   "TestGetUser-name",
			Email:      "TestGetUser-getUser-email-valid-ByID",
			InviteCode: "TestGetUser",
		}
		cred := storage.Credential{
			Password: "TestGetUser-password",
		}
		// Store password as CreateSignup will overwrite it with the hash
		pw := cred.Password
		want, err := ts.CreateUser(context.TODO(), *in, cred)
		if err != nil {
			t.Fatalf("CreateSignup() = got error %+v, want nil", err)
		}
		user, err := ts.GetUserByID(context.TODO(), want.ID)
		if err != nil {
			t.Fatalf("GetUserByID() = got error %v, want nil", err)
		}
		if user.Email != in.Email {
			t.Fatalf("GetUserByID() = got bad email %v, want %v", user.Email, in.Email)
		}
		if user.Username != in.Username {
			t.Fatalf("GetUserByID() = got bad username %v, want %v", user.Username, in.Username)
		}
		usr, err := ts.GetUser(context.TODO(), in.Username, cred.Password)
		if err != nil {
			t.Fatalf("GetUserByID() = password hash doesn't match: %+v \n pw-db=%s pw-in=%s",
				err, cred.Password, pw)
		}
		if !cmp.Equal(want.ID, usr.ID) {
			t.Error(cmp.Diff(want.ID, usr.ID))
		}

		if !cmp.Equal(in.OrgID, usr.OrgID) {
			t.Error(cmp.Diff(in.OrgID, usr.OrgID))
		}
	})

	t.Run("Invalid-ByID", func(t *testing.T) {
		t.Parallel()
		const badUUID string = "12345678-aaaa-bbbb-cccc-1234567890ab"
		_, err := ts.GetUserByID(context.TODO(), badUUID)
		if err != storage.NotFound {
			t.Fatalf("GetUserByID() = got error %v, want storage.ErrNotFound", err)
		}
	})

	t.Run("Valid-ByEmail", func(t *testing.T) {
		t.Parallel()
		in := &storage.User{
			OrgID:      uuid.New().String(),
			Username:   "TestGetUser-namevalid",
			Email:      "TestGetUser-getUser-email-valid-ByEmail",
			InviteCode: "TestGetUser2",
		}
		cred := storage.Credential{
			Password: "TestGetUser-password",
		}
		// Store password as CreateSignup will overwrite it with the hash
		pw := cred.Password
		u, err := ts.CreateUser(context.TODO(), *in, cred)
		if err != nil {
			t.Fatalf("CreateSignup() = got error %+v, want nil", err)
		}
		in.ID = u.ID
		user, err := ts.GetUserByEmail(context.TODO(), in.Email)
		if err != nil {
			t.Fatalf("GetUserByID() = got error %v, want nil", err)
		}
		if user.Email != in.Email {
			t.Fatalf("GetUserByID() = got bad email %v, want %v", user.Email, in.Email)
		}
		if user.Username != in.Username {
			t.Fatalf("GetUserByID() = got bad username %v, want %v", user.Username, in.Username)
		}
		usr, err := ts.GetUser(context.TODO(), in.Username, cred.Password)
		if err != nil {
			t.Fatalf("GetUserByID() = password hash doesn't match: %+v \n pw-db=%s pw-in=%s",
				err, cred.Password, pw)
		}
		if !cmp.Equal(in.ID, usr.ID) {
			t.Error(cmp.Diff(in.ID, usr.ID))
		}
		if !cmp.Equal(in.OrgID, usr.OrgID) {
			t.Error(cmp.Diff(in.OrgID, usr.OrgID))
		}
	})

	t.Run("Invalid-ByEmail", func(t *testing.T) {
		t.Parallel()
		_, err := ts.GetUserByEmail(context.TODO(), "bad-email@brank.as")
		if err != storage.NotFound {
			t.Fatalf("GetUserByID() = got error %v, want storage.ErrNotFound", err)
		}
	})
}

func TestEmailVerification(t *testing.T) {
	t.Parallel()
	ts := newTestStorage(t)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		in := &storage.User{
			OrgID:      uuid.New().String(),
			Username:   "TestEmailVerification-name",
			Email:      "TestEmailVerification-verify-email-valid",
			InviteCode: "TestEmailVerification",
		}
		cred := storage.Credential{
			Password: "TestEmailVerification-password",
		}
		user, err := ts.CreateUser(context.TODO(), *in, cred)
		if err != nil {
			t.Fatalf("CreateSignup() = got error %v, want nil", err)
		}
		code, err := ts.CreateConfirmationCode(context.TODO(), user.ID)
		if len(code) != 36 {
			t.Fatalf("CreateConfirmationCode() = bad confirmation code UUIDv4: %s", code)
		}
		if err != nil {
			t.Fatalf("CreateConfirmationCode() = got error %v, want nil", err)
		}
		u, err := ts.VerifyConfirmationCode(context.TODO(), code)
		if err != nil {
			t.Fatalf("VerifyConfirmationCode() = got error %v, want nil", err)
		}
		if !cmp.Equal(in.Email, u.Email) {
			t.Error(cmp.Diff(in.Email, u.Email))
		}
		if !cmp.Equal(user.ID, u.ID) {
			t.Error(cmp.Diff(user.ID, u.ID))
		}
		if !cmp.Equal(in.OrgID, u.OrgID) {
			t.Error(cmp.Diff(in.OrgID, u.OrgID))
		}
	})

	t.Run("Bad user ID", func(t *testing.T) {
		t.Parallel()
		const badUUID string = "12345678-aaaa-bbbb-cccc-1234567890ab"
		in := &storage.User{
			OrgID:      uuid.New().String(),
			Username:   "TestEmailVerification-name-baduid",
			Email:      "TestEmailVerification-verify-email-bad-id",
			InviteCode: "badUserID",
		}
		cred := storage.Credential{
			Password: "TestEmailVerification-password",
		}
		_, err := ts.CreateUser(context.TODO(), *in, cred)
		if err != nil {
			t.Fatalf("CreateSignup() = got error %v, want nil", err)
		}
		_, err = ts.CreateConfirmationCode(context.TODO(), badUUID)
		if err == nil {
			t.Fatal("CreateConfirmationCode() = got error nil, want error")
		}
	})

	t.Run("Bad verification code", func(t *testing.T) {
		t.Parallel()
		const badUUID string = "12345678-aaaa-bbbb-cccc-1234567890ab"
		in := &storage.User{
			OrgID:      uuid.New().String(),
			Username:   "TestEmailVerification-name-badver",
			Email:      "TestEmailVerification-verify-email-bad-code",
			InviteCode: "badVerCode",
		}
		cred := storage.Credential{
			Password: "TestEmailVerification-password",
		}
		user, err := ts.CreateUser(context.TODO(), *in, cred)
		if err != nil {
			t.Fatalf("CreateSignup() = got error %v, want nil", err)
		}
		_, err = ts.CreateConfirmationCode(context.TODO(), user.ID)
		if err != nil {
			t.Fatalf("CreateConfirmationCode() = got error %v, want nil", err)
		}
		_, err = ts.VerifyConfirmationCode(context.TODO(), badUUID)
		if err == nil {
			t.Fatal("VerifyConfirmationCode() = got error nil, want error")
		}
	})
}

func TestGetUsers(t *testing.T) {
	t.Parallel()
	ts := newTestStorage(t)

	ctx := context.Background()
	oid := uuid.New().String()
	o := cmpopts.IgnoreFields(storage.User{}, "ID", "InviteExpiry", "Created", "Updated", "Deleted", "InviteSender", "PreferredMFA", "Count")

	fs := []storage.User{
		{
			OrgID:         oid,
			Username:      "Username",
			FirstName:     "A-First-name",
			LastName:      "A-Last-name",
			Email:         "example@gmail.com",
			EmailVerified: false,
			InviteSender:  "Invite-Sender",
			InviteStatus:  "Invited",
			InviteCode:    "53215",
			MFALogin:      false,
		},
		{
			OrgID:         oid,
			Username:      "Username2",
			FirstName:     "B-First-name",
			LastName:      "B-Last-name",
			Email:         "example@gmail2.com",
			EmailVerified: false,
			InviteSender:  "Invite-Sender2",
			InviteStatus:  "Invited",
			InviteCode:    "53216",
			MFALogin:      false,
		},
		{
			OrgID:         oid,
			Username:      "Username3",
			FirstName:     "C-First-name",
			LastName:      "C-Last-name",
			Email:         "example@gmail3.com",
			EmailVerified: false,
			InviteSender:  "Invite-Sender2",
			InviteStatus:  "Invited",
			InviteCode:    "53316",
			MFALogin:      false,
		},
	}
	cred := storage.Credential{
		Password: "test-user-list-password",
	}
	for i, f := range fs {
		u, err := ts.CreateUser(ctx, f, cred)
		if err != nil {
			t.Error(err)
		}
		fs[i].ID = u.ID
	}

	tests := []struct {
		name string
		f    storage.FilterList
		want []storage.User
	}{
		{
			name: "No Limit",
			f: storage.FilterList{
				OrgID:  oid,
				Offset: 0,
			},
			want: []storage.User{fs[0], fs[1], fs[2]},
		},
		{
			name: "Descending List",
			f: storage.FilterList{
				OrgID:        oid,
				SortBy:       "DESC",
				SortByColumn: "first_name",
				Offset:       0,
			},
			want: []storage.User{fs[2], fs[1], fs[0]},
		},
		{
			name: "Ascending List",
			f: storage.FilterList{
				OrgID:        oid,
				SortBy:       "ASC",
				SortByColumn: "first_name",
				Offset:       0,
			},
			want: []storage.User{fs[0], fs[1], fs[2]},
		},
		{
			name: "ID List",
			f: storage.FilterList{
				OrgID:        oid,
				ID:           []string{fs[1].ID, fs[2].ID},
				SortBy:       "DESC",
				SortByColumn: "first_name",
				Limit:        2,
				Offset:       0,
			},
			want: []storage.User{fs[2], fs[1]},
		},
		{
			name: "Limit",
			f: storage.FilterList{
				OrgID:        oid,
				SortBy:       "DESC",
				SortByColumn: "first_name",
				Limit:        2,
				Offset:       0,
			},
			want: []storage.User{fs[2], fs[1]},
		},
		{
			name: "Offset",
			f: storage.FilterList{
				OrgID:        oid,
				SortBy:       "DESC",
				SortByColumn: "first_name",
				Limit:        0,
				Offset:       1,
			},
			want: []storage.User{fs[1], fs[0]},
		},
		{
			name: "Limit+Offset",
			f: storage.FilterList{
				OrgID:        oid,
				SortBy:       "ASC",
				SortByColumn: "first_name",
				Limit:        1,
				Offset:       1,
			},
			want: []storage.User{fs[1]},
		},
		{
			name: "Search+By+Name",
			f: storage.FilterList{
				OrgID:  oid,
				Offset: 0,
				Name:   "A-First-name",
			},
			want: []storage.User{fs[0]},
		},
		{
			name: "Search+By+PartialName",
			f: storage.FilterList{
				OrgID:  oid,
				Offset: 0,
				Name:   "First-name",
			},
			want: []storage.User{fs[0], fs[1], fs[2]},
		},
		{
			name: "Search+By+PartialName2",
			f: storage.FilterList{
				OrgID:  oid,
				Offset: 0,
				Name:   "A-",
			},
			want: []storage.User{fs[0]},
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			got, err := ts.GetUsers(ctx, test.f)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
			}
		})
	}
}
