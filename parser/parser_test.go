package parser

import (
	"testing"
	"strings"
)

func TestSimpleParse(t *testing.T) {
	doc := `IMPORT foo`
	input := strings.NewReader(doc)

	h, err := Parse(input)
	if err != nil {
		t.Errorf("Surprise! Error: %s", err)
	}

	handy := h.(*handler)

	if len(handy.imports) != 1 {
		t.Errorf("Expected 1 import, got %d", len(handy.imports))
	}

	if handy.imports[0] != "`foo`" {
		t.Errorf("Expected \"foo\", got %s", handy.imports[0])
	}
}

func TestParseImports(t *testing.T) {
	doc := `IMPORT foo bar baz IMPORT foo2 IMPORT bar2`
	input := strings.NewReader(doc)

	h, err := Parse(input)
	if err != nil {
		t.Errorf("Surprise! Error: %s", err)
	}

	handy := h.(*handler)

	if len(handy.imports) != 5 {
		t.Errorf("Expected 1 import, got %d", len(handy.imports))
	}

	expects := []string{"`foo`", "`bar`", "`baz`", "`foo2`", "`bar2`"}
	for i, expect := range expects {
		if handy.imports[i] != expect {
			t.Errorf("Expected %s, got %s", expect, handy.imports[i])
		}
	}

}

func TestParseSimpleRoutes(t *testing.T) {
	doc := `IMPORT foo
	ROUTE foo bar
	ROUTE "FOO" "BAR"
	ROUTE
	ROUTE "foo" "bar bar bar"


	`
	input := strings.NewReader(doc)

	h, err := Parse(input)
	if err != nil {
		t.Errorf("Surprise! Error: %s", err)
	}

	handy := h.(*handler)

	// Canary
	if len(handy.imports) != 1 {
		t.Errorf("What? No imports?")
	}

	if len(handy.routes) != 4 {
		t.Errorf("Expected 4 routes, got %d", len(handy.routes))
	}

	names := []string{"`foo`", "`FOO`", "", "`foo`"}
	descs := []string{"`bar`", "`BAR`", "", "`bar bar bar`"}

	for i, name := range names {
		if handy.routes[i].name != name {
			t.Errorf("Expected name to be %s, got %s", name, handy.routes[i].name)
		}
		if handy.routes[i].description != descs[i] {
			t.Errorf("Expected name to be %s, got %s", descs[i], handy.routes[i].description)
		}
	}
}

func TestParseFullRoute(t *testing.T) {
	doc := `
IMPORT github.com/Masterminds/codl

ROUTE matt "Butcher"
	DOES «foo.Bar» "foo"
		USING "param1" "default value"
		FROM cxt:query get:q
		USING "param2" FROM cxt:p2

ROUTE route2 "Another description"
	DOES «foo.Baz» "baz"
		USING "param1" USING param2 USING param3

ROUTE route3 "Route Three"
	INCLUDES one
	INCLUDES two
	DOES «three»
		`

	input := strings.NewReader(doc)

	h, err := Parse(input)
	if err != nil {
		t.Errorf("Surprise! Error: %s", err)
	}

	handy := h.(*handler)

	if handy.routes[0].name != "`matt`" {
		t.Errorf("Expected `matt`")
	}
	if handy.routes[0].commands[0].name != "`foo`" {
		t.Errorf("Expected first command to be named foo.")
	}
	if handy.routes[0].commands[0].cmd != "foo.Bar" {
		t.Errorf("Expected first command to be foo.Bar.")
	}
	if handy.routes[0].commands[0].params[0].from[1] != "`get:q`" {
		t.Errorf("Expected second FROM to be get:q")
	}
	if handy.routes[0].commands[0].params[1].name != "`param2`" {
		t.Errorf("Expected second USING to be param2")
	}
	if handy.routes[0].commands[0].params[1].defval != "" {
		t.Errorf("Expected second USING to have empty default value.")
	}

	// Test route 3
	if handy.routes[2].commands[1].cmdType != cmdInclude {
		t.Errorf("Expected 3rd command to have an INCLUDES in slot 2")
	}
	if handy.routes[2].commands[1].name !=  "`two`" {
		t.Errorf("Expected 3rd command to be two, got %s", handy.routes[2].commands[1].name)
	}
}
