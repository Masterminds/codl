package cmd

import (
	"github.com/Masterminds/cookoo"
	"path/filepath"
)

// This file contains commands related to the file system.

// FindCodl finds all CODL files (*.codl) in a given directory.
func FindCodl(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	dir := cookoo.GetString("dir", ".", p)

	where := filepath.Join(dir, "*.go")
	files, err := filepath.Glob(where)

	return files, err
}
