package user

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"testing"

	"brank.as/petnet/api/core/static"
	"brank.as/petnet/serviceutil/auth/hydra"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	uc "brank.as/petnet/api/core/user"
	"brank.as/petnet/api/integration/perahub"
	"brank.as/petnet/api/storage/postgres"
	upb "brank.as/petnet/gunk/drp/v1/profile"
)

var _testStorage *postgres.Storage

func TestMain(m *testing.M) {
	const dbConnEnv = "DATABASE_CONNECTION"
	ddlConnStr := os.Getenv(dbConnEnv)
	if ddlConnStr == "" {
		log.Printf("%s is not set, skipping", dbConnEnv)
		return
	}

	var teardown func()
	_testStorage, teardown = postgres.NewTestStorage(ddlConnStr, filepath.Join("..", "..", "migrations", "sql"))

	exitCode := m.Run()

	if teardown != nil {
		teardown()
	}

	os.Exit(exitCode)
}

func newTestStorage(tb testing.TB) *postgres.Storage {
	if testing.Short() {
		tb.Skip("skipping tests that use postgres on -short")
	}

	return _testStorage
}

func newTestSvc(t *testing.T, st *postgres.Storage) (*Svc, *perahub.HTTPMock) {
	cl := perahub.NewTestHTTPMock(st, perahub.MockConfig{})
	ph, err := perahub.New(cl,
		"dev",
		"https://newkycgateway.dev.perahub.com.ph/gateway/",
		"https://privatedrp.dev.perahub.com.ph/v1/remit/nonex/",
		"https://privatedrp.dev.perahub.com.ph/v1/billspay/wrapper/api/",
		"https://privatedrp.dev.perahub.com.ph/v1/transactions/api/",
		"https://privatedrp.dev.perahub.com.ph/v1/billspay/",
		"partner-id",
		"client-key",
		"api-key",
		"",
		"",
		map[string]perahub.OAuthCreds{
			static.WISECode: {
				ClientID:     "id",
				ClientSecret: "secret",
			},
		},
	)
	if err != nil {
		t.Fatal("setting up perahub integration: ", err)
	}
	return New(uc.New(st, ph), NewValidators()), cl
}

type tester interface {
	Test(t *testing.T)
}

func TestEnforce(t *testing.T) {
	vs := NewValidators()
	for _, v := range vs {
		for _, sp := range static.Partners["PH"] {
			if v.Kind() == sp.Code {
				ts, ok := v.(tester)
				if !ok {
					t.Fatal("add register testcase for partner: ", v.Kind())
				}
				ts.Test(t)
			}
		}
	}
}

func (WISEVal) Test(t *testing.T) {
	st := newTestStorage(t)

	tests := []struct {
		desc        string
		wantErr     codes.Code
		conflictErr bool
	}{
		{
			desc: "Success",
		},
		{
			desc:        "Conflict",
			wantErr:     codes.AlreadyExists,
			conflictErr: true,
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid, "partner", static.WISECode))
	ctx := nmd.ToIncoming(context.Background())

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			s, m := newTestSvc(t, st)
			if test.conflictErr {
				m.SetConflictError()
			}
			_, err := s.RegisterUser(ctx, &upb.RegisterUserRequest{
				RemitPartner: static.WISECode,
				Email:        "test@mail.com",
			})
			if gotErr := status.Code(err); gotErr != test.wantErr {
				t.Errorf("error mismatch, want: %v, got: %v", test.wantErr, gotErr)
			}
		})
	}

	s, _ := newTestSvc(t, st)
	t.Run("Create Profile", func(t *testing.T) {
		if _, err := s.CreateProfile(ctx, &upb.CreateProfileRequest{
			RemitPartner: static.WISECode,
			Email:        "test@mail.com",
			Type:         "personal",
			FirstName:    "Brankas",
			LastName:     "Sender",
			BirthDate:    "1990-01-10",
			Phone: &upb.PhoneNumber{
				CountryCode: "62",
				Number:      "123456789",
			},
			Address: &upb.Address{
				Address1:   "East Offices Bldg., 114 Aguirre St.,Legaspi Village,",
				City:       "Makati",
				Country:    "PH",
				PostalCode: "1229",
			},
			Occupation: "Software Engineer",
		}); err != nil {
			t.Error(err)
		}
	})

	t.Run("Get Profile", func(t *testing.T) {
		want := &upb.GetProfileResponse{
			Profile: &upb.Profile{
				ID:         "16325688",
				Type:       "personal",
				FirstName:  "Brankas",
				LastName:   "Sender",
				BirthDate:  "1990-01-10",
				Phone:      "+639999999999",
				Occupation: "Software Engineer",
				Address: &upb.Address{
					Address1:   "East Offices Bldg., 114 Aguirre St.,Legaspi Village,",
					City:       "Makati",
					PostalCode: "1229",
					Country:    "ph",
				},
			},
		}
		got, err := s.GetProfile(ctx, &upb.GetProfileRequest{
			RemitPartner: static.WISECode,
			Email:        "test@mail.com",
		})
		if err != nil {
			t.Error(err)
		}

		o := cmp.Options{
			cmpopts.IgnoreUnexported(
				upb.GetProfileResponse{},
				upb.Profile{},
				upb.Date{},
				upb.Address{},
			),
		}
		if !cmp.Equal(want, got, o) {
			t.Error("(-want +got): ", cmp.Diff(want, got, o))
		}
	})
}

