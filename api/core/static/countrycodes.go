package static

import (
	"context"

	"brank.as/petnet/api/storage"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// CountryTo4217 returns the ISO4217 currency code associated with the ISO3166 country code.
func (s *Svc) CountryTo4217(ctx context.Context, country, currency, dest string) (string, error) {
	m, err := s.cty.CurrencyCodes(ctx, country, currency)
	if err != nil {
		return "", err
	}
	if c, ok := m[dest]; ok {
		return c, nil
	}
	return "", status.Error(codes.InvalidArgument, "unsupported destination country")
}

func (s *Svc) GetISO(ctx context.Context, cc string) (*storage.ISOCty, error) {
	return s.st.GetISO(ctx, cc)
}
