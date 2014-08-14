package parser

import (
	"fmt"
	"io"
)

// Insertion modes
const (
	TopMode = iota
	ImportMode
	RouteMode
	IncludeMode
	UsingMode
	DoesMode
	FromMode
)

const (
	cmdInclude = iota
	cmdDoes
)

type using struct {
	from []string
	name, defval string
}

type command struct {
	cmdType int
	name, cmd string
	params []*using
	currentParam *using
}

type route struct {
	name, description string
	commands []*command
	currentCommand *command
}

type handler struct {
	mode int
	imports []string
	routes []*route
	err error

	currentRoute *route
}

func Parse(input io.Reader) (EventHandler, error) {
	l := &handler {
		mode: TopMode,
		imports: []string{},
		routes: []*route{},
	}
	z := NewTokenizer(input, l)

	for l.err == nil {
		z.Next()
	}

	if l.err != io.EOF {
		return l, l.err
	}

	return l, nil
}

func (l *handler) Err() error {
	return l.err
}

func (l *handler) Error(err error) {
	l.err = err
}

func (l *handler) Literal(str string) {
	switch l.mode {
	case TopMode, ImportMode, RouteMode, FromMode, IncludeMode:
		l.err = fmt.Errorf("Literals are only allowed in DOES and USING: %s", str)
	case DoesMode:
		cc := l.currentRoute.currentCommand
		if len(cc.cmd) > 0 {
			l.err = fmt.Errorf("DOES %s already has a command.", cc.name)
			return
		}
		cc.cmd = str
	case UsingMode:
		cc := l.currentRoute.currentCommand
		// In Using mode, we can take a default that is a literal.
		if len(cc.currentParam.name) == 0 {
			l.err = fmt.Errorf("USING requires a name that is not a literal.")
			return
		} else if len(cc.currentParam.defval) > 0 {
			l.err = fmt.Errorf("USING only allows one default value")
			return
		}
		cc.currentParam.defval = str
	}
}
func (l *handler) Strval(str string){

	str = asString(str)

	switch l.mode {
	case TopMode:
		l.err = fmt.Errorf("String value is in the top scope: %s", str)
	case ImportMode:
		l.imports = append(l.imports, str)
	case RouteMode:
		if len(l.currentRoute.name) == 0 {
			l.currentRoute.name = str
		} else if len(l.currentRoute.description) == 0 {
			l.currentRoute.description = str
		} else {
			l.err = fmt.Errorf("ROUTE takes one name and one description. No place for %s", str)
		}
	case IncludeMode:
		if len(l.currentRoute.currentCommand.name) > 0 {
			fmt.Errorf("INCLUDE takes only one string. No place for %s", str)
			return
		}
		l.currentRoute.currentCommand.name = str
	case DoesMode:
		if len(l.currentRoute.currentCommand.cmd) == 0 {
			l.err = fmt.Errorf("DOES requires a `literal` for a command, not a string %s.", str)
		} else if len(l.currentRoute.currentCommand.name) == 0 {
			l.currentRoute.currentCommand.name = str
		} else {
			l.err = fmt.Errorf("DOES takes one literal and one string. No place for %s", str)
		}
	case UsingMode:
		cp := l.currentRoute.currentCommand.currentParam
		if len(cp.name) == 0 {
			cp.name = str
		} else if len(cp.defval) == 0 {
			cp.defval = str
		} else {
			l.err = fmt.Errorf("USING takes one literal and one string or literal. No place for %s", str)
		}
	case FromMode:
		cp := l.currentRoute.currentCommand.currentParam
		cp.from = append(cp.from, str)
	}

}

func asString(str string) string {
	return fmt.Sprintf("`%s`", str)
}

func (l *handler) Import(){
	if l.mode != TopMode && l.mode != ImportMode {
		l.err = fmt.Errorf("IMPORT must be before first ROUTE (mode: %d != %d)", l.mode, TopMode)
		return
	}

	l.mode = ImportMode
}
func (l *handler) Include(){
	switch l.mode {
	case TopMode, ImportMode:
		l.err = fmt.Errorf("INCLUDE is only allowed inside of a ROUTE")
	//case RouteMode, UsingMode, DoesMode, FromMode, IncludeMode:
	default:
		c := &command{ cmdType: cmdInclude }
		l.mode = IncludeMode
		l.currentRoute.currentCommand = c
		l.currentRoute.commands = append(l.currentRoute.commands, c)
	}
}

func (l *handler) Route(){
	// No modes override this.
	l.mode = RouteMode
	r := new(route)
	l.currentRoute = r
	l.routes = append(l.routes, r)
}

func (l *handler) Using() {
	switch l.mode {
	case TopMode, ImportMode, IncludeMode, RouteMode:
		l.err = fmt.Errorf("USING is only allowed inside of a DOES")
	case DoesMode, UsingMode, FromMode:
		u := new(using)
		cc := l.currentRoute.currentCommand
		cc.currentParam = u
		cc.params = append(cc.params, u)
		l.mode = UsingMode
	}
}

func (l *handler) Does(){
	switch l.mode {
	case TopMode, ImportMode:
		l.err = fmt.Errorf("DOES can only appear inside of a ROUTE.")
	default:
		l.mode = DoesMode
		c := new(command)
		l.currentRoute.commands = append(l.currentRoute.commands, c)
		l.currentRoute.currentCommand = c
	}
}
func (l *handler) From(){
	if l.mode != UsingMode {
		l.err = fmt.Errorf("FROM can only appear insude of a USING")
		return
	}
	l.mode = FromMode
}
