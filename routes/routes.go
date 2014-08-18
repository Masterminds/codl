package routes

// This directory contains support go code for routes.

import (
	"flag"
)

var buildFlags *flag.FlagSet

func init() {
	buildFlags = flag.NewFlagSet("build", flag.PanicOnError)
	buildFlags.Bool("h", false, "Show build help")
	buildFlags.String("d", ".", "The directory to look for CODL files.")
}
