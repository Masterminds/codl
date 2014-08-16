package main

import (
	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/cli"
	"github.com/Masterminds/codl/routes"
	//"github.com/Masterminds/codl/parser"
	"path/filepath"
	"strings"
	"flag"
	"fmt"
	"os"
)

// Overridden during compilation
var version = "Development"
var Summary = "Converts CODL files to Go source code"
var Description = `The CODL transformer transform CODL files into Go source code.

Commands:

- help: Show help text and exit.
- build: Convert ".codl" files to ".go" files.
`

func main() {
	reg, router, cxt := cookoo.Cookoo()
	flags := flag.NewFlagSet("global", flag.PanicOnError)
	flags.Bool("h", false, "Show help text and exit.")

	routes.AppRoutes(reg)

	err := cli.New(reg, router, cxt).Help(Summary, Description, flags).RunSubcommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}

}

func startup(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {


	files := os.Args[1:]
	for _, f := range files {
		p, _ := filepath.Abs(f)
		dir := filepath.Dir(p)
		name := strings.TrimSuffix(filepath.Base(p), ".codl")
		fmt.Printf("Got %s %s %s\n", dir, p, name)
	}

	return nil, nil
}
