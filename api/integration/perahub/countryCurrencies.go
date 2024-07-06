package perahub

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"brank.as/petnet/serviceutil/logging"
)

type currencyCache struct {
	*sync.Mutex
	lst map[string]country
}

type country struct {
	last    time.Time
	exp     time.Time
	curr    []ISOCombined
	currMap map[string]string
}

type ISOCombined struct {
	CountryName  string `json:"COUNTRY_LONG"`
	CountryNum   string `json:"ISO_COUNTRY_NUM_CD"`
	CountryCd    string `json:"ISO_COUNTRY_CD"`
	CurrencyCd   string `json:"CURRENCY_CD"`
	CurrencyNum  string `json:"ISO_CURRENCY_NUM_CD"`
	CurrencyName string `json:"CURRENCY_NAME"`
}

// CurrencyCodes ...
func (s *Svc) CurrencyCodes(ctx context.Context, cty, cur string) (map[string]string, error) {
	s.curr.Lock()
	defer s.curr.Unlock()
	ctyCur := fmt.Sprintf("%s %s", cty, cur)
	if s.curr.lst[ctyCur].exp.Before(time.Now()) {
		c, err := s.CurrencyCodesRaw(ctx, ctyCur)
		if err != nil {
			return nil, err
		}
		cy := country{
			last: time.Now(),
			exp:  time.Now().Add(s.exp),
			curr: c,
		}
		m := make(map[string]string, len(c))
		for _, v := range c {
			m[v.CountryCd] = v.CurrencyCd
		}
		cy.currMap = m
		s.curr.lst[ctyCur] = cy
	}
	return s.curr.lst[ctyCur].currMap, nil
}

func (s *Svc) CurrencyCodesRaw(ctx context.Context, ctycur string) ([]ISOCombined, error) {
	const mod, modReq = "wudas", "GetCountriesCurrencies"
	req, err := s.newParahubRequest(ctx, mod, modReq, struct {
		Cty string `json:"OriginatingCountry"`
	}{Cty: ctycur})
	if err != nil {
		return nil, err
	}

	resp, err := s.post(ctx, s.moduleURL(mod, ""), *req)
	if err != nil {
		return nil, err
	}

	res := []ISOCombined{}
	if err := json.Unmarshal(resp, &res); err != nil {
		logging.FromContext(ctx).WithField("body", string(resp)).Error("unmarshal failed")
		return nil, err
	}

	return res, nil
}
