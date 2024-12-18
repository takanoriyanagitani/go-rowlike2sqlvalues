package main

import (
	"context"
	"fmt"
	"iter"
	"log"
	"os"
	"strings"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"

	sw "github.com/takanoriyanagitani/go-rowlike2sqlvalues/writer"

	pp "github.com/takanoriyanagitani/go-rowlike2sqlvalues/rowlike/rdb/postgresql/pgcopy"
	ph "github.com/takanoriyanagitani/go-rowlike2sqlvalues/rowlike/rdb/postgresql/pgcopy/header"
)

var EnvValByKey func(string) IO[string] = Lift(
	func(key string) (string, error) {
		val, found := os.LookupEnv(key)
		switch found {
		case true:
			return val, nil
		default:
			return "", fmt.Errorf("env var %s missing", key)
		}
	},
)

var columnInfoJsonFilename IO[string] = EnvValByKey("ENV_COL_INFO_JSONL_NAME")

var columnInfoMap IO[map[int16]pp.ColumnInfo] = Bind(
	columnInfoJsonFilename,
	pp.FilenameToColumnInfoMap,
)

var pgheader IO[ph.PgcopySimpleHeader] = ph.HeaderFromStdinDefault

var stdin2pgrows IO[iter.Seq2[pp.PgRow, error]] = Bind(
	pgheader,
	func(_ ph.PgcopySimpleHeader) IO[iter.Seq2[pp.PgRow, error]] {
		return pp.PgRowsFromStdin
	},
)

var valueMap IO[iter.Seq2[map[string]sw.Value, error]] = Bind(
	columnInfoMap,
	func(
		cmap map[int16]pp.ColumnInfo,
	) IO[iter.Seq2[map[string]sw.Value, error]] {
		return Bind(
			stdin2pgrows,
			pp.ColumnMapToPgRows(cmap),
		)
	},
)

func valueMap2stdout(vmap iter.Seq2[map[string]sw.Value, error]) IO[Void] {
	return func(ctx context.Context) (Void, error) {
		smap := map[string]string{}
		var buf strings.Builder
		var vw sw.ValueWriter = sw.ValueWriterNew(&buf)
		for m, e := range vmap {
			if nil != e {
				return Empty, e
			}

			clear(smap)
			for key, val := range m {
				buf.Reset()
				_, e := val(vw)(ctx)
				if nil != e {
					return Empty, e
				}
				smap[key] = buf.String()
			}

			fmt.Printf("map: %v\n", smap)
		}
		return Empty, nil
	}
}

var stdin2values2stdout IO[Void] = Bind(
	valueMap,
	valueMap2stdout,
)

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return stdin2values2stdout(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		log.Printf("%v\n", e)
	}
}