func (CEBVal) Test(t *testing.T) {
	st := newTestStorage(t)

	tests := []struct {
		desc     string
		wantErr  codes.Code
		crUsrErr bool
	}{
		{
			desc: "Success",
		},
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid, "partner", static.CEBCode))
	ctx := nmd.ToIncoming(context.Background())

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			s, _ := newTestSvc(t, st)
			_, err := s.RegisterUser(ctx, &upb.RegisterUserRequest{
				RemitPartner:    static.CEBCode,
				FirstName:       "John",
				LastName:        "Doe",
				BirthDate:       "1981-08-12",
				MobileCountryID: "63",
				ContactNo:       "63",
				PhoneCountryID:  "0",
				PhoneArCode:     "1",
				CountryAddID:    "166",
				ProvinceAdd:     "1953 Ph 3b Block 6 Lot 9 Camarin Caloocan City",
				CurrentAdd:      "2863",
				UserID:          "2863",
				SourceOfID:      "1",
				Tin:             "1",
				PhoneNo:         "1",
				AgentCode:       "01030063",
			})
			if gotErr := status.Code(err); gotErr != test.wantErr {
				t.Errorf("error mismatch, want: %v, got: %v", test.wantErr, gotErr)
			}
		})
	}
	t.Run("Get User", func(t *testing.T) {
		s, _ := newTestSvc(t, st)
		want := &upb.GetUserResponse{
			Code:    0,
			Message: "Successful",
			Result: &upb.GUResult{
				User: &upb.User{
					UserID:         3663,
					UserNumber:     "EWFHL0000070155824",
					FirstName:      "GREMAR",
					MiddleName:     "REAS",
					LastName:       "NAPARATE",
					BirthDate:      "1994-11-04T00:00:00",
					MobileCountry:  0,
					PhoneCountry:   0,
					CountryAddress: 0,
					SourceOfFund:   0,
				},
			},
			RemcoID: 9,
		}
		got, err := s.GetUser(ctx, &upb.GetUserRequest{
			RemitPartner: static.CEBCode,
			FirstName:    "John",
			LastName:     "Doe",
			BirthDate:    "1981-08-12",
			UserNumber:   "EWJCJ0000100155830",
		})
		if err != nil {
			t.Error(err)
		}

		o := cmp.Options{
			cmpopts.IgnoreUnexported(
				upb.GetUserResponse{},
				upb.GUResult{},
				upb.User{},
			),
		}
		if !cmp.Equal(want, got, o) {
			t.Error("(-want +got): ", cmp.Diff(want, got, o))
		}
	})
	s, _ := newTestSvc(t, st)
	t.Run("Create Recipient", func(t *testing.T) {
		if _, err := s.CreateRecipient(ctx, &upb.CreateRecipientRequest{
			RemitPartner:     static.CEBCode,
			FirstName:        "Newwiee",
			MiddleName:       "Tay",
			LastName:         "Tawan",
			SenderUserID:     3663,
			BirthDate:        "1994-01-30",
			MobileCountryID:  "63",
			ContactNumber:    "123",
			PhoneCountryID:   "0",
			PhoneAreaCode:    "63",
			PhoneNumber:      "12345",
			CountryAddressID: "166",
			BirthCountryID:   "166",
			ProvinceAddress:  "prov",
			Address:          "1953 Ph 3b Block 6 Lot 9 Camarin Caloocan City",
			UserID:           2863,
			Occupation:       "Software Engineer",
			PostalCode:       "1850",
			StateIDAddress:   "0",
			Tin:              "tin",
		}); err != nil {
			t.Error(err)
		}
	})
	t.Run("Get Recipients", func(t *testing.T) {
		s, _ := newTestSvc(t, st)
		want := &upb.GetRecipientsResponse{
			Recipients: []*upb.Recipient{
				{
					RecipientID:    "3663",
					FirstName:      "GREMAR",
					MiddleName:     "REAS",
					LastName:       "NAPARATE",
					BirthDate:      "1994-11-04T00:00:00",
					StateIDAddress: "0",
					MobileCountry:  0,
					PhoneCountry:   0,
					CountryAddress: 0,
					BirthCountry:   0,
				},
			},
		}
		got, err := s.GetRecipients(ctx, &upb.GetRecipientsRequest{
			RemitPartner: static.CEBCode,
			SenderUserID: "3663",
		})
		if err != nil {
			t.Error(err)
		}

		o := cmp.Options{
			cmpopts.IgnoreUnexported(
				upb.GetRecipientsResponse{},
				upb.Recipient{},
			),
		}
		if !cmp.Equal(want, got, o) {
			t.Error("(-want +got): ", cmp.Diff(want, got, o))
		}
	})
}
