package permission

import (
	"context"

	"github.com/sirupsen/logrus"

	ppb "brank.as/petnet/gunk/dsa/v2/profile"
	"brank.as/petnet/serviceutil/logging"
	"brank.as/petnet/svcutil/permission"

	"brank.as/petnet/profile/storage"

	"brank.as/petnet/profile/storage/postgres"
	rbppb "brank.as/rbac/gunk/v1/permissions"
)

func BootstrapAdminPermissions(ctx context.Context, log *logrus.Entry, pecl rbppb.PermissionServiceClient, prcl rbppb.ProductServiceClient, st *postgres.Storage) error {
	var nps []*rbppb.ServicePermission
	for k, v := range permission.Permissions {
		np := &rbppb.ServicePermission{
			Name:        k,
			Description: v.Description,
			Actions:     v.Actions,
			Resource:    v.Resource,
		}
		nps = append(nps, np)
	}

	res, err := pecl.CreatePermission(ctx, &rbppb.CreatePermissionRequest{
		ServiceName: permission.AdminServiceName,
		Description: "The available permissions for a petnet admin to assign to roles",
		Permissions: nps,
	})
	if err != nil {
		logging.WithError(err, log).Error("creating permissions")
		return err
	}

	pfs, err := st.GetOrgProfiles(ctx, storage.FilterList{})
	if err != nil {
		logging.WithError(err, log).Error("getting org profiles")
		return err
	}
	var oid string
	for _, p := range pfs {
		if p.OrgType == int(ppb.OrgType_PetNet) {
			oid = p.OrgID
		}
	}

	if _, err := prcl.GrantService(ctx, &rbppb.GrantServiceRequest{
		ServiceID: res.GetServiceID(),
		OrgID:     oid,
	}); err != nil {
		logging.WithError(err, log).Error("granting service")
		return err
	}

	log.Info("Bootstrapped Petnet Admin permissions")
	return nil
}
