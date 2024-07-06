package postgres_test

import (
	"context"
	"testing"
	"time"

	"brank.as/petnet/profile/storage"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/google/uuid"
)

// test for Question
func TestQuestion(t *testing.T) {
	ts := newTestStorage(t)
	ctx := context.Background()
	oid := uuid.NewString()
	uid := uuid.NewString()
	_, err := ts.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uid,
	})
	if err != nil {
		t.Fatal(err)
	}
	o := cmpopts.IgnoreFields(storage.Question{}, "ID", "Created", "Updated")
	fs := []storage.Question{
		{
			OrgID:   oid,
			UserID:  uid,
			QID:     "1",
			QType:   "1",
			ANS:     "1",
			Created: time.Now(),
			Updated: time.Now(),
		},
		{
			OrgID:   oid,
			UserID:  uid,
			QType:   "1",
			QID:     "2",
			ANS:     "2",
			Created: time.Now(),
			Updated: time.Now(),
		},
		{
			OrgID:   oid,
			UserID:  uid,
			QID:     "3",
			QType:   "1",
			ANS:     "3",
			Created: time.Now(),
			Updated: time.Now(),
		},
		{
			OrgID:   oid,
			UserID:  uid,
			QID:     "4",
			QType:   "1",
			ANS:     "4",
			Created: time.Now(),
			Updated: time.Now(),
		},
	}
	editfs := []storage.Question{
		{
			OrgID:   oid,
			UserID:  uid,
			QID:     "1",
			QType:   "1",
			ANS:     "2",
			Created: time.Now(),
			Updated: time.Now(),
		},
		{
			OrgID:   oid,
			UserID:  uid,
			QType:   "1",
			QID:     "2",
			ANS:     "1",
			Created: time.Now(),
			Updated: time.Now(),
		},
		{
			OrgID:   oid,
			UserID:  uid,
			QID:     "3",
			QType:   "1",
			ANS:     "1",
			Created: time.Now(),
			Updated: time.Now(),
		},
		{
			OrgID:   oid,
			UserID:  uid,
			QID:     "4",
			QType:   "1",
			ANS:     "2",
			Created: time.Now(),
			Updated: time.Now(),
		},
	}
	for _, v := range fs {
		got, err := ts.UpsertQuestion(ctx, &v)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(&v, got, o) {
			t.Error("(-want +got): ", cmp.Diff(v, got, o))
		}
	}
	for _, v := range editfs {
		got, err := ts.UpsertQuestion(ctx, &v)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(&v, got, o) {
			t.Error("(-want +got): ", cmp.Diff(v, got, o))
		}
	}
	gotList, err := ts.ListQuestion(ctx, &storage.Question{
		OrgID:  oid,
		UserID: uid,
	})
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(editfs, gotList, o) {
		t.Error("(-want +got): ", cmp.Diff(editfs, gotList, o))
	}
}

// test for MLTF Question
func TestMLTFQuestion(t *testing.T) {
	ts := newTestStorage(t)
	ctx := context.Background()
	oid := uuid.NewString()
	uid := uuid.NewString()
	_, err := ts.CreateOrgProfile(ctx, &storage.OrgProfile{
		OrgID:  oid,
		UserID: uid,
	})
	if err != nil {
		t.Fatal(err)
	}
	o := cmpopts.IgnoreFields(storage.Question{}, "ID", "Created", "Updated")
	fs := []storage.Question{
		{
			OrgID:          oid,
			UserID:         uid,
			QID:            "1",
			QType:          "1",
			CustomersTotal: "20",
			HrTotal:        "15",
			ImpactScore:    "12",
			Created:        time.Now(),
			Updated:        time.Now(),
		},
		{
			OrgID:          oid,
			UserID:         uid,
			QType:          "1",
			QID:            "2",
			CustomersTotal: "30",
			HrTotal:        "25",
			ImpactScore:    "22",
			Created:        time.Now(),
			Updated:        time.Now(),
		},
		{
			OrgID:          oid,
			UserID:         uid,
			QID:            "3",
			QType:          "1",
			CustomersTotal: "40",
			HrTotal:        "35",
			ImpactScore:    "32",
			Created:        time.Now(),
			Updated:        time.Now(),
		},
		{
			OrgID:          oid,
			UserID:         uid,
			QID:            "4",
			QType:          "1",
			CustomersTotal: "50",
			HrTotal:        "45",
			ImpactScore:    "42",
			Created:        time.Now(),
			Updated:        time.Now(),
		},
	}
	editfs := []storage.Question{
		{
			OrgID:          oid,
			UserID:         uid,
			QID:            "1",
			QType:          "1",
			CustomersTotal: "20",
			HrTotal:        "15",
			ImpactScore:    "12",
			Created:        time.Now(),
			Updated:        time.Now(),
		},
		{
			OrgID:          oid,
			UserID:         uid,
			QType:          "1",
			QID:            "2",
			CustomersTotal: "30",
			HrTotal:        "25",
			ImpactScore:    "22",
			Created:        time.Now(),
			Updated:        time.Now(),
		},
		{
			OrgID:          oid,
			UserID:         uid,
			QID:            "3",
			QType:          "1",
			CustomersTotal: "40",
			HrTotal:        "35",
			ImpactScore:    "32",
			Created:        time.Now(),
			Updated:        time.Now(),
		},
		{
			OrgID:          oid,
			UserID:         uid,
			QID:            "4",
			QType:          "1",
			CustomersTotal: "50",
			HrTotal:        "45",
			ImpactScore:    "42",
			Created:        time.Now(),
			Updated:        time.Now(),
		},
	}
	for _, v := range fs {
		got, err := ts.UpsertMlTfQuestion(ctx, &v)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(&v, got, o) {
			t.Error("(-want +got): ", cmp.Diff(v, got, o))
		}
	}
	for _, v := range editfs {
		got, err := ts.UpsertMlTfQuestion(ctx, &v)
		if err != nil {
			t.Error(err)
		}
		if !cmp.Equal(&v, got, o) {
			t.Error("(-want +got): ", cmp.Diff(v, got, o))
		}
	}
	gotList, err := ts.ListMlTfQuestion(ctx, &storage.Question{
		OrgID:  oid,
		UserID: uid,
	})
	if err != nil {
		t.Error(err)
	}
	if !cmp.Equal(editfs, gotList, o) {
		t.Error("(-want +got): ", cmp.Diff(editfs, gotList, o))
	}
}
