package partner

import (
	"context"

	ppb "brank.as/petnet/gunk/drp/v1/partner"
	"brank.as/petnet/serviceutil/logging"
)

func (s *Svc) GetPartnersRemco(ctx context.Context) (*ppb.GetPartnersRemcoResponse, error) {
	log := logging.FromContext(ctx).WithField("method", "core.partner.GetPartnersRemco")
	res, err := s.ph.PerahubGetRemcoID(ctx)
	if err != nil {
		log.WithError(err).Error("get partners remco from perahub")
		return nil, err
	}

	if res == nil {
		log.WithError(err).Error("failed to get partners remco from perahub")
		return nil, err
	}

	if len(res.Result) == 0 {
		log.WithError(err).Error("failed to get partners remco from perahub")
		return nil, err
	}

	var partners []*ppb.PerahubGetRemcoIDResult
	for _, v := range res.Result {
		partners = append(partners, &ppb.PerahubGetRemcoIDResult{
			ID:   v.ID.String(),
			Name: v.Name,
		})
	}

	return &ppb.GetPartnersRemcoResponse{
		Code:    res.Code.String(),
		Message: res.Message,
		Result:  partners,
	}, nil
}
