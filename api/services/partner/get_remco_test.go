package partner

import (
	"context"
	"testing"

	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestGetPartnersRemco(t *testing.T) {
	st := newTestStorage(t)
	ctx := context.Background()
	o := cmp.Options{
		cmpopts.IgnoreUnexported(
			ppb.GetPartnersRemcoResponse{},
			ppb.PerahubGetRemcoIDResult{},
		),
	}
	tests := []struct {
		desc string
		env  string
		want *ppb.GetPartnersRemcoResponse
	}{
		{
			desc: "get partners remco id",
			want: &ppb.GetPartnersRemcoResponse{
				Code:    "200",
				Message: "Good",
				Result: []*ppb.PerahubGetRemcoIDResult{
					{
						ID:   "1",
						Name: "iRemit",
					},
					{
						ID:   "2",
						Name: "BPI",
					},
				},
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			h, _ := newTestSvc(t, st)
			got, err := h.GetPartnersRemco(ctx, nil)
			if err != nil {
				t.Fatal(err)
			}
			if !cmp.Equal(test.want, got, o) {
				t.Error("(-want +got): ", cmp.Diff(test.want, got, o))
			}
		})
	}
}
