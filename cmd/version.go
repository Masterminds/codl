package cmd

import (
	"github.com/Masterminds/cookoo"
	"fmt"
)

func Version(c cookoo.Context, p *cookoo.Params) (interface{}, cookoo.Interrupt) {
	ver := cookoo.GetString("version", "unstable", p)
	fmt.Printf("Version %s\n", ver)

	return ver, nil
}
