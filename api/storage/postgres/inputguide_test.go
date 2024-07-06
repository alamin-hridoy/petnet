package postgres

import (
	"context"
	"testing"

	"brank.as/petnet/api/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func TestInputGuide(t *testing.T) {
	ts := newTestStorage(t)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		in := &storage.InputGuide{
			Partner: "RM",
			Data: storage.InputGuideData{
				storage.IGIDsLabel: {
					{
						Value:       "AFP ID",
						Name:        "GOVERNMENT_ISSUED_ID",
						Description: "descr",
					},
				},
			},
		}
		_, err := ts.CreateInputGuide(context.TODO(), *in)
		if err != nil {
			t.Fatalf("CreateInputGuide() = got error %v, want nil", err)
		}

		got, err := ts.GetInputGuide(context.TODO(), in.Partner)
		if err != nil {
			t.Fatalf("GetInputGuide() = got error %v, want nil", err)
		}
		o := []cmp.Option{
			cmpopts.IgnoreFields(storage.InputGuide{}, "Created", "Updated"),
		}
		if !cmp.Equal(in, got, o...) {
			t.Error(cmp.Diff(in, got, o...))
		}
		if got.Created.IsZero() || got.Updated.IsZero() {
			t.Error("created and updated shouldn't be empty")
		}
		updated := got.Updated

		in.Data = storage.InputGuideData{
			storage.IGIDsGroup: {
				{
					Value:       "code updated",
					Name:        "name updated",
					Description: "desc updated",
				},
			},
		}
		got, err = ts.UpdateInputGuide(context.TODO(), *in)
		if err != nil {
			t.Fatalf("UpdateInputGuide() = got error %v, want nil", err)
		}
		if got.Updated.IsZero() {
			t.Fatal("updated field is zero")
		}
		if got.Updated == updated {
			t.Fatal("update field hasn't changed after update")
		}

		got, err = ts.GetInputGuide(context.TODO(), in.Partner)
		if err != nil {
			t.Fatalf("GetInputGuide() = got error %v, want nil", err)
		}
		if !cmp.Equal(in, got, o...) {
			t.Error(cmp.Diff(in, got, o...))
		}
	})
}
