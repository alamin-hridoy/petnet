package revenue_commission

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_unmarshalDSACommission(t *testing.T) {
	bodyStr := `{
		"id": 1,
		"dsa_code": "UB",
		"commission_type": "UNIONBANK",
		"tier": "1",
		"commissiont_amount": "40",
		"commission_currency": "1",
		"updated_by": "SONNY",
		"created_at": "2022-07-21T17:35:17.000000Z",
		"updated_at": "2022-07-21T17:35:17.000000Z",
		"effective_date": null
		}`

	var comm DSACommission
	err := json.Unmarshal([]byte(bodyStr), &comm)

	assert.Nil(t, err)
}
