package perahub

import "testing"

func TestIsDomestic(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		srcCtry  string
		destCtry string
		want     int
	}{
		{
			name:     "Domestic",
			srcCtry:  "PH",
			destCtry: "PH",
			want:     1,
		},
		{
			name:     "International",
			srcCtry:  "PH",
			destCtry: "GB",
			want:     0,
		},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			r := IsDomestic(test.srcCtry, test.destCtry)
			if r != test.want {
				t.Errorf("want: %v, got: %v", r, test.want)
			}
		})
	}
}
