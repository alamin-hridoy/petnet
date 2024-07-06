package microinsurance

import (
	"context"

	coreerror "brank.as/petnet/api/core/error"
	"brank.as/petnet/api/integration/microinsurance"
	migunk "brank.as/petnet/gunk/drp/v1/microinsurance"
)

// GetRelationships ...
func (s *MICoreSvc) GetRelationships(ctx context.Context) (*migunk.GetRelationshipsResult, error) {
	res, err := s.cl.GetRelationships(ctx)
	if err != nil {
		return nil, coreerror.ToCoreError(err)
	}

	list := make([]*migunk.Relationship, 0, len(res))
	for _, rl := range res {
		list = append(list,
			&migunk.Relationship{
				Relationship:      rl.Relationship,
				RelationshipValue: rl.RelationshipValue,
			})
	}

	return &migunk.GetRelationshipsResult{
		Relationships: list,
	}, nil
}

func toRelationship(r microinsurance.Relationship) *migunk.Relationship {
	return &migunk.Relationship{
		Relationship:      r.Relationship,
		RelationshipValue: r.RelationshipValue,
	}
}
