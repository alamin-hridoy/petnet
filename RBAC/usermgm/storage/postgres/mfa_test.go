package postgres

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"

	"brank.as/rbac/svcutil/random"
	"brank.as/rbac/usermgm/storage"
)

func TestMFA(t *testing.T) {
	t.Parallel()
	ts := newTestStorage(t)
	code, err := random.String(7)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	u, err := ts.CreateUser(ctx, storage.User{
		OrgID:         uuid.New().String(),
		Username:      "mfatestuser",
		FirstName:     "MFA",
		LastName:      "User",
		Email:         "test@example.com",
		EmailVerified: true,
		InviteCode:    randomString(10),
	}, storage.Credential{
		Username: "mfatestuser",
		Password: randomString(20),
	})
	if err != nil {
		t.Fatal(err)
	}

	cfm := cmp.FilterPath(func(p cmp.Path) bool { return p.Last().String() == ".Confirmed" },
		cmp.Comparer(func(a, b sql.NullTime) bool { return a.Valid || b.Valid }),
	)

	t.Run("Validate", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		ma := storage.MFA{
			UserID:   u.ID,
			MFAType:  storage.PINCode,
			Token:    code,
			Active:   true,
			Deadline: time.Now().Add(time.Minute).Truncate(time.Millisecond),
		}
		m, err := ts.CreateMFA(ctx, ma)
		if err != nil {
			t.Error(err)
		}
		o := cmpopts.IgnoreFields(ma, "ID", "Token", "Created", "Updated")
		if !cmp.Equal(ma, *m, o, cfm) {
			t.Error(cmp.Diff(ma, *m, o, cfm))
		}

		ce, err := ts.CreateMFAEvent(ctx, storage.MFAEvent{
			UserID:   u.ID,
			MFAID:    m.ID,
			MFAType:  storage.PINCode,
			Desc:     "validation",
			Deadline: time.Now().Add(10 * time.Second),
		})
		if err != nil {
			t.Fatal(err)
		}

		eev, err := ts.ConfirmMFAEvent(ctx, storage.MFAEvent{
			UserID:  u.ID,
			MFAID:   m.ID,
			MFAType: storage.PINCode,
			Token:   randomString(8),
		})
		if err == nil {
			t.Error("invalid success", eev)
		}

		cf, err := ts.ConfirmMFAEvent(ctx, storage.MFAEvent{
			UserID:  u.ID,
			MFAID:   m.ID,
			MFAType: storage.PINCode,
			Token:   code,
		})
		if err != nil {
			t.Error(err)
		}

		if o := cmpopts.IgnoreFields(*ce, "Active", "Confirmed"); !cmp.Equal(ce, cf, o) {
			t.Error(cmp.Diff(ce, cf, o))
		}
		if cf.Active {
			t.Error("event still active")
		}
		if !cf.Confirmed.Valid {
			t.Error("event confirm timestamp missing")
		}
	})

	t.Run("Activate", func(t *testing.T) {
		t.Parallel()
		ctx := context.Background()
		m, err := ts.CreateMFA(ctx, storage.MFA{
			UserID:   u.ID,
			MFAType:  storage.SMS,
			Token:    randomString(10),
			Deadline: time.Now().Add(time.Minute),
		})
		if err != nil {
			t.Error(err)
		}
		tok := randomString(8)
		ev, err := ts.CreateMFAEvent(ctx, storage.MFAEvent{
			UserID:   u.ID,
			MFAID:    m.ID,
			MFAType:  storage.SMS,
			Token:    tok,
			Desc:     "activate",
			Deadline: time.Now().Add(20 * time.Second),
		})
		if err != nil {
			t.Error(err)
		}

		conf, err := ts.ConfirmMFAEvent(ctx, *ev)
		if err != nil {
			t.Error(err)
		}
		o := cmpopts.IgnoreFields(storage.MFAEvent{}, "Active", "Token")
		if !cmp.Equal(ev, conf, o, cfm) {
			t.Error(cmp.Diff(ev, conf, o, cfm))
		}

		dup, err := ts.ConfirmMFAEvent(ctx, *ev)
		if err == nil {
			t.Error("duplicate event", dup)
		}

		aev, err := ts.EnableMFA(ctx, u.ID, conf.MFAID)
		if err != nil {
			t.Error(err)
		}

		got, err := ts.GetMFAByID(ctx, m.ID)
		if err != nil {
			t.Error(err)
		}

		if !cmp.Equal(*aev, got.Confirmed.Time) {
			t.Error(cmp.Diff(*aev, got.Confirmed.Time))
		}
	})
}
