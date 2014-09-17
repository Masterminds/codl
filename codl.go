package main

import (
	"github.com/Masterminds/cookoo"
	"github.com/Masterminds/cookoo/cli"
	"github.com/Masterminds/codl/routes"
	//"github.com/Masterminds/codl/parser"
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
- watch: Watch all .codl files in a directory for changes, and transform them.

Examples:

$ codl build -d routes/   # Convert all .codl files in routes/
$ codl watch -d routes/   # Watch routes/ for changed .codl files.
$ codl -h                 # Show global help.
$ codl watch -h           # Show help for the 'codl watch' command.
`

func main() {
	reg, router, cxt := cookoo.Cookoo()
	flags := flag.NewFlagSet("global", flag.PanicOnError)
	flags.Bool("h", false, "Show help text and exit.")

	// Used by repeat.
	cxt.Put("router", router)
	cxt.Put("version", version)

	routes.AppRoutes(reg)

	err := cli.New(reg, router, cxt).Help(Summary, Description, flags).RunSubcommand()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s", err)
		os.Exit(1)
	}

}
