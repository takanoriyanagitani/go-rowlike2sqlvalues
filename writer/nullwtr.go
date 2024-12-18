package writer

import (
	"context"
	"io"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

func NullWriterNew(w io.Writer) func(struct{}) IO[Void] {
	return func(_ struct{}) IO[Void] {
		return func(_ context.Context) (Void, error) {
			_, e := w.Write([]byte("Null"))
			return Empty, e
		}
	}
}
