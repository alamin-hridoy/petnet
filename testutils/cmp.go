package testutils

import (
	"strings"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/protobuf/testing/protocmp"
)

func CmpProtoDiff(x, y interface{}, opts ...cmp.Option) string {
	opts = append(opts, protocmp.Transform())
	// opts = append(opts, cmpopts.EquateErrors())
	return cmp.Diff(x, y, opts...)
}

func CmpProtoEqual(x, y interface{}, opts ...cmp.Option) bool {
	opts = append(opts, protocmp.Transform())
	// opts = append(opts, cmpopts.EquateErrors())
	return cmp.Equal(x, y, opts...)
}

func IgnoreContent(cs ...string) cmp.Option {
	is := make(map[string]bool, 2*len(cs))

	for _, c := range cs {
		c := strings.ToLower(c)
		is["[\""+c+"\"]"] = true
		is["."+c] = true
	}

	return cmp.FilterPath(func(p cmp.Path) bool {
		return is[strings.ToLower(p.Last().String())]
	}, cmp.Ignore())
}
