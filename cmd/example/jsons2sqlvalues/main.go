package main

import (
	"context"
	"fmt"
	"iter"
	"log"
	"strings"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"

	sw "github.com/takanoriyanagitani/go-rowlike2sqlvalues/writer"

	js "github.com/takanoriyanagitani/go-rowlike2sqlvalues/rowlike/json/std"
)

var values IO[iter.Seq2[map[string]sw.Value, error]] = js.
	AnyToValueToStdinToValuesDefault

func mapsSink(m iter.Seq2[map[string]sw.Value, error]) IO[Void] {
	return func(ctx context.Context) (Void, error) {
		var smap map[string]string = map[string]string{}
		var buf strings.Builder

		var vw sw.ValueWriter = sw.ValueWriterNew(&buf)

		for row, e := range m {
			clear(smap)

			select {
			case <-ctx.Done():
				return Empty, ctx.Err()
			default:
			}

			if nil != e {
				return Empty, e
			}

			for key, val := range row {
				buf.Reset()
				_, e := val(vw)(ctx)
				if nil != e {
					return Empty, e
				}
				smap[key] = buf.String()
			}

			fmt.Printf("%v\n", smap)
		}
		return Empty, nil
	}
}

var vals2sink IO[Void] = Bind(
	values,
	mapsSink,
)

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return vals2sink(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
