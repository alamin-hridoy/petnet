package signup

import (
	"strings"
	"testing"

	upb "brank.as/petnet/gunk/dsa/v1/user"
)

func TestPasswordValidation(t *testing.T) {
	t.Parallel()

	tests := []struct {
		desc      string
		pwd       string
		wantError error
	}{
		{
			desc: "Success",
			pwd:  "Password1!",
		},
		{
			desc:      "MissingSpecial",
			pwd:       "Password1",
			wantError: errMissingSpecial,
		},
		{
			desc:      "MissingUpperCase",
			pwd:       "password1!",
			wantError: errMissingUpper,
		},
		{
			desc:      "MissingNumber",
			pwd:       "Password!",
			wantError: errMissingNumber,
		},
		{
			desc:      "ToShort",
			pwd:       "Passd1!",
			wantError: errIncorrectLength,
		},
		{
			desc:      "ToLong",
			pwd:       strings.Repeat("a", 62) + "P1!",
			wantError: errIncorrectLength,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()
			err := validate(&upb.SignupRequest{
				Username:  "user",
				FirstName: "first",
				LastName:  "last",
				Email:     "email@example.com",
				Password:  test.pwd,
			})
			if err != test.wantError {
				t.Fatalf("got error: %v want: %v", err, test.wantError)
			}
		})
	}
}
