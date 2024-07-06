package partners

import (
	spb "brank.as/petnet/gunk/dsa/v2/partner"
	sVcpb "brank.as/petnet/gunk/dsa/v2/service"
	"brank.as/petnet/profile/storage"
)

var CicoPartnersList = []*storage.PartnerList{
	{
		Stype:       "GCASH",
		Name:        "GCASH",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_CASHINCASHOUT.String(),
	},
	{
		Stype:       "PAYMAYA",
		Name:        "PAYMAYA",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_CASHINCASHOUT.String(),
	},
	{
		Stype:       "DRAGONPAY",
		Name:        "DRAGONPAY",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_CASHINCASHOUT.String(),
	},
	{
		Stype:       "COINS",
		Name:        "COINS",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_CASHINCASHOUT.String(),
	},
	{
		Stype:       "PERAHUB",
		Name:        "PERAHUB",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_CASHINCASHOUT.String(),
	},
	{
		Stype:       "DISKARTECH",
		Name:        "DISKARTECH",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_CASHINCASHOUT.String(),
	},
}

var RTAPartnersList = []*storage.PartnerList{
	{
		Stype:       "METROBANK",
		Name:        "METROBANK",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTOACCOUNT.String(),
	},
	{
		Stype:       "BPI",
		Name:        "BPI",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTOACCOUNT.String(),
	},
	{
		Stype:       "UNIONBANK",
		Name:        "UNIONBANK",
		Status:      spb.PartnerStatusType_ENABLED.String(),
		ServiceName: sVcpb.ServiceType_REMITTOACCOUNT.String(),
	},
}
