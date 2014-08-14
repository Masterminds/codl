package parser

import (
	"github.com/Masterminds/sprig"
	"io"
	"text/template"
)

const bodyTpl = `package {{.Package}}

// This file is auto-generated by Codl.

import (
	"github.com/Masterminds/cookoo"
	{{range .Imports}}{{.}}
	{{end}}
)

func {{.Package | title }}Routes(reg *cookoo.Registry) {
	{{range .Routes}}reg.Route({{.Name}}, {{.Description}}){{range .Commands}}.
		Does({{.Cmd}}, {{.Name}}){{range .Params}}.
			Using({{.Name}}){{if .DefaultVal}}.WithDefault({{.DefaultVal}}){{end}}{{if .From}}.From({{.From}}){{end}}{{end}}{{end}}
	{{end}}
}
`

type Registry interface {
	Package() string
	Routes() []*Route
	Imports() []string
}

type Serializer struct {
	out io.Writer
	reg Registry
	tpl *template.Template
}

func NewSerializer(out io.Writer, reg Registry) *Serializer {
	s := &Serializer{out: out, reg: reg}
	s.compile()

	return s
}

func (s *Serializer) Write() error {
	return s.tpl.Execute(s.out, s.reg)
}

func (s *Serializer) compile() {
	s.tpl = template.Must(template.New("body").Funcs(sprig.TxtFuncMap()).Parse(bodyTpl))
}

