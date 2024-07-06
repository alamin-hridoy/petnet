package postgres

import (
	"context"
	"testing"

	"brank.as/rbac/usermgm/storage"
)

func TestCreateOrganization(t *testing.T) {
	ts := newTestStorage(t)

	t.Run("Valid", func(t *testing.T) {
		t.Parallel()
		in := &storage.Organization{
			OrgName: "TestOrg-CompanyName",
		}
		orgID, err := ts.CreateOrg(context.TODO(), *in)
		if err != nil {
			t.Fatalf("CreateOrganization() = got error %v, want nil", err)
		}
		if orgID == "" {
			t.Fatal("CreateOrganization() = returned empty ID")
		}
	})
}

func TestGetOrg(t *testing.T) {
	ts := newTestStorage(t)

	t.Run("Valid-ByID", func(t *testing.T) {
		t.Parallel()
		in := &storage.Organization{
			OrgName: "TestOrg-CompanyName",
		}
		orgID, err := ts.CreateOrg(context.TODO(), *in)
		if err != nil {
			t.Fatalf("CreateOrganization() = got error %+v, want nil", err)
		}
		org, err := ts.GetOrgByID(context.TODO(), orgID)
		if err != nil {
			t.Fatalf("GetOrgByID() = got error %v, want nil", err)
		}
		if org.OrgName != in.OrgName {
			t.Fatalf("GetOrgByID() = got bad company name %v, want %v", org.OrgName, in.OrgName)
		}
		if org.Active {
			t.Fatalf("GetOrgByID() = new organizations should be disabled by default")
		}
	})

	t.Run("Invalid-ByID", func(t *testing.T) {
		t.Parallel()
		const badUUID string = "12345678-aaaa-bbbb-cccc-1234567890ab"
		_, err := ts.GetOrgByID(context.TODO(), badUUID)
		if err != storage.NotFound {
			t.Fatalf("GetOrgByID() = got error %v, want storage.ErrNotFound", err)
		}
	})
}

func TestGetOrgs(t *testing.T) {
	ts := newTestStorage(t)

	t.Run("Valid-ByID", func(t *testing.T) {
		t.Parallel()
		in := []storage.Organization{
			{OrgName: "TestOrg-CompanyName"},
			{
				OrgName: "TestOrg-CompanyName-2",
			},
		}
		var err error
		for i := range in {
			i := len(in) - i - 1 // insert in reverse order to match sorting
			in[i].ID, err = ts.CreateOrg(context.TODO(), in[i])
			if err != nil {
				t.Fatalf("CreateOrganization() = got error %+v, want nil", err)
			}
		}
		orgs, err := ts.GetOrgs(context.TODO())
		if err != nil {
			t.Fatalf("GetOrgs() = got error %v, want nil", err)
		}
		for i := range in {
			org := orgs[i]
			want := in[i]
			if org.OrgName != want.OrgName {
				t.Errorf("GetOrgs() = got bad company name %v, want %v", org.OrgName, want.OrgName)
			}
			if org.Active {
				t.Errorf("GetOrgs() = new organizations should be disabled by default")
			}
		}
	})
}

func TestUpdateOrg(t *testing.T) {
	ts := newTestStorage(t)

	t.Run("Valid-UpdateByID", func(t *testing.T) {
		t.Parallel()
		in := &storage.Organization{
			OrgName: "TestOrg-CompanyName",
		}
		orgID, err := ts.CreateOrg(context.TODO(), *in)
		if err != nil {
			t.Fatalf("CreateOrganization() = got error %+v, want nil", err)
		}
		const newCompanyName = "TestOrg-NewCompanyName"
		const totalEmployess = 100
		in.OrgName = newCompanyName
		in.ID = orgID
		in.Active = true

		_, err = ts.UpdateOrgByID(context.TODO(), *in)
		if err != nil {
			t.Fatalf("UpdateOrgByID() = got error %v, want nil", err)
		}
		updatedOrg, err := ts.GetOrgByID(context.TODO(), orgID)
		if err != nil {
			t.Fatalf("GetOrgByID() = got error %v, want nil", err)
		}
		if updatedOrg.OrgName != newCompanyName {
			t.Fatalf("UpdateOrgByID() = got bad company name %v, want %v", updatedOrg.OrgName, newCompanyName)
		}
		if !updatedOrg.Active {
			t.Fatalf("UpdateOrgByID() = got bad company isActive value")
		}
	})

	t.Run("Invalid-UpdateByID", func(t *testing.T) {
		t.Parallel()
		const badUUID string = "12345678-aaaa-bbbb-cccc-1234567890ab"
		in := &storage.Organization{
			ID:      badUUID,
			OrgName: "TestOrg-CompanyName",
		}
		_, err := ts.UpdateOrgByID(context.TODO(), *in)
		if err != storage.NotFound {
			t.Fatalf("UpdateOrgByID() = got error %v, want storage.ErrNotFound", err)
		}
	})
}
