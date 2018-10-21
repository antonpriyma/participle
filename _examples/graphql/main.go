package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/alecthomas/participle"
	"github.com/alecthomas/repr"
)

type File struct {
	Entries []*Entry `{ @@ }`
}

type Entry struct {
	Type   *Type   `  @@`
	Schema *Schema `| @@`
	Enum   *Enum   `| @@`
	Scalar string  `| "scalar" @Ident`
}

type Enum struct {
	Name  string   `"enum" @Ident`
	Cases []string `"{" { @Ident } "}"`
}

type Schema struct {
	Fields []*Field `"schema" "{" { @@ } "}"`
}

type Type struct {
	Name       string   `"type" @Ident`
	Implements string   `[ "implements" @Ident ]`
	Fields     []*Field `"{" { @@ } "}"`
}

type Field struct {
	Name       string      `@Ident`
	Arguments  []*Argument `[ "(" [ @@ { "," @@ } ] ")" ]`
	Type       *TypeRef    `":" @@`
	Annotation string      `[ "@" @Ident ]`
}

type Argument struct {
	Name    string   `@Ident`
	Type    *TypeRef `":" @@`
	Default *Value   `[ "=" @@ ]`
}

type TypeRef struct {
	Array       *TypeRef `(   "[" @@ "]"`
	Type        string   `  | @Ident )`
	NonNullable bool     `[ @"!" ]`
}

type Value struct {
	Symbol string `@Ident`
}

var (
	parser = participle.MustBuild(&File{})

	cli struct {
		Files []string `arg:"" type:"existingfile" required:"" help:"GraphQL schema files to parse."`
	}
)

func main() {
	ctx := kong.Parse(&cli)
	for _, file := range cli.Files {
		ast := &File{}
		r, err := os.Open(file)
		ctx.FatalIfErrorf(err)
		err = parser.Parse(r, ast)
		r.Close()
		ctx.FatalIfErrorf(err)
		repr.Println(ast)
	}
}
