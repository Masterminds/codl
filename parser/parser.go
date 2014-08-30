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
	cmdCommand = iota
	cmdInclude
	cmdDoes
)

type Using struct {
	From []string
	Name, DefaultVal string
}

type Command struct {
	cmdType int
	Name, Cmd string
	Params []*Using
	currentParam *Using
}

func (c *Command) IsIncludes() bool {
	return c.cmdType == cmdInclude
}

type Route struct {
	Name, Description string
	Commands []*Command
	currentCommand *Command
}

type handler struct {
	mode int
	imports []string
	routes []*Route
	err error

	currentRoute *Route
}

func Parse(input io.Reader) (EventHandler, error) {
	l := &handler {
		mode: TopMode,
		imports: []string{},
		routes: []*Route{},
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

func (l *handler) Package() string {
	return "routes"
}

func (l *handler) Routes() []*Route {
	return l.routes
}
func (l *handler) Imports() []string {
	return l.imports
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
		if len(cc.Cmd) > 0 {
			l.err = fmt.Errorf("DOES %s already has a command.", cc.Name)
			return
		}
		cc.Cmd = str
	case UsingMode:
		cc := l.currentRoute.currentCommand
		// In Using mode, we can take a default that is a literal.
		if len(cc.currentParam.Name) == 0 {
			l.err = fmt.Errorf("USING requires a name that is not a literal.")
			return
		} else if len(cc.currentParam.DefaultVal) > 0 {
			l.err = fmt.Errorf("USING only allows one default value")
			return
		}
		cc.currentParam.DefaultVal = str
	}
}
func (l *handler) Strval(str string){
	orig := str

	str = asString(str)

	switch l.mode {
	case TopMode:
		l.err = fmt.Errorf("String value is in the top scope: %s", str)
	case ImportMode:
		l.imports = append(l.imports, str)
	case RouteMode:
		if len(l.currentRoute.Name) == 0 {
			l.currentRoute.Name = str
		} else if len(l.currentRoute.Description) == 0 {
			l.currentRoute.Description = str
		} else {
			l.err = fmt.Errorf("ROUTE takes one name and one description. No place for %s", str)
		}
	case IncludeMode:
		if len(l.currentRoute.currentCommand.Name) > 0 {
			fmt.Errorf("INCLUDE takes only one string. No place for %s", str)
			return
		}
		l.currentRoute.currentCommand.Name = str
	case DoesMode:
		if len(l.currentRoute.currentCommand.Cmd) == 0 {
			// We're gonna strategically ignore this rule. For pragmatic reasons,
			// it's a better user experience to "pretend" this string is a literal.
			//l.err = fmt.Errorf("DOES requires a `literal` for a command, not a string %s.", str)
			//fmt.Printf("Got a str for a command: %s\n", orig)
			l.currentRoute.currentCommand.Cmd = orig
		} else if len(l.currentRoute.currentCommand.Name) == 0 {
			l.currentRoute.currentCommand.Name = str
		} else {
			l.err = fmt.Errorf("DOES takes one literal and one string. No place for %s", str)
		}
	case UsingMode:
		cp := l.currentRoute.currentCommand.currentParam
		if len(cp.Name) == 0 {
			cp.Name = str
		} else if len(cp.DefaultVal) == 0 {
			cp.DefaultVal = str
		} else {
			l.err = fmt.Errorf("USING takes one literal and one string or literal. No place for %s", str)
		}
	case FromMode:
		cp := l.currentRoute.currentCommand.currentParam
		cp.From = append(cp.From, str)
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
func (l *handler) Includes(){
	switch l.mode {
	case TopMode, ImportMode:
		l.err = fmt.Errorf("INCLUDE is only allowed inside of a ROUTE")
	//case RouteMode, UsingMode, DoesMode, FromMode, IncludeMode:
	default:
		c := &Command{ cmdType: cmdInclude }
		l.mode = IncludeMode
		l.currentRoute.currentCommand = c
		l.currentRoute.Commands = append(l.currentRoute.Commands, c)
	}
}

func (l *handler) Route(){
	// No modes override this.
	l.mode = RouteMode
	r := new(Route)
	l.currentRoute = r
	l.routes = append(l.routes, r)
}

func (l *handler) Using() {
	switch l.mode {
	case TopMode, ImportMode, IncludeMode, RouteMode:
		l.err = fmt.Errorf("USING is only allowed inside of a DOES")
	case DoesMode, UsingMode, FromMode:
		u := new(Using)
		cc := l.currentRoute.currentCommand
		cc.currentParam = u
		cc.Params = append(cc.Params, u)
		l.mode = UsingMode
	}
}

func (l *handler) Does(){
	switch l.mode {
	case TopMode, ImportMode:
		l.err = fmt.Errorf("DOES can only appear inside of a ROUTE.")
	default:
		l.mode = DoesMode
		c := new(Command)
		l.currentRoute.Commands = append(l.currentRoute.Commands, c)
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
