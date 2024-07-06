package common

import (
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_dynamicUrlModify(t *testing.T) {
	req, err := http.NewRequest("GET", "https://localhost/v1/drp/remco-sf", nil)
	require.Nil(t, err)

	rt := dynamicUrlModify(req, strings.ReplaceAll(req.URL.Path, "/v1/drp/", ""))

	assert.Equal(t, "GET_remco-sf", rt)
}
