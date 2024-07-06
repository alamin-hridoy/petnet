package contextwriter

import (
	"context"
	"io"
)

func New(w io.WriterTo) WriterTo {
	return WriterTo{WriterTo: w}
}

type WriterTo struct {
	io.WriterTo
}

// WriteTo the writer.
func (wt WriterTo) WriteTo(ctx context.Context, w io.Writer) error {
	if _, err := wt.WriterTo.WriteTo(w); err != nil {
		return err
	}
	return ctx.Err()
}
