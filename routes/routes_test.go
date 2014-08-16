package routes

import (
	"github.com/Masterminds/codl/parser"
	"path/filepath"
	gparser "go/parser"
	gtoken "go/token"
	"testing"
	"strings"
	"bytes"
	//"fmt"
	"os"
)

// TestRoutes tests the *.codl files.
//
// It has the "side effect" of bootstrapping the CODL files, too.
func TestRoutes(t *testing.T) {
	//files, err := filepath.Glob(os.Getenv("GOPATH") + "/github.com/Masterminds/codl/routes/*.codl")
	files, err := filepath.Glob("*.codl")
	if err != nil {
		t.Errorf("! Fatal error: %s", err)
		return
	}

	//fmt.Printf("Found %v\n", files)

	for _, f := range files {

		outbase := strings.TrimSuffix(filepath.Base(f), ".codl")
		outname := outbase + ".go"

		input, err := os.Open(f)
		if err != nil {
			t.Errorf("Could not open %s. Skipping: %s", f, err)
			continue
		}

		h, err := parser.Parse(input)
		if err != nil {
			t.Errorf("Surprise! Error: %s", err)
		}

		var gosrc bytes.Buffer

		reg := h.(parser.Registry)
		ser := parser.NewSerializer(outbase, "routes", &gosrc, reg)
		if err := ser.Write(); err != nil {
			t.Errorf("Failed to serialize: %s", err)
		}

		fs := gtoken.NewFileSet()
		_, err = gparser.ParseFile(fs, outname, gosrc.String(), 0)
		if err != nil {
			t.Errorf("Failed to parse the resulting Go file: %s", err)
		}

		outfile, err := os.Create(outname)
		if err != nil {
			t.Errorf("Could not create an output file: %s", err)
		}
		defer outfile.Close()
		outfile.Write(gosrc.Bytes())

	}
}
