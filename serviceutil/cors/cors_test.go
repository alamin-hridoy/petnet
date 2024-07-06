package cors

import (
	"errors"
	"testing"
)

func TestNewMatcher(t *testing.T) {
	t.Parallel()
	origins := []string{"https://brank.as", "*.brank.as", "*.bnk.to"}
	invalidOrigins := []string{"*"}

	tests := []struct {
		desc      string
		origins   []string
		wantError error
	}{
		{
			desc:    "Success",
			origins: origins,
		},
		{
			desc:      "WildcardOriginNotAllowed",
			origins:   invalidOrigins,
			wantError: errors.New("wildcard origin is not allowed"),
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			_, gotError := NewOriginMatcher(test.origins)
			if gotError != nil && errors.Is(gotError, test.wantError) {
				t.Fatalf("got code %v want %v", gotError, test.wantError)
			}
		})
	}
}

func TestIsAllowedOrigins(t *testing.T) {
	t.Parallel()
	origins := []string{"https://brank.as", "*.brank.as", "*.bnk.to"}
	validOrigin := "https://brank.as"

	tests := []struct {
		desc    string
		origins []string
		input   string
		want    bool
	}{
		{
			desc:    "Success",
			origins: origins,
			input:   validOrigin,
			want:    true,
		},
		{
			desc:    "OriginNotIncluded",
			origins: origins,
			input:   "some-origin",
			want:    false,
		},
		{
			desc:    "MissingOrigin",
			origins: origins,
			input:   "",
			want:    false,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.desc, func(t *testing.T) {
			t.Parallel()

			m, _ := NewOriginMatcher(test.origins)

			got := m.IsAllowedOrigin(test.input)

			want := test.want
			if got != want {
				t.Fatalf("got bool %v want %v", got, want)
			}
		})
	}
}
