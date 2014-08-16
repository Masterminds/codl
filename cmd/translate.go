package cmd

import (
	"github.com/Masterminds/codl/parser"
	"github.com/Masterminds/cookoo"
	"strings"
	"path"
	"fmt"
	"os"
	"io"
)

const ExitNoFiles = 2

func Translate(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	files := p.Get("files", []string{}).([]string)

	if len(files) == 0 {
		fmt.Fprintf(os.Stderr, "No CODL files found. Quitting.\n")
		os.Exit(ExitNoFiles)
	}

	created := []string{}
	for _, fname := range files {
		basedir := path.Dir(fname)
		basename := strings.TrimSuffix(path.Base(fname), ".codl")
		newname := path.Join(basedir, basename + ".go")

		var input io.Reader
		var output io.Writer
		var err error
		if input, err = os.Open(fname); err != nil {
			return created, err
		}

		if output, err = os.Create(newname); err != nil {
			return created, err
		}

		if err := translate(basename, input, output); err != nil {
			return created, fmt.Errorf("Fatal error in %s: %s", fname, err)
		}

		created = append(created, newname)
	}

	return created, nil
}

func translate(basename string, in io.Reader, out io.Writer) error {
	h, err := parser.Parse(in)
	if err != nil {
		return err
	}

	reg := h.(parser.Registry)
	ser := parser.NewSerializer(basename, out, reg)
	if err := ser.Write(); err != nil {
		return err
	}
	return nil
}
