package pgcopy2header

import (
	"context"
	"io"
	"os"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

type PgcopySimpleHeader struct {
	Signature       [11]byte
	Flags           [4]byte
	HeaderExtension [4]byte
}

func ReaderToHeader(rdr io.Reader) (PgcopySimpleHeader, error) {
	var ret PgcopySimpleHeader
	var buf [19]byte
	_, e := io.ReadFull(rdr, buf[:])
	if nil != e {
		return ret, e
	}

	copy(ret.Signature[:], buf[0:11])
	copy(ret.Flags[:], buf[11:15])
	copy(ret.HeaderExtension[:], buf[15:19])
	return ret, nil
}

func StdinToHeader(_ context.Context) (PgcopySimpleHeader, error) {
	return ReaderToHeader(os.Stdin)
}

var HeaderFromStdinDefault IO[PgcopySimpleHeader] = StdinToHeader
