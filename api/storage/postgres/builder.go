package postgres

import (
	"fmt"
	"strings"

	"brank.as/petnet/api/storage"
	"brank.as/petnet/svcutil/random"
	"github.com/lib/pq"
)

type option func(*config)

type config struct {
	jsonParent string
	compDate   bool
}

type Builder struct {
	query    string
	allwhere string
	args     map[string]interface{}
}

func NewBuilder(q string) *Builder {
	m := make(map[string]interface{})
	return &Builder{
		query: q,
		args:  m,
	}
}

type operator string

const (
	eq     = "="
	gt     = ">"
	lt     = "<"
	gtOrEq = ">="
	ltOrEq = "<="
	notEq  = "!="
	notIn  = "NOT IN"
	in     = "IN"
)

func CompareDate() option {
	return func(c *config) { c.compDate = true }
}

func JSONBParent(v string) option {
	return func(c *config) { c.jsonParent = v }
}

// Where adds to the query string and sets the args IF the value is not empty, this is for
// constructing dynamic queries. Without any option the "=" operator is used by default
func (b *Builder) Where(col string, op operator, v string, opts ...option) *Builder {
	c := &config{}
	for _, opt := range opts {
		opt(c)
	}

	if v == "" {
		return b
	}

	clause := "WHERE"
	if strings.Contains(b.query, "WHERE") {
		clause = "AND"
	}

	pmName := col + "_" + random.NumberString(10)
	switch {
	case c.jsonParent != "":
		wq := fmt.Sprintf(" %s %s'%s' %s :%s", clause, c.jsonParent, col, op, pmName)
		b.query = b.query + wq
		b.allwhere = b.allwhere + wq
	case c.compDate:
		wq := fmt.Sprintf(" %s cast(%s as date) %s :%s", clause, col, op, pmName)
		b.query = b.query + wq
		b.allwhere = b.allwhere + wq
	default:
		wq := fmt.Sprintf(" %s %s %s :%s", clause, col, op, pmName)
		b.query = b.query + wq
		b.allwhere = b.allwhere + wq
	}
	b.args[pmName] = v
	return b
}

// In adds to the query string and sets arguments if v is not 0
func (b *Builder) Any(col string, v []string) *Builder {
	if len(v) == 0 {
		return b
	}

	clause := "WHERE"
	if strings.Contains(b.query, "WHERE") {
		clause = "AND"
	}
	n := col + "_in"
	wq := fmt.Sprintf(" %s %s = ANY(:%s)", clause, col, n)
	b.query = b.query + wq
	b.allwhere = b.allwhere + wq
	b.args[n] = pq.Array(v)
	return b
}

// In adds to the query string and sets arguments if v is not 0
func (b *Builder) NotAny(col string, v []string) *Builder {
	if len(v) == 0 {
		return b
	}

	clause := "WHERE"
	if strings.Contains(b.query, "WHERE") {
		clause = "AND"
	}
	n := col + "_in"
	wq := fmt.Sprintf(" %s %s != ANY(:%s)", clause, col, n)
	b.query = b.query + wq
	b.allwhere = b.allwhere + wq
	b.args[n] = pq.Array(v)
	return b
}

// In adds to the query string and sets arguments if v is not 0
func (b *Builder) WhereNotIn(col string, v []string) *Builder {
	if len(v) == 0 {
		return b
	}

	clause := "WHERE"
	if strings.Contains(b.query, "WHERE") {
		clause = "AND"
	}
	n := col + "_in"
	wq := fmt.Sprintf(" %s %s NOT IN(:%s)", clause, col, n)
	b.query = b.query + wq
	b.allwhere = b.allwhere + wq
	b.args[n] = pq.Array(v)
	return b
}

func (b *Builder) Limit(v int) *Builder {
	if v == 0 {
		return b
	}

	b.query = b.query + fmt.Sprintf(" LIMIT %d", v)
	return b
}

func (b *Builder) Offset(v int) *Builder {
	if v == 0 {
		return b
	}

	b.query = b.query + fmt.Sprintf(" OFFSET %d", v)
	return b
}

func (b *Builder) SortByColumn(col string, o storage.SortOrder) *Builder {
	if col == "" {
		return b
	}

	if o == "" {
		o = storage.Asc
	}
	b.query = b.query + fmt.Sprintf(" ORDER BY %s %s", col, o)
	return b
}

func (b *Builder) AddTotalQuery(tc string) *Builder {
	b.query = strings.ReplaceAll(b.query, tc, b.allwhere)
	return b
}

func (b *Builder) SortByColumnTextToJson(col string, key string, o storage.SortOrder) *Builder {
	if col == "" {
		return b
	}

	if o == "" {
		o = storage.Asc
	}
	b.query = b.query + fmt.Sprintf(" ORDER BY cast(cast(%s AS json)->>'%s' as int) %s", col, key, o)
	return b
}
