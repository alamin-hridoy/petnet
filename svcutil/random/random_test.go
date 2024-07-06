package random

import (
	"testing"
)

func TestRandom(t *testing.T) {
	invCode := InvitationCode(32)
	invCode2 := InvitationCode(32)
	if invCode == invCode2 {
		t.Fatal("codes are equal")
	}
}
