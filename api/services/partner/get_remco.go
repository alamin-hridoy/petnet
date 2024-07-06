package partner

import (
	"context"

	"brank.as/petnet/serviceutil/logging"
	"google.golang.org/protobuf/types/known/emptypb"

	ppb "brank.as/petnet/gunk/drp/v1/partner"
	spb "brank.as/petnet/gunk/dsa/v2/partnerlist"
)

func (s *Svc) GetPartnersRemco(ctx context.Context, in *emptypb.Empty) (*ppb.GetPartnersRemcoResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "service.partner.GetPartnersRemco")
	res, err := s.partner.GetPartnersRemco(ctx)
	if err != nil {
		log.WithError(err).Error("failed to get partners remco id")
		return nil, err
	}

	if res == nil {
		log.WithError(err).Error("failed to get partners remco id")
		return nil, err
	}

	for _, v := range res.GetResult() {
		gRes, err := s.plcl.GetPartnerList(ctx, &spb.GetPartnerListRequest{
			Name: v.GetName(),
		})
		if err != nil {
			log.WithError(err).Errorf("failed to get partner:%s", v.GetName())
			continue
		}

		if gRes == nil || gRes.PartnerList == nil || len(gRes.PartnerList) == 0 {
			log.WithError(err).Errorf("failed to get partner:%s", v.GetName())
			continue
		}

		_, err = s.plcl.UpdatePartnerList(ctx, &spb.UpdatePartnerListRequest{
			PartnerList: &spb.PartnerList{
				Stype:   gRes.PartnerList[0].Stype,
				RemcoID: v.GetID(),
			},
		})
		if err != nil {
			log.WithError(err).Errorf("failed to update remco id for partner:%s", v.GetName())
		}
	}

	return res, nil
}
