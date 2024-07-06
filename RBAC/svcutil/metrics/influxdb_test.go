package metrics

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/influxdata/influxdb-client-go/v2/api/write"

	"brank.as/rbac/svcutil/testutils"
)

type MockWrite struct {
	p []write.Point
}

func NewMock() *MockWrite {
	return &MockWrite{p: []write.Point{}}
}

func (mw *MockWrite) WritePoint(p *write.Point) {
	b := &strings.Builder{}
	b.WriteString(p.Name())
	b.WriteString(": ")
	for _, f := range p.TagList() {
		b.WriteString(f.Key)
		b.WriteString(": ")
		b.WriteString(f.Value)
		b.WriteString(",  ")
	}
	for _, f := range p.FieldList() {
		b.WriteString(f.Key)
		b.WriteString(": ")
		fmt.Fprintf(b, "%v", f.Value)
		b.WriteString(",  ")
	}
	fmt.Println(b.String())
	mw.p = append(mw.p, *p)
}

func (mw *MockWrite) WriteRecord(_ string) {}

func (mw *MockWrite) Close() {}

func (mw *MockWrite) Flush() {}

func (mw *MockWrite) Errors() <-chan error {
	return nil
}

func (mw *MockWrite) Last() *write.Point {
	l := len(mw.p)
	if l == 0 {
		return nil
	}
	return &mw.p[l-1]
}

func TestTagSpan(t *testing.T) {
	idb := &Influxdb{reqMeasurement: "test_measure"}
	ctx := context.Background()

	t.Run("CreateTags", func(t *testing.T) {
		tags := map[string]string{"test": "create"}
		tctx := WithTags(ctx, tags)
		got := Tags(tctx)
		if !testutils.CmpProtoEqual(tags, got) {
			t.Error(testutils.CmpProtoDiff(tags, got))
		}
	})

	t.Run("Span", func(t *testing.T) {
		tags := map[string]string{"test": "span"}
		tctx := WithTags(ctx, tags)
		tctx = idb.ServerSpan(tctx, "span")
		got := Tags(tctx)
		if !testutils.CmpProtoEqual(tags, got) {
			t.Error(testutils.CmpProtoDiff(tags, got))
		}
		span := getSpan(tctx)
		if !testutils.CmpProtoEqual(tags, span.tags) {
			t.Error(testutils.CmpProtoDiff(tags, span.tags))
		}
	})

	t.Run("NestedSpan", func(t *testing.T) {
		tags1 := map[string]string{"test": "span"}
		tags2 := map[string]string{"test": "nested span"}
		tctx := WithTags(ctx, tags1)
		tctx = idb.ServerSpan(tctx, "span")
		nestCtx := WithTags(tctx, tags2)
		if got := Tags(tctx); !testutils.CmpProtoEqual(tags1, got) {
			t.Error(testutils.CmpProtoDiff(tags1, got))
		}
		if got := Tags(nestCtx); !testutils.CmpProtoEqual(tags2, got) {
			t.Error(testutils.CmpProtoDiff(tags2, got))
		}
	})
}

func TestFieldSpan(t *testing.T) {
	t.Run("CreateFields", func(t *testing.T) {
		ctx := context.Background()
		fields := map[string]interface{}{"test": "create"}
		tctx := WithFields(ctx, fields)
		got := Fields(tctx)
		if !testutils.CmpProtoEqual(fields, got) {
			t.Error(testutils.CmpProtoDiff(fields, got))
		}
	})

	t.Run("Span", func(t *testing.T) {
		m := NewMock()
		idb := &Influxdb{reqMeasurement: "test_measure", w: m}
		ctx := context.Background()

		fields := map[string]interface{}{"test": "span"}
		tctx := WithFields(ctx, fields)
		tctx = idb.ServerSpan(tctx, "span")
		got := Fields(tctx)
		if !testutils.CmpProtoEqual(fields, got) {
			t.Error(testutils.CmpProtoDiff(fields, got))
		}
		span := getSpan(tctx)
		if !testutils.CmpProtoEqual(fields, span.fields) {
			t.Error(testutils.CmpProtoDiff(fields, span.fields))
		}
	})

	t.Run("NestedSpan", func(t *testing.T) {
		m := NewMock()
		idb := &Influxdb{reqMeasurement: "test_measure", w: m}
		ctx := context.Background()

		fields1 := map[string]interface{}{"test": "span"}
		fields2 := map[string]interface{}{"test": "nested span"}
		tctx := WithFields(ctx, fields1)
		tctx = idb.ServerSpan(tctx, "span")
		nestCtx := WithFields(tctx, fields2)
		if got := Fields(tctx); !testutils.CmpProtoEqual(fields1, got) {
			t.Error(testutils.CmpProtoDiff(fields1, got))
		}
		if got := Fields(nestCtx); !testutils.CmpProtoEqual(fields2, got) {
			t.Error(testutils.CmpProtoDiff(fields2, got))
		}
	})

	t.Run("RequestRecord", func(t *testing.T) {
		m := NewMock()
		idb := &Influxdb{reqMeasurement: "test_measure", w: m}
		ctx := context.Background()

		fields1 := map[string]interface{}{"test": "request record"}
		tctx := WithFields(ctx, fields1)
		if IsRequestRecorder(tctx) {
			t.Error("should not be a request recorder")
		}
		tctx = idb.WithRequestMetrics(tctx)
		if !IsRequestRecorder(tctx) {
			t.Error("not a request recorder")
		}
		SetTag(tctx, "tag1", "test")
		idb.FinalizeRequest(tctx, map[string]string{
			"tag2": "test2",
		}, nil)
		wantTags := map[string]string{
			"tag1": "test",
			"tag2": "test2",
		}
		got := map[string]string{}
		for _, t := range m.Last().TagList() {
			got[t.Key] = t.Value
		}
		if !testutils.CmpProtoEqual(wantTags, got) {
			t.Error(testutils.CmpProtoDiff(wantTags, got))
		}
	})
}

func TestTransport(t *testing.T) {
	t.Skip()
	idb := &Influxdb{reqMeasurement: "test_measure", w: NewMock()}
	ctx := context.Background()
	cl := &http.Client{
		Timeout:   2 * time.Second,
		Transport: idb.NewTransport("", nil),
	}
	u, _ := url.Parse("https://google.com")
	req, _ := idb.NewRequest(ctx, http.MethodGet, u, nil)
	r, err := cl.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer r.Body.Close()
	fmt.Println(r)
}

func TestTransportInject(t *testing.T) {
	t.Skip()
	idb := &Influxdb{reqMeasurement: "test_measure", w: NewMock()}
	ctx := context.Background()
	cl := &http.Client{
		Timeout:   2 * time.Second,
		Transport: idb.NewTransport("", nil),
	}
	u, _ := url.Parse("https://google.com")
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	r, err := cl.Do(req)
	if err != nil {
		t.Error(err)
	}
	defer r.Body.Close()
	fmt.Println(r)
}
