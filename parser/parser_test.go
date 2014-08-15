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
		if handy.routes[i].Name != name {
			t.Errorf("Expected name to be %s, got %s", name, handy.routes[i].Name)
		}
		if handy.routes[i].Description != descs[i] {
			t.Errorf("Expected name to be %s, got %s", descs[i], handy.routes[i].Description)
		}
	}
}

func TestParseRuleBreaker(t *testing.T) {
	doc := `ROUTE foo DOES bare.Word thingy`
	input := strings.NewReader(doc)

	h, err := Parse(input)
	if err != nil {
		t.Errorf("Surprise! Error: %s", err)
	}

	handy := h.(*handler)

	if handy.routes[0].Commands[0].Name != "`thingy`" {
		t.Errorf("Unexpected name: %s", handy.routes[0].Commands[0].Name)
	}
	if handy.routes[0].Commands[0].Cmd != "bare.Word" {
		t.Errorf("Expected rule-breaker to be allowed: %s", handy.routes[0].Commands[0].Cmd)
	}
}

func TestParseFullRoute(t *testing.T) {
	doc := `
// Import codl
IMPORT github.com/Masterminds/codl

// First route.
ROUTE matt "Butcher" // That's my name!
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

	if handy.routes[0].Name != "`matt`" {
		t.Errorf("Expected `matt`")
	}
	if handy.routes[0].Commands[0].Name != "`foo`" {
		t.Errorf("Expected first command to be named foo.")
	}
	if handy.routes[0].Commands[0].Cmd != "foo.Bar" {
		t.Errorf("Expected first command to be foo.Bar.")
	}
	if handy.routes[0].Commands[0].Params[0].From[1] != "`get:q`" {
		t.Errorf("Expected second From to be get:q")
	}
	if handy.routes[0].Commands[0].Params[1].Name != "`param2`" {
		t.Errorf("Expected second USING to be param2")
	}
	if handy.routes[0].Commands[0].Params[1].DefaultVal != "" {
		t.Errorf("Expected second USING to have empty default value.")
	}

	// Test route 3
	if handy.routes[2].Commands[1].cmdType != cmdInclude {
		t.Errorf("Expected 3rd command to have an INCLUDES in slot 2")
	}
	if handy.routes[2].Commands[1].Name !=  "`two`" {
		t.Errorf("Expected 3rd command to be two, got %s", handy.routes[2].Commands[1].Name)
	}
}
