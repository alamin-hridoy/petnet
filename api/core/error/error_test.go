package error

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"

	"brank.as/petnet/api/integration/perahub"
)

func TestToCoreError(t *testing.T) {
	err := &perahub.Error{
		Code:       "07",
		GRPCCode:   codes.InvalidArgument,
		Msg:        "Test Error",
		UnknownErr: "",
		Type:       "PARTNER",
		Errors:     nil,
	}

	ee := ToCoreError(err)

	assert.Equal(t, "Test Error", ee.Message)
	assert.Equal(t, codes.InvalidArgument, ee.Code)
}
