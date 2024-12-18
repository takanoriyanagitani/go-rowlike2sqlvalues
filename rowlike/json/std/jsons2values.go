package jsons2vals

import (
	"bufio"
	"context"
	"encoding/json"
	"io"
	"iter"
	"os"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"

	sw "github.com/takanoriyanagitani/go-rowlike2sqlvalues/writer"
)

func DecoderToValues(
	dec *json.Decoder,
	any2value func(any) sw.Value,
) iter.Seq2[map[string]sw.Value, error] {
	return func(yield func(map[string]sw.Value, error) bool) {
		var buf map[string]any
		var err error

		var mapd map[string]sw.Value = map[string]sw.Value{}
		for {
			clear(buf)
			clear(mapd)

			err = dec.Decode(&buf)
			if io.EOF == err {
				return
			}
			if nil == err {
				for key, val := range buf {
					var v sw.Value = any2value(val)
					mapd[key] = v
				}
			}

			if !yield(mapd, err) {
				return
			}
		}
	}
}

func DecoderToValuesDefault(
	dec *json.Decoder,
) iter.Seq2[map[string]sw.Value, error] {
	return DecoderToValues(dec, sw.AnyToVal)
}

func ReaderToValues(
	rdr io.Reader,
	any2value func(any) sw.Value,
) iter.Seq2[map[string]sw.Value, error] {
	var br io.Reader = bufio.NewReader(rdr)
	var dec *json.Decoder = json.NewDecoder(br)
	return DecoderToValues(dec, any2value)
}

func StdinToValues(
	any2value func(any) sw.Value,
) iter.Seq2[map[string]sw.Value, error] {
	return ReaderToValues(os.Stdin, any2value)
}

func AnyToValueToStdinToValues(
	any2value func(any) sw.Value,
) IO[iter.Seq2[map[string]sw.Value, error]] {
	return func(
		_ context.Context,
	) (iter.Seq2[map[string]sw.Value, error], error) {
		return StdinToValues(any2value), nil
	}
}

var AnyToValueToStdinToValuesDefault IO[iter.
	Seq2[map[string]sw.Value, error]] = AnyToValueToStdinToValues(sw.AnyToVal)
