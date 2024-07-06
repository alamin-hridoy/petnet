package partner

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/metadata"

	"brank.as/petnet/api/core/partner/ayannah"
	"brank.as/petnet/api/core/partner/bpi"
	cebuanaint "brank.as/petnet/api/core/partner/cebint"
	"brank.as/petnet/api/core/partner/cebuana"
	"brank.as/petnet/api/core/partner/instacash"
	"brank.as/petnet/api/core/partner/intelexpress"
	"brank.as/petnet/api/core/partner/iremit"
	"brank.as/petnet/api/core/partner/japanremit"
	"brank.as/petnet/api/core/partner/metrobank"
	"brank.as/petnet/api/core/partner/perahubremit"
	"brank.as/petnet/api/core/partner/remitly"
	"brank.as/petnet/api/core/partner/ria"
	"brank.as/petnet/api/core/partner/transfast"
	"brank.as/petnet/api/core/partner/uniteller"
	"brank.as/petnet/api/core/partner/ussc"
	"brank.as/petnet/api/core/partner/wise"
	"brank.as/petnet/api/core/static"
	"brank.as/petnet/api/storage"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"brank.as/petnet/serviceutil/auth/hydra"
)

func TestInputGuide(t *testing.T) {
	st := newTestStorage(t)

	tests := []struct {
		desc string
		in   *ppb.InputGuideRequest
		want *ppb.InputGuideResponse
	}{
		{
			desc: "Remitly Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.RMCode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGIDsLabel: {
						Field: "receiver.identification.type",
						Inputs: []*ppb.Input{
							{
								Value: "GOVERNMENT_ISSUED_ID",
								Name:  "AFP ID",
							},
							{
								Value: "DRIVERS_LICENSE",
								Name:  "Driver License",
							},
						},
					},
					storage.IGOtherInfoLabel: {
						OtherInfo: remitly.OtherInfoFieldsRM,
					},
				},
			},
		},
		{
			desc: "Transfast Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.TFCode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGIDsLabel: {
						Field: "receiver.identification.type",
						Inputs: []*ppb.Input{
							{
								Value:       "1",
								Name:        "FAMILY MAINTENANCE",
								CountryCode: "PH",
							},
							{
								Value:       "2",
								Name:        "EDUCATION",
								CountryCode: "PH",
							},
						},
					},
					storage.IGRelationsLabel: {
						Field: "receiver.relationship",
						Inputs: []*ppb.Input{
							{
								Value: "8",
								Name:  "FRIEND",
							},
							{
								Value: "16",
								Name:  "Self",
							},
						},
					},
					storage.IGOccupsLabel: {
						Field: "receiver.employment.occupation_id,receiver.employment.occupation",
						Inputs: []*ppb.Input{
							{
								Value: "1",
								Name:  "HOUSEWIFE",
							},
							{
								Value: "2",
								Name:  "STUDENT",
							},
						},
					},
					storage.IGPurposesLabel: {
						Field: "receiver.transaction_purpose",
						Inputs: []*ppb.Input{
							{
								Value:       "1",
								Name:        "FAMILY MAINTENANCE",
								CountryCode: "PH",
							},
							{
								Value:       "2",
								Name:        "EDUCATION",
								CountryCode: "PH",
							},
						},
					},
					storage.IGOtherInfoLabel: {
						OtherInfo: transfast.OtherInfoFieldsTF,
					},
				},
			},
		},
		{
			desc: "TransferWise Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.WISECode,
				CountryCode:  "US",
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGCountryLabel: {
						Field: "receiver.country",
						Inputs: []*ppb.Input{
							{
								Value: "AX",
								Name:  "Åland Islands",
							},
							{
								Value: "AL",
								Name:  "Albania",
							},
						},
					},
					storage.IGStateLabel: {
						Field: "receiver.state",
						Inputs: []*ppb.Input{
							{
								Value: "AL",
								Name:  "Alabama",
							},
							{
								Value: "AK",
								Name:  "Alaska",
							},
						},
					},
					storage.IGCurrencyLabel: {
						Field: "receiver.currency",
						Inputs: []*ppb.Input{
							{
								Value:       "AED",
								Description: "UAE Dirham",
							},
							{
								Value:       "ARS",
								Description: "Argentine peso",
							},
						},
					},
					storage.IGOtherInfoLabel: {
						OtherInfo: wise.OtherInfoFieldsWISE,
					},
				},
			},
		},
		{
			desc: "uniteller Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.UNTCode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGCountryLabel: {
						Field: "receiver.country",
						Inputs: []*ppb.Input{
							{
								Value: "AS",
								Name:  "AMERICAN SAMOA",
							},
							{
								Value: "AU",
								Name:  "AUSTRALIAN",
							},
						},
					},
					storage.IGCurrencyLabel: {
						Field: "receiver.currency",
						Inputs: []*ppb.Input{
							{
								Value: "USD",
								Name:  "US DOLLAR",
							},
							{
								Value: "PHP",
								Name:  "PHILIPPINE PESOS",
							},
						},
					},
					storage.IGOccupsLabel: {
						Field: "receiver.occupation",
						Inputs: []*ppb.Input{
							{
								Value: "OTH",
								Name:  "OTHER",
							},
							{
								Value: "HW",
								Name:  "HOUSEWIFE",
							},
						},
					},
					storage.IGIDsLabel: {
						Field: "receiver.identification.type",
						Inputs: []*ppb.Input{
							{
								Value: "PHILIPPINES",
								Name:  "LICENSE",
							},
							{
								Value: "PHILIPPINES",
								Name:  "PASSPORT",
							},
						},
					},
					storage.IGStateLabel: {
						Field: "receiver.states",
						Inputs: []*ppb.Input{
							{
								Value:       "PH-DVO",
								StateName:   "DAVAO OCCIDENTAL",
								CountryName: "PHILIPPINES",
							},
							{
								StateName:   "PH-NA",
								Value:       "PH-NA",
								CountryName: "PHILIPPINES",
							},
						},
					},
					storage.IGUsStateLabel: {
						Field: "receiver.usa-states",
						Inputs: []*ppb.Input{
							{
								Value:       "NJ",
								StateName:   "NEW JERSEY",
								CountryName: "USA",
							},
							{
								Value:       "PA",
								StateName:   "PENNSYLVANIA",
								CountryName: "USA",
							},
						},
					},
					storage.IGOtherInfoLabel: {
						OtherInfo: uniteller.OtherInfoFieldsUNT,
					},
				},
			},
		},
		{
			desc: "cebuana Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.CEBCode,
				AgentCode:    "01030063",
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGCountryLabel: {
						Field: "*country",
						Inputs: []*ppb.Input{
							{
								Value:       "1",
								Name:        "AFGHANISTAN",
								CountryCode: "AF",
							},
							{
								Value:       "2",
								Name:        "ALBANIA",
								CountryCode: "AL",
							},
						},
					},
					storage.IGCurrencyLabel: {
						Field: "*currency",
						Inputs: []*ppb.Input{
							{
								Value:        "6",
								CurrencyCode: "Php",
								Description:  "PHILIPPINE PESO",
							},
						},
					},
					storage.IGFundsLabel: {
						Field: "receiver.source_funds,sender.source_funds",
						Inputs: []*ppb.Input{
							{
								Value: "1",
								Name:  "Employed",
							},
							{
								Value: "2",
								Name:  "Self-Employed",
							},
						},
					},
					storage.IGIDsLabel: {
						Field: "receiver.identification.type,sender.identification.type",
						Inputs: []*ppb.Input{
							{
								Value:       "1",
								Name:        "CR",
								Description: "Bank Credit Card",
							},
							{
								Value:       "3",
								Name:        "BY",
								Description: "Barangay ID",
							},
						},
					},
					storage.IGOtherInfoLabel: {
						OtherInfo: cebuana.OtherInfoFieldsCEB,
					},
				},
			},
		},
		{
			desc: "ussc Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.USSCCode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGOtherInfoLabel: {
						OtherInfo: ussc.OtherInfoFieldsUSSC,
					},
				},
			},
		},
		{
			desc: "iremit Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.IRCode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGOtherInfoLabel: {
						OtherInfo: iremit.OtherInfoFieldsIR,
					},
				},
			},
		},
		{
			desc: "ria Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.RIACode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGOtherInfoLabel: {
						OtherInfo: ria.OtherInfoFieldsRIA,
					},
				},
			},
		},
		{
			desc: "metrobank Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.MBCode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGOtherInfoLabel: {
						OtherInfo: metrobank.OtherInfoFieldsMB,
					},
				},
			},
		},
		{
			desc: "bpi Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.BPICode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGOtherInfoLabel: {
						OtherInfo: bpi.OtherInfoFieldsBPI,
					},
				},
			},
		},
		{
			desc: "instacash Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.ICCode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGOtherInfoLabel: {
						OtherInfo: instacash.OtherInfoFieldsIC,
					},
				},
			},
		},
		{
			desc: "japan remit Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.JPRCode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGOtherInfoLabel: {
						OtherInfo: japanremit.OtherInfoFieldsJPR,
					},
				},
			},
		},
		{
			desc: "ayannah Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.AYACode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGOtherInfoLabel: {
						OtherInfo: ayannah.OtherInfoFieldsAYA,
					},
				},
			},
		},
		{
			desc: "cebint Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.CEBINTCode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGOtherInfoLabel: {
						OtherInfo: cebuanaint.OtherInfoFieldsCEBINT,
					},
				},
			},
		},
		{
			desc: "intel express Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.IECode,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGOtherInfoLabel: {
						OtherInfo: intelexpress.OtherInfoFieldsIE,
					},
				},
			},
		},
		{
			desc: "PerahubRemit Success",
			in: &ppb.InputGuideRequest{
				RemitPartner: static.PerahubRemit,
				City:         "city1",
				ID:           1,
			},
			want: &ppb.InputGuideResponse{
				InputGuide: map[string]*ppb.Guide{
					storage.IGProvincesCityLabel: {
						Inputs: []*ppb.Input{
							{
								Value:     "MANILA",
								Name:      "MANILA",
								StateName: "METRO MANILA",
							},
							{
								Value:     "CITY OF MAKATI",
								Name:      "CITY OF MAKATI",
								StateName: "METRO MANILA",
							},
							{
								Value:     "CITY OF MUNTINLUPA",
								Name:      "CITY OF MUNTINLUPA",
								StateName: "METRO MANILA",
							},
							{
								Value:     "CITY OF PARANAQUE",
								Name:      "CITY OF PARANAQUE",
								StateName: "METRO MANILA",
							},
							{
								Value:     "PASAY CITY",
								Name:      "PASAY CITY",
								StateName: "METRO MANILA",
							},
						},
					},
					storage.IGBrgyLabel: {
						Inputs: []*ppb.Input{
							{
								Name:  "Barangay 1",
								Value: "1013",
							},
						},
					},
					storage.IGPurposesLabel: {
						Inputs: []*ppb.Input{
							{
								Value: "Family Support/Living Expenses",
								Name:  "Family Support/Living Expenses",
							},
							{
								Value: "Saving/Investments",
								Name:  "Saving/Investments",
							},
							{
								Value: "Gift",
								Name:  "Gift",
							},
						},
					},
					storage.IGRelationsLabel: {
						Inputs: []*ppb.Input{
							{
								Value: "Family",
								Name:  "Family",
							},
							{
								Value: "Friend",
								Name:  "Friend",
							},
						},
					},
					storage.IGPartnerLabel: {
						Inputs: []*ppb.Input{
							{
								Value: "DRP",
								Name:  "BRANKAS",
							},
						},
					},
					storage.IGOccupsLabel: {
						Inputs: []*ppb.Input{
							{
								Value: "Airline/Maritime Employee",
								Name:  "Airline/Maritime Employee",
							},
							{
								Value: "Art/Entertainment/Media/Sports Professional",
								Name:  "Art/Entertainment/Media/Sports Professional",
							},
							{
								Value: "Civil/Government Employee",
								Name:  "Civil/Government Employee",
							},
							{
								Value: "Domestic Helper",
								Name:  "Domestic Helper",
							},
							{
								Value: "Driver",
								Name:  "Driver",
							},
						},
					},
					storage.IGFundsLabel: {
						Inputs: []*ppb.Input{
							{
								Value: "Salary",
								Name:  "Salary",
							},
							{
								Value: "Savings",
								Name:  "Savings",
							},
							{
								Value: "Borrowed Funds/Loan",
								Name:  "Borrowed Funds/Loan",
							},
							{
								Value: "Pension/Government/Welfare",
								Name:  "Pension/Government/Welfare",
							},
							{
								Value: "Gift",
								Name:  "Gift",
							},
						},
					},
					storage.IGEmploymentLabel: {
						Inputs: []*ppb.Input{
							{
								Value: "Administrative/Human Resources",
								Name:  "Administrative/Human Resources",
							},
							{
								Value: "Agriculture",
								Name:  "Agriculture",
							},
							{
								Value: "Banking /Financial Services",
								Name:  "Banking /Financial Services",
							},
							{
								Value: "Computer and Information Tech Services",
								Name:  "Computer and Information Tech Services",
							},
							{
								Value: "Construction/Contractors",
								Name:  "Construction/Contractors",
							},
						},
					},
					storage.IGOtherInfoLabel: {
						OtherInfo: perahubremit.OtherInfoFieldsPerahubRemit,
					},
				},
			},
		},
	}

	vs := NewValidators()
	for _, v := range vs {
		testExists := false
		for _, test := range tests {
			if v.Kind() == test.in.RemitPartner {
				testExists = true
			}
		}
		if !testExists && v.Kind() != static.WUCode {
			t.Fatal("add input guide testcase for partner: ", v.Kind())
		}
	}

	uid := uuid.New().String()
	nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, "owner", uid))
	ctx := nmd.ToIncoming(context.Background())
	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h, _ := newTestSvc(t, st)
			got, err := h.InputGuide(ctx, test.in)
			if err != nil {
				t.Fatal(err)
			}
			o := cmp.Options{
				cmpopts.IgnoreUnexported(
					ppb.InputGuideResponse{},
					ppb.Guide{},
					ppb.Input{},
				),
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
			}
		})
	}
}
