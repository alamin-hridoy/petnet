package partners

import (
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	sVcpb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
)

var PartnersList = []*storage.PartnerList{
	{
		Stype:       "WU",
		Name:        "Western Union",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "IR",
		Name:        "IRemit",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "TF",
		Name:        "Transfast",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "RIA",
		Name:        "Ria",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "MB",
		Name:        "MetroBank",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "RM",
		Name:        "Remitly",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "BPI",
		Name:        "BPI",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "USSC",
		Name:        "USSC",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "JPR",
		Name:        "JapanRemit",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "IC",
		Name:        "InstantCash",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "UNT",
		Name:        "Uniteller",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "WISE",
		Name:        "Transfer Wise",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "CEBI",
		Name:        "Cebuana Intl",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "AYA",
		Name:        "Ayannah",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "CEB",
		Name:        "Cebuana",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
	{
		Stype:       "IE",
		Name:        "IntelExpress",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTANCE.String(),
	},
}
