package partner

import (
	"context"
	"testing"

	"brank.as/petnet/api/core/static"
	ppb "brank.as/petnet/gunk/drp/v1/partner"
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	pl "brank.as/petnet/gunk/dsa/v2/partnerlist"
	pfSvc "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/serviceutil/auth/hydra"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
	"github.com/grpc-ecosystem/go-grpc-middleware/util/metautils"
	"google.golang.org/grpc/metadata"
)

func TestList(t *testing.T) {
	st := newTestStorage(t)

	allPtnr := &ppb.RemitPartnersResponse{
		Partners: map[string]*ppb.RemitPartner{
			static.WUCode: {
				PartnerCode: static.WUCode,
				PartnerName: "Western Union",
				SupportedSendTypes: map[string]*ppb.RemitType{
					"Send": {
						Code:        "Send",
						Description: "Send money to an individual for pick up at any Western Union location.",
					},
					"Direct": {
						Code:        "Direct",
						Description: "Send money directly to a bank account.",
					},
					"Mobile": {
						Code:        "Mobile",
						Description: "Send money to a mobile number.",
					},
					"QuickPay": {
						Code:        "QuickPay",
						Description: "Send money to a buiness that accepts Western Union QuickPay.",
					},
				},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.IRCode: {
				PartnerCode:        static.IRCode,
				PartnerName:        "iRemit",
				SupportedSendTypes: map[string]*ppb.RemitType{},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.TFCode: {
				PartnerCode:        static.TFCode,
				PartnerName:        "Transfast",
				SupportedSendTypes: map[string]*ppb.RemitType{},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.RMCode: {
				PartnerCode:        static.RMCode,
				PartnerName:        "Remitly",
				SupportedSendTypes: map[string]*ppb.RemitType{},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.RIACode: {
				PartnerCode:        static.RIACode,
				PartnerName:        "Ria",
				SupportedSendTypes: map[string]*ppb.RemitType{},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.MBCode: {
				PartnerCode:        static.MBCode,
				PartnerName:        "Metrobank",
				SupportedSendTypes: map[string]*ppb.RemitType{},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.BPICode: {
				PartnerCode:        static.BPICode,
				PartnerName:        "BPI",
				SupportedSendTypes: map[string]*ppb.RemitType{},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.USSCCode: {
				PartnerCode: static.USSCCode,
				PartnerName: "USSC",
				SupportedSendTypes: map[string]*ppb.RemitType{
					"Send": {
						Code:        "Send",
						Description: "Send money transaction.",
					},
				},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.WISECode: {
				PartnerCode: static.WISECode,
				PartnerName: "TransferWise",
				SupportedSendTypes: map[string]*ppb.RemitType{
					"Send": {
						Code:        "Send",
						Description: "Send money transaction.",
					},
				},
				SupportedDisburseTypes: map[string]*ppb.RemitType{},
			},
			static.ICCode: {
				PartnerCode:        static.ICCode,
				PartnerName:        "InstaCash",
				SupportedSendTypes: map[string]*ppb.RemitType{},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.JPRCode: {
				PartnerCode:        static.JPRCode,
				PartnerName:        "JapanRemit",
				SupportedSendTypes: map[string]*ppb.RemitType{},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.UNTCode: {
				PartnerCode:        static.UNTCode,
				PartnerName:        "Uniteller",
				SupportedSendTypes: map[string]*ppb.RemitType{},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.CEBCode: {
				PartnerCode: static.CEBCode,
				PartnerName: "Cebuana",
				SupportedSendTypes: map[string]*ppb.RemitType{
					"Send": {
						Code:        "Send",
						Description: "Send money transaction.",
					},
				},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.CEBINTCode: {
				PartnerCode:        static.CEBINTCode,
				PartnerName:        "Cebuana intl",
				SupportedSendTypes: map[string]*ppb.RemitType{},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.AYACode: {
				PartnerCode: static.AYACode,
				PartnerName: "Ayannah",
				SupportedSendTypes: map[string]*ppb.RemitType{
					"Send": {
						Code:        "Send",
						Description: "Send money transaction.",
					},
				},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.IECode: {
				PartnerCode: static.IECode,
				PartnerName: "IntelExpress",
				SupportedSendTypes: map[string]*ppb.RemitType{
					"Send": {
						Code:        "Send",
						Description: "Send money transaction.",
					},
				},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
			static.PerahubRemit: {
				PartnerCode: static.PerahubRemit,
				PartnerName: "PerahubRemit",
				SupportedSendTypes: map[string]*ppb.RemitType{
					"Send": {
						Code:        "Send",
						Description: "Send money transaction.",
					},
				},
				SupportedDisburseTypes: map[string]*ppb.RemitType{
					"Payout": {
						Code:        "Payout",
						Description: "Payout a transaction.",
					},
				},
			},
		},
	}

	tests := []struct {
		desc string
		env  string
		in   *ppb.RemitPartnersRequest
		want *ppb.RemitPartnersResponse
	}{
		{
			desc: "Live All Enabled",
			env:  "live",
			in: &ppb.RemitPartnersRequest{
				Country: "PH",
			},
			want: allPtnr,
		},
		{
			desc: "Live WU IR Enabled",
			env:  "live",
			in: &ppb.RemitPartnersRequest{
				Country: "PH",
			},
			want: &ppb.RemitPartnersResponse{
				Partners: map[string]*ppb.RemitPartner{
					static.WUCode: {
						PartnerCode: static.WUCode,
						PartnerName: "Western Union",
						SupportedSendTypes: map[string]*ppb.RemitType{
							"Send": {
								Code:        "Send",
								Description: "Send money to an individual for pick up at any Western Union location.",
							},
							"Direct": {
								Code:        "Direct",
								Description: "Send money directly to a bank account.",
							},
							"Mobile": {
								Code:        "Mobile",
								Description: "Send money to a mobile number.",
							},
							"QuickPay": {
								Code:        "QuickPay",
								Description: "Send money to a buiness that accepts Western Union QuickPay.",
							},
						},
						SupportedDisburseTypes: map[string]*ppb.RemitType{
							"Payout": {
								Code:        "Payout",
								Description: "Payout a transaction.",
							},
						},
					},
					static.IRCode: {
						PartnerCode:        static.IRCode,
						PartnerName:        "iRemit",
						SupportedSendTypes: map[string]*ppb.RemitType{},
						SupportedDisburseTypes: map[string]*ppb.RemitType{
							"Payout": {
								Code:        "Payout",
								Description: "Payout a transaction.",
							},
						},
					},
				},
			},
		},
		{
			desc: "Live No Partners Enabled",
			env:  "live",
			in: &ppb.RemitPartnersRequest{
				Country: "PH",
			},
			want: &ppb.RemitPartnersResponse{
				Partners: map[string]*ppb.RemitPartner{},
			},
		},
		{
			desc: "Sandbox No Partners Enabled",
			env:  "sandbox",
			in: &ppb.RemitPartnersRequest{
				Country: "PH",
			},
			want: &ppb.RemitPartnersResponse{
				Partners: map[string]*ppb.RemitPartner{},
			},
		},
	}
	uid := uuid.New().String()
	oid := uuid.New().String()

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			nmd := metautils.NiceMD(metadata.Pairs(hydra.ClientIDKey, uid, hydra.OrgIDKey, oid, "owner", uid, "environment", test.env))
			ctx := nmd.ToIncoming(context.Background())

			h, m := newTestSvc(t, st)
			if test.env == "live" {
				m.SetGetPartnersResp(mockPartnerList(test.want))
				m.SetGetPartnerListResponse(mockPartnerListGlobally(test.want))
			} else {
				m.SetGetPartnersResp(mockPartnerList(&ppb.RemitPartnersResponse{
					Partners: map[string]*ppb.RemitPartner{},
				}))
				m.SetGetPartnerListResponse(mockPartnerListGlobally(test.want))
			}
			got, err := h.RemitPartners(ctx, test.in)
			if err != nil {
				t.Fatal(err)
			}
			o := cmp.Options{
				cmpopts.IgnoreUnexported(
					ppb.RemitPartnersResponse{},
					ppb.RemitPartner{},
					ppb.RemitType{},
				),
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
			}
		})
	}
}

func mockPartnerList(ptnr *ppb.RemitPartnersResponse) *pfSvc.ListServiceRequestResponse {
	var svc []*pfSvc.ServiceRequest
	for _, v := range ptnr.GetPartners() {
		svc = append(svc, &pfSvc.ServiceRequest{
			Partner: v.GetPartnerCode(),
			Status:  pfSvc.ServiceRequestStatus_ACCEPTED,
			Enabled: true,
		})
	}

	return &pfSvc.ListServiceRequestResponse{
		ServiceRequst: svc,
		Total:         int32(len(svc)),
	}
}

func mockPartnerListGlobally(ptnr *ppb.RemitPartnersResponse) *pl.GetPartnerListResponse {
	uid := uuid.New().String()
	var svc []*pl.PartnerList
	for _, v := range ptnr.GetPartners() {
		svc = append(svc, &pl.PartnerList{
			ID:     uid,
			Stype:  v.GetPartnerCode(),
			Status: spb.PartnerStatusType_ENABLED.String(),
		})
	}

	return &pl.GetPartnerListResponse{
		PartnerList: svc,
	}
}
