package cors

import (
	"errors"
	"strings"
)

type Matcher struct {
	allowedOrigins  []string
	allowedWOrigins []wildcard
}

func NewOriginMatcher(origins []string) (*Matcher, error) {
	m := &Matcher{}

	for _, origin := range origins {
		origin = strings.ToLower(origin)

		if origin == "*" {
			return nil, errors.New("wildcard origin is not allowed")
		} else if i := strings.IndexByte(origin, '*'); i >= 0 {
			w := wildcard{origin[0:i], origin[i+1:]}
			m.allowedWOrigins = append(m.allowedWOrigins, w)

		} else {
			m.allowedOrigins = append(m.allowedOrigins, origin)
		}
	}

	return m, nil
}

func (m *Matcher) IsAllowedOrigin(origin string) bool {
	if origin == "" {
		return false
	}

	origin = strings.ToLower(origin)
	for _, o := range m.allowedOrigins {
		if o == origin {
			return true
		}
	}

	for _, w := range m.allowedWOrigins {
		if w.match(origin) {
			return true
		}
	}

	return false
}

type wildcard struct {
	prefix string
	suffix string
}

func (w wildcard) match(s string) bool {
	return len(s) >= len(w.prefix)+len(w.suffix) && strings.HasPrefix(s, w.prefix) && strings.HasSuffix(s, w.suffix)
}
