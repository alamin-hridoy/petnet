package partners

import (
	"context"

	"github.com/sirupsen/logrus"

	sVcpb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
	"brank.as/petnet/profile/storage/postgres"
	"brank.as/petnet/svcutil/partners"
)

func BootstrapAdminPartners(ctx context.Context, log *logrus.Entry, st *postgres.Storage) error {
	if len(partners.PartnersList) > 0 {
		for _, ptnr := range partners.PartnersList {
			lr, err := st.GetPartnerList(ctx, &storage.PartnerList{
				Stype:       ptnr.Stype,
				ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
			})
			if err != nil {
				if _, err := st.CreatePartnerList(ctx, ptnr); err != nil {
					log.Info("Bootstrapped Petnet Admin partner Create Failed")
				}
			}
			if len(lr) > 0 {
				if _, err := st.UpdatePartnerList(ctx, ptnr); err != nil {
					log.Info("Bootstrapped Petnet Admin partner update Failed")
				}
			} else {
				if _, err := st.CreatePartnerList(ctx, ptnr); err != nil {
					log.Info("Bootstrapped Petnet Admin partner Create Failed")
				}
			}
		}
	}
	log.Info("Bootstrapped Petnet Admin partners")
	return nil
}

func BootstrapAdminCicoPartners(ctx context.Context, log *logrus.Entry, st *postgres.Storage) error {
	if len(partners.CicoPartnersList) > 0 {
		for _, ptnr := range partners.CicoPartnersList {
			lr, err := st.GetPartnerList(ctx, &storage.PartnerList{
				Stype:       ptnr.Stype,
				ServiceName: sVcpb.ServiceType_CASHINCASHOUT.String(),
			})
			if err != nil {
				if _, err := st.CreatePartnerList(ctx, ptnr); err != nil {
					log.Info("Bootstrapped Petnet Admin cico partner Create Failed")
				}
			}
			if len(lr) > 0 {
				if _, err := st.UpdatePartnerList(ctx, ptnr); err != nil {
					log.Info("Bootstrapped Petnet Admin cico partner update Failed")
				}
			} else {
				if _, err := st.CreatePartnerList(ctx, ptnr); err != nil {
					log.Info("Bootstrapped Petnet Admin cico partner Create Failed")
				}
			}
		}
	}
	log.Info("Bootstrapped Petnet Admin cico partners")
	return nil
}

func BootstrapAdminRTAPartners(ctx context.Context, log *logrus.Entry, st *postgres.Storage) error {
	if len(partners.RTAPartnersList) > 0 {
		for _, ptnr := range partners.RTAPartnersList {
			lr, err := st.GetPartnerList(ctx, &storage.PartnerList{
				Stype:       ptnr.Stype,
				ServiceName: sVcpb.ServiceType_REMITTOACCOUNT.String(),
			})
			if err != nil {
				if _, err := st.CreatePartnerList(ctx, ptnr); err != nil {
					log.Info("Bootstrapped Petnet Admin RTA partner Create Failed")
				}
			}
			if len(lr) > 0 {
				if _, err := st.UpdatePartnerList(ctx, ptnr); err != nil {
					log.Info("Bootstrapped Petnet Admin RTA partner update Failed")
				}
			} else {
				if _, err := st.CreatePartnerList(ctx, ptnr); err != nil {
					log.Info("Bootstrapped Petnet Admin RTA partner Create Failed")
				}
			}
		}
	}
	log.Info("Bootstrapped Petnet Admin RTA partners")
	return nil
}
