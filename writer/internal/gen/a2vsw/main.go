package main

import (
	"bufio"
	"context"
	"io"
	"os"
	"text/template"

	. "github.com/takanoriyanagitani/go-rowlike2sqlvalues/util"
)

var filename IO[string] = Of("any2valsw.go")

type TypePair struct {
	TypeHint  string
	Primitive string
}

type Config struct {
	Filename string
	Pairs    []TypePair
}

var config IO[Config] = Bind(
	All(filename),
	Lift(func(s []string) (Config, error) {
		return Config{
			Filename: s[0],
			Pairs: []TypePair{
				{TypeHint: "String", Primitive: "string"},
				{TypeHint: "Bytes", Primitive: "[]byte"},
				{TypeHint: "Int", Primitive: "int32"},
				{TypeHint: "Long", Primitive: "int64"},
				{TypeHint: "Float", Primitive: "float32"},
				{TypeHint: "Double", Primitive: "float64"},
				{TypeHint: "Boolean", Primitive: "bool"},
				{TypeHint: "Null", Primitive: "struct{}"},
				{TypeHint: "Time", Primitive: "time.Time"},
				{TypeHint: "Uuid", Primitive: "uuid.UUID"},
			},
		}, nil
	}),
)

var tmpl *template.Template = template.Must(template.ParseFiles(
	"./internal/gen/a2vsw/a2vsw.tmpl",
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
