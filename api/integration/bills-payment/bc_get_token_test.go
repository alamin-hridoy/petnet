package bills_payment

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestBCGetToken(t *testing.T) {
	t.Parallel()

	st := newTestStorage(t)
	tests := []struct {
		name    string
		in      BCGetTokenRequest
		want    *BCGetTokenResponse
		wantErr bool
	}{
		{
			name: "Success",
			in: BCGetTokenRequest{
				GrantType: "client_credentials",
				TpaID:     "PP01",
				Scope:     "mecom-auth/all",
			},
			want: &BCGetTokenResponse{
				Code:    200,
				Message: "Success",
				Result: BCGetTokenResult{
					AccessToken: "eyJraWQiOiJraWoydFFERTZiSWxnOFE3enZMSmFZaE5jNXdlWHRzaVM0OW1vYVR4YWs0PSIsImFsZyI6IlJTMjU2In0.eyJzdWIiOiI0M2JoZThlcjk1YzM4ajg3MDdtbWpsYWxkIiwidG9rZW5fdXNlIjoiYWNjZXNzIiwic2NvcGUiOiJtZWNvbS1hdXRoXC9hbGwiLCJhdXRoX3RpbWUiOjE2NjY3NTEwMjgsImlzcyI6Imh0dHBzOlwvXC9jb2duaXRvLWlkcC5hcC1zb3V0aGVhc3QtMS5hbWF6b25hd3MuY29tXC9hcC1zb3V0aGVhc3QtMV9aZkJqVWVTeTMiLCJleHAiOjE2NjY3NTQ2MjgsImlhdCI6MTY2Njc1MTAyOCwidmVyc2lvbiI6MiwianRpIjoiMGRjNjZjNDUtODY2My00NWY5LWIyNWItNDVkMDg0YzZiMDY2IiwiY2xpZW50X2lkIjoiNDNiaGU4ZXI5NWMzOGo4NzA3bW1qbGFsZCJ9.dJpsEgqPQRSLrjNlz9QS74-gX3m-DNqiDfyBSl_NtaEupCguzFP_G0pRn1I2oxdfNsAUQog0a-NA2KYSQqA_CHmo81JoSPVaXmc7EhlPWl2ANkn7brVVSroRn3Of_cktS_gsxMWDe2t7Wxb8cSt5sFKaT2USA2fMqB1r4RoCmz3s6k9yvs_Niukmmzw_o6_4brte0-gxW6F1Jfx8dF2M27RXvMBVZ3fJQYHz_Njq-pfhuY-xlbUZvgyCywGOihZ7V-jvxp2p9vh3mIgYVtrEkY08vJeHdXpv2VxbZYd3snqdI4CO68jNONpqOqzZHSsRyE3hV6JNGtJThWZ0ezRqFw",
					ExpiresIn:   3600,
					TokenType:   "Bearer",
				},
				RemcoID: 2,
			},
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			s, _ := newTestSvc(t, st)
			got, err := s.BCGetToken(context.Background(), test.in)
			if (err != nil) != test.wantErr {
				t.Fatalf("BCGetToken() error = %v, wantErr %v", err, test.wantErr)
			}
			if !cmp.Equal(test.want, got) {
				t.Error(cmp.Diff(test.want, got))
			}
		})
	}
}
