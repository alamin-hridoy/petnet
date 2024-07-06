package profile

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/sirupsen/logrus/hooks/test"

	"brank.as/petnet/profile/core/profile"
	"brank.as/petnet/profile/storage/postgres"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	tspb "google.golang.org/protobuf/types/known/timestamppb"
)

// todo make into integration test with postgres
func TestOrgProfile(t *testing.T) {
	t.Parallel()
	ctx := context.Background()
	ts := tspb.New(time.Unix(1515151515, 0))
	ts2 := tspb.New(time.Unix(1414141414, 0))
	tests := []struct {
		desc            string
		profileConflict bool
		req             *ppb.UpsertProfileRequest
		ureq            *ppb.UpsertProfileRequest
		reqq            *ppb.UpdateOrgProfileUserIDRequest
		wantErr         bool
	}{
		{
			desc: "All Success",
			req: &ppb.UpsertProfileRequest{
				Profile: &ppb.OrgProfile{
					UserID:    "10000000-0000-0000-0000-000000000000",
					OrgID:     "20000000-0000-0000-0000-000000000000",
					OrgType:   ppb.OrgType_PetNet,
					Status:    ppb.Status_Accepted,
					RiskScore: ppb.RiskScore_Low,
					BusinessInfo: &ppb.BusinessInfo{
						CompanyName:   "company-name",
						StoreName:     "store-name",
						PhoneNumber:   "123456789",
						FaxNumber:     "123456789",
						Website:       "https://website.com",
						CompanyEmail:  "company@mail.com",
						ContactPerson: "contact-person",
						Position:      "position",
						Address: &ppb.Address{
							Address1:   "address1",
							City:       "city",
							State:      "state",
							PostalCode: "12345",
						},
					},
					AccountInfo: &ppb.AccountInfo{
						Bank:                    "bank",
						BankAccountNumber:       "bank-acct-no",
						BankAccountHolder:       "bank-acct-holder",
						AgreeTermsConditions:    ppb.Boolean_True,
						AgreeOnlineSupplierForm: ppb.Boolean_True,
						Currency:                ppb.Currency_PHP,
					},
					ReminderSent:      ppb.Boolean_True,
					DateApplied:       ts,
					Deleted:           ts,
					DsaCode:           "123",
					TerminalIdOtc:     "456",
					TerminalIdDigital: "789",
					IsProvider:        false,
				},
			},
			ureq: &ppb.UpsertProfileRequest{
				Profile: &ppb.OrgProfile{
					UserID:    "10000000-0000-0000-0000-000000000000",
					OrgID:     "20000000-0000-0000-0000-000000000000",
					OrgType:   ppb.OrgType_PetNet,
					Status:    ppb.Status_Pending,
					RiskScore: ppb.RiskScore_High,
					BusinessInfo: &ppb.BusinessInfo{
						CompanyName:   "company-name-u",
						StoreName:     "store-name-u",
						PhoneNumber:   "1234567891",
						FaxNumber:     "1234567891",
						Website:       "https://website-u.com",
						CompanyEmail:  "company@mail-u.com",
						ContactPerson: "contact-person-u",
						Position:      "position-u",
						Address: &ppb.Address{
							Address1:   "address1-u",
							City:       "cityu",
							State:      "stateu",
							PostalCode: "123450",
						},
					},
					AccountInfo: &ppb.AccountInfo{
						Bank:                    "bank-u",
						BankAccountNumber:       "bank-acct-no-u",
						BankAccountHolder:       "bank-acct-holder-u",
						AgreeTermsConditions:    ppb.Boolean_False,
						AgreeOnlineSupplierForm: ppb.Boolean_False,
						Currency:                ppb.Currency_SGD,
					},
					ReminderSent:      ppb.Boolean_False,
					DateApplied:       ts2,
					Deleted:           ts2,
					DsaCode:           "1234",
					TerminalIdOtc:     "4567",
					TerminalIdDigital: "7890",
					IsProvider:        false,
				},
			},
			reqq: &ppb.UpdateOrgProfileUserIDRequest{
				OldOrgID: "20000000-0000-0000-0000-000000000000",
				NewOrgID: "40000000-0000-0000-0000-000000000004",
				UserID:   "10000000-0000-0000-0000-000000000000",
			},
		},
		{
			desc:    "Missing orgID",
			wantErr: true,
			req: &ppb.UpsertProfileRequest{
				Profile: &ppb.OrgProfile{
					UserID: "10000000-0000-0000-0000-000000000000",
				},
			},
		},
		{
			desc:    "Email invalid",
			wantErr: true,
			req: &ppb.UpsertProfileRequest{
				Profile: &ppb.OrgProfile{
					UserID: "10000000-0000-0000-0000-000000000000",
					OrgID:  "20000000-0000-0000-0000-000000000000",
					BusinessInfo: &ppb.BusinessInfo{
						CompanyEmail: "companymail.com",
					},
				},
			},
		},
	}
	test.NewNullLogger()

	o := cmp.Options{
		cmpopts.IgnoreUnexported(
			ppb.GetProfileResponse{},
			ppb.OrgProfile{},
			ppb.BusinessInfo{},
			ppb.AccountInfo{},
			ppb.Address{},
		),
		cmpopts.IgnoreFields(ppb.OrgProfile{}, "ID", "Created", "Updated", "Deleted", "DateApplied"),
	}
	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			st, cleanup := postgres.NewTestStorage(os.Getenv("DATABASE_CONNECTION"), filepath.Join("..", "..", "migrations", "sql"))
			t.Cleanup(cleanup)

			h := New(st, Mock{}, profile.New(st))
			_, err := h.UpsertProfile(ctx, test.req)
			if err != nil && !test.wantErr {
				t.Fatal("h.UpsertProfile: ", err)
			}
			if err == nil && test.wantErr {
				t.Fatal("h.UpsertProfile: want error got nil")
			}
			if err != nil && test.wantErr {
				return
			}

			got, err := h.GetProfile(ctx, &ppb.GetProfileRequest{
				OrgID: test.req.Profile.OrgID,
			})
			if err != nil {
				t.Fatal("h.GetProfile: ", err)
			}
			if !cmp.Equal(test.req.Profile, got.Profile, o) {
				t.Fatal("GetProfile (-want +got): ", cmp.Diff(test.req.Profile, got.Profile, o))
			}

			_, errDC := h.GetProfileByDsaCode(ctx, &ppb.GetProfileByDsaCodeRequest{
				DsaCode: test.req.Profile.DsaCode,
			})
			if errDC != nil {
				t.Fatal("h.GetProfileByDsaCode: ", errDC)
			}

			_, err = h.UpsertProfile(ctx, test.ureq)
			if err != nil && !test.wantErr {
				t.Fatal("h.UpsertProfile: ", err)
			}

			got, err = h.GetProfile(ctx, &ppb.GetProfileRequest{
				OrgID: test.ureq.Profile.OrgID,
			})
			if err != nil {
				t.Fatal("h.GetProfile: ", err)
			}
			if !cmp.Equal(test.ureq.Profile, got.Profile, o) {
				t.Fatal("GetProfile (-want +got): ", cmp.Diff(test.ureq.Profile, got.Profile, o))
			}

			_, err = h.UpdateOrgProfileUserID(ctx, test.reqq)
			if err != nil {
				t.Fatal("h.UpdateOrgProfileUserID: ", err)
			}

			_, errDC = h.GetProfileByDsaCode(ctx, &ppb.GetProfileByDsaCodeRequest{
				DsaCode: test.ureq.Profile.DsaCode,
			})
			if errDC != nil {
				t.Fatal("h.GetProfile: ", errDC)
			}
			if !cmp.Equal(test.ureq.Profile, got.Profile, o) {
				t.Fatal("GetProfile (-want +got): ", cmp.Diff(test.ureq.Profile, got.Profile, o))
			}

			_, err = h.ListProfiles(ctx, &ppb.ListProfilesRequest{})
			if err != nil {
				if !test.wantErr {
					t.Fatal("h.ListProfiles: ", err)
				}
				return
			}
		})
	}
}
