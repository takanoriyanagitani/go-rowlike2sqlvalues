package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"os"
	"strings"
	"text/template"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

var argLen int = len(os.Args)

var GetArgByIndex func(int) IO[string] = Lift(
	func(ix int) (string, error) {
		if ix < argLen {
			return os.Args[ix], nil
		}
		return "", fmt.Errorf("invalid argument index: %v", ix)
	},
)

// e.g, Boolean
var typeHint IO[string] = GetArgByIndex(1)

// e.g, boolean2value.go
var filename IO[string] = Bind(
	typeHint,
	Lift(func(s string) (string, error) {
		var low string = strings.ToLower(s)
		return low + "2value.go", nil
	}),
)

type Config struct {
	TypeHint string
	Filename string
}

var config IO[Config] = Bind(
	All(typeHint, filename),
	Lift(func(s []string) (Config, error) {
		return Config{
			TypeHint: s[0],
			Filename: s[1],
		}, nil
	}),
)

var tmpl *template.Template = template.Must(template.ParseFiles(
	"./internal/gen/toval/toval.tmpl",
))

func ExecuteTemplateToWriter(
	wtr io.Writer,
	t *template.Template,
	cfg Config,
) error {
	var bw *bufio.Writer = bufio.NewWriter(wtr)
	defer bw.Flush()
	return t.Execute(bw, cfg)
}

func ExecuteTemplateToFilename(
	t *template.Template,
	cfg Config,
) error {
	var filename string = cfg.Filename
	f, e := os.Create(filename)
	if nil != e {
		return e
	}
	defer f.Close()
	return ExecuteTemplateToWriter(f, t, cfg)
}

var executeTemplate IO[Void] = Bind(
	config,
	Lift(func(cfg Config) (Void, error) {
		return Empty, ExecuteTemplateToFilename(tmpl, cfg)
	}),
)

var sub IO[Void] = func(ctx context.Context) (Void, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	return executeTemplate(ctx)
}

func main() {
	_, e := sub(context.Background())
	if nil != e {
		panic(e)
	}
}
