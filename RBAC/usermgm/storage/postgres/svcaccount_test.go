package postgres

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/uuid"

	"brank.as/rbac/usermgm/storage"
)

func TestSvcAcct(t *testing.T) {
	t.Parallel()
	ts := newTestStorage(t)
	tests := []struct {
		name string
		acct storage.SvcAccount
		err  error
	}{
		{
			name: "valid",
			acct: storage.SvcAccount{
				OrgID:        uuid.New().String(),
				Environment:  "sandbox",
				ClientName:   "client1",
				ClientID:     "id-1",
				CreateUserID: "user1",
			},
		},
		{
			name: "missing org",
			acct: storage.SvcAccount{
				ClientName:   "client2",
				ClientID:     "id-2",
				CreateUserID: "user2",
			},
			err: fmt.Errorf("invalid account client_id, client_name, org_id, and create_user_id are required"),
		},
		{
			name: "valid key",
			acct: storage.SvcAccount{
				OrgID:        uuid.New().String(),
				Environment:  "sandbox",
				ClientName:   "client3",
				ClientID:     "key-1",
				CreateUserID: "user1",
				Challenge:    randomString(32),
			},
		},
	}

	opt := []cmp.Option{
		cmp.FilterPath(func(p cmp.Path) bool {
			return p.Last().String() == ".Created"
		}, cmp.Comparer(func(a, b time.Time) bool {
			return a.Sub(b) > -10*time.Second || a.Sub(b) < 10*time.Second
		})),
		cmp.FilterPath(func(p cmp.Path) bool {
			return p.Last().String() == ".Challenge"
		}, cmp.Ignore()),
	}
	for _, tst := range tests {
		tst := tst
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			clID, err := ts.CreateSvcAccount(context.TODO(), tst.acct)
			if err != nil {
				switch {
				case tst.err == nil:
					t.Fatalf("unexpected error %v", err)
				case !cmp.Equal(tst.err.Error(), err.Error(), opt...):
					t.Error(cmp.Diff(tst.err.Error(), err.Error(), opt...))
				default:
					return
				}
			}
			if !cmp.Equal(tst.acct.ClientID, clID, opt...) {
				t.Error(cmp.Diff(tst.acct.ClientID, clID, opt...))
			}
			cl, err := ts.GetSvcAccountByID(context.TODO(), clID)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(&tst.acct, cl, opt...) {
				t.Error(cmp.Diff(&tst.acct, cl, opt...))
			}
			orgCl, err := ts.GetSvcAccountByOrgID(context.TODO(), cl.OrgID)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal([]storage.SvcAccount{*cl}, orgCl, opt...) {
				t.Error(cmp.Diff([]storage.SvcAccount{*cl}, orgCl, opt...))
			}
		})
	}
}

func TestValidateAPIKey(t *testing.T) {
	t.Parallel()
	keys := []string{
		randomString(32),
		randomString(32),
		randomString(32),
		randomString(32),
	}
	ts := newTestStorage(t)
	tests := []struct {
		name string
		acct storage.SvcAccount
		chlg string
		err  error
	}{
		{
			name: "valid",
			acct: storage.SvcAccount{
				OrgID:        uuid.New().String(),
				Environment:  "sandbox",
				ClientName:   "client1",
				ClientID:     "apikey-1",
				CreateUserID: "user1",
				Challenge:    keys[0],
			},
			chlg: keys[0],
		},
		{
			name: "missing challenge",
			acct: storage.SvcAccount{
				OrgID:        uuid.New().String(),
				ClientName:   "client2",
				ClientID:     "apikey-2",
				CreateUserID: "user2",
				Challenge:    keys[1],
			},
			chlg: "",
			err:  fmt.Errorf("missing service account key"),
		},
		{
			name: "incorrect key",
			acct: storage.SvcAccount{
				OrgID:        uuid.New().String(),
				Environment:  "sandbox",
				ClientName:   "client3",
				ClientID:     "apikey-3",
				CreateUserID: "user1",
				Challenge:    keys[2],
			},
			chlg: keys[3],
			err:  storage.NotFound,
		},
	}
	opt := []cmp.Option{
		cmp.FilterPath(func(p cmp.Path) bool {
			return p.Last().String() == ".Created"
		}, cmp.Comparer(func(a, b time.Time) bool {
			return a.Sub(b) > -10*time.Second || a.Sub(b) < 10*time.Second
		})),
		cmp.FilterPath(func(p cmp.Path) bool {
			return p.Last().String() == ".Challenge"
		}, cmp.Ignore()),
	}
	for _, tst := range tests {
		tst := tst
		t.Run(tst.name, func(t *testing.T) {
			t.Parallel()
			clID, err := ts.CreateSvcAccount(context.TODO(), tst.acct)
			if err != nil {
				t.Fatalf("unexpected error %v", err)
			}
			if !cmp.Equal(tst.acct.ClientID, clID, opt...) {
				t.Error(cmp.Diff(tst.acct.ClientID, clID, opt...))
			}
			cl, err := ts.GetSvcAccountByID(context.TODO(), clID)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(&tst.acct, cl, opt...) {
				t.Error(cmp.Diff(&tst.acct, cl, opt...))
			}
			org, err := ts.ValidateSvcAccount(context.TODO(), tst.acct.ClientID, tst.chlg)
			if err != nil {
				switch {
				case tst.err == nil:
					t.Fatalf("unexpected error %v", err)
				case !cmp.Equal(tst.err.Error(), err.Error(), opt...):
					t.Error(cmp.Diff(tst.err.Error(), err.Error(), opt...))
				default:
					return
				}
			}
			if !cmp.Equal(&tst.acct, org, opt...) {
				t.Error(cmp.Diff(&tst.acct, org, opt...))
			}
		})
	}
}

// randomString using PRNG for repeatable tests.
func randomString(length int) string {
	chars := []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZÅÄÖ" +
		"abcdefghijklmnopqrstuvwxyzåäö" +
		"0123456789")
	var b strings.Builder
	for i := 0; i < length; i++ {
		b.WriteRune(chars[rand.Intn(len(chars))])
	}
	return b.String()
}
