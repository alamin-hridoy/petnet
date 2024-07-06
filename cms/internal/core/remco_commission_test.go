package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	revcom "brank.as/petnet/gunk/drp/v1/revenue-commission"
)

func Test_getCreateAndDeleteList(t *testing.T) {
	drpList := []revcom.RemcoCommissionFee{
		{
			FeeID:               11,
			RemcoID:             1,
			MinAmount:           "10",
			MaxAmount:           "100",
			ServiceFee:          "11",
			CommissionAmount:    "11",
			CommissionAmountOTC: "0",
			CommissionType:      revcom.CommissionType_CommissionTypeRange,
			TrxType:             revcom.TrxType_TrxTypeInbound,
			UpdatedBy:           "drpShouldNotDeleteOrCreate",
		},
		{
			FeeID:               22,
			RemcoID:             1,
			MinAmount:           "101",
			MaxAmount:           "1000",
			ServiceFee:          "22",
			CommissionAmount:    "22",
			CommissionAmountOTC: "0",
			CommissionType:      revcom.CommissionType_CommissionTypeRange,
			TrxType:             revcom.TrxType_TrxTypeInbound,
			UpdatedBy:           "shouldCreate",
		},
	}

	phList := []revcom.RemcoCommissionFee{
		{
			FeeID:               111,
			RemcoID:             1,
			MinAmount:           "10",
			MaxAmount:           "100",
			ServiceFee:          "11",
			CommissionAmount:    "11",
			CommissionAmountOTC: "0",
			CommissionType:      revcom.CommissionType_CommissionTypeRange,
			TrxType:             revcom.TrxType_TrxTypeInbound,
			UpdatedBy:           "phShouldNotDeleteOrCreate",
		},
		{
			FeeID:               222,
			RemcoID:             1,
			MinAmount:           "101",
			MaxAmount:           "1000",
			ServiceFee:          "22",
			CommissionAmount:    "20",
			CommissionAmountOTC: "0",
			CommissionType:      revcom.CommissionType_CommissionTypeRange,
			TrxType:             revcom.TrxType_TrxTypeInbound,
			UpdatedBy:           "shouldDelete",
		},
	}

	createList, deleteList := getCreateAndDeleteList(drpList, phList)

	assert.NotEmpty(t, createList)
	assert.NotEmpty(t, deleteList)

	require.Equal(t, 1, len(deleteList))
	require.Equal(t, 1, len(createList))

	assert.Equal(t, deleteList[0].FeeID, uint32(222))
	assert.Equal(t, createList[0].FeeID, uint32(22))
}
