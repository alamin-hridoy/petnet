package wise

import (
	"context"
	"fmt"

	"brank.as/petnet/api/core"
)

func (s *Svc) GetUser(ctx context.Context, req core.GetUserRequest) (*core.GetUserResponse, error) {
	return nil, fmt.Errorf("service not available for WISE")
}
