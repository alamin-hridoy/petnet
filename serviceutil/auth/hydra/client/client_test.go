package client

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

func TestAdminOptions(t *testing.T) {
	host := os.Getenv("HYDRA_HOST")
	if host == "" {
		t.Skipf("missing %q env", "HYDRA_HOST")
	}

	tests := []struct {
		desc      string
		host      string
		runCreate bool
		wantErr   error
	}{
		{
			desc: "ValidWithHost",
			host: host,
		},
		{
			desc:    "InvalidHost",
			host:    "127.0 .0.1",
			wantErr: Error{Message: `hydra admin client initialization failed`},
		},
	}

	for _, test := range tests {
		client, err := NewAdminClient(test.host)
		if err != nil {
			if test.runCreate {
				_, err := client.GetClient(context.TODO(), "test")
				if err.Error() != test.wantErr.Error() {
					t.Errorf("%s: error mismatch got: %v, want: %v", test.desc, err, test.wantErr)
				}
				// skip assertion below
				continue
			}
			if err.Error() != test.wantErr.Error() {
				t.Errorf("%s: error mismatch got: %v, want: %v", test.desc, err, test.wantErr)
			}
		}
	}
}

func TestAdminWithoutNameSpace(t *testing.T) {
	host := os.Getenv("HYDRA_HOST")
	if host == "" {
		t.Skipf("missing %q env", "HYDRA_HOST")
	}
	svc, err := NewAdminClient(host)
	if err != nil {
		t.Fatal(err)
	}

	executeClientCreation(t, svc, "", "")
}

func executeClientCreation(t *testing.T, svc *AdminClient, clientID, secret string) {
	cl := AuthClient{
		ClientID:                          clientID,
		ClientName:                        "TestClientName",
		RedirectURIs:                      []string{"http://localhost:4000/callback"},
		PostLogoutRedirectURIs:            []string{"http://localhost:4000/logout"},
		FrontChannelLogoutURI:             "http://localhost:4000/frontchannel_logout",
		FrontChannelLogoutSessionRequired: true,
		BackChannelLogoutURI:              "http://localhost:4000/backchannel_logout",
		BackChannelLogoutSessionRequired:  true,
		GrantTypes: []string{
			"client_credentials",
			"authorization_code",
			"implicit",
			"refresh_token",
		},
		ResponseTypes: []string{
			"id_token",
			"code",
			"token",
		},
		Scopes: []string{
			"as.brank.ecebuana/read",
			"as.brank.ecebuana/write",
			"offline_access",
			"offline",
			"openid",
		},
		Audience: []string{
			"api.testclient.staging.bnk.to",
			"auth.testclient.staging.bnk.to",
		},
		Secret:      secret,
		SubjectType: "public",
		AuthMethod:  "client_secret_post",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// ensure that back-to-back runs will not fail.
	if err := svc.DeleteClient(ctx, clientID); err != nil {
		t.Logf("cleanup failed (ignore if this is the first test run): %v", err)
	}

	// test the full create/delete cycle, then create again to inspect the results
	// abort if any stage fails, to allow inspection of the hydra state.
	res, err := svc.CreateClient(ctx, cl)
	if err != nil {
		t.Fatal(err)
	}

	if res.ClientID == "" || res.Secret == "" {
		t.Fatal("failed to generate client id and secret")
	}

	clientID = res.ClientID
	got, err := svc.GetClient(ctx, clientID)
	if err != nil {
		t.Fatal(err)
	}

	want := cl
	want.ClientID = res.ClientID
	want.Secret = ""
	if !cmp.Equal(want, *got) {
		t.Error(cmp.Diff(want, *got))
	}
	got.Secret = "NewSecretValue"
	got.Scopes = []string{
		"as.brank.ecebuana/read",
		"as.brank.ecebuana/write",
	}
	want = *got
	err = svc.UpdateClient(ctx, *got)
	if err != nil {
		t.Error(err)
	}

	if diff := cmp.Diff(want, got); diff == "" {
		t.Error("failed to update client")
	}

	if err := svc.DeleteClient(ctx, clientID); err != nil {
		t.Fatal(err)
	}
	if _, err := svc.CreateClient(ctx, cl); err != nil {
		t.Fatal(err)
	}
}
