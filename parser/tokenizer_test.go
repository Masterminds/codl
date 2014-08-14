package parser

import (
	"testing"
	"strings"
	"io"
)

func TestMainLoop(t *testing.T) {
	r := strings.NewReader("    ")
	l := new(ListenerFixture)
	z := NewTokenizer(r, l)

	// Just make sure we are advancing.
	for z.lastErr == nil {
		z.Next()
	}

	if l.err != io.EOF {
		t.Errorf("Listener got error %s", l.err)
	}

	if z.lastErr != io.EOF {
		t.Errorf("Got error that evaded the listener: %s", z.lastErr)
	}

}

func TestStrings(t *testing.T) {

	expect := "this is a string"
	wrapped := `"` + expect + `"`

	// First: quoted strings
	r := strings.NewReader(wrapped)
	l := new(ListenerFixture)
	z := NewTokenizer(r, l)

	z.Next()
	if l.last != expect {
		t.Errorf("Expected '%s', got '%s'", expect, l.last)
	}

	// Second: barewords
	expect = "bareword"
	r = strings.NewReader(expect)
	z = NewTokenizer(r, l)

	z.Next()
	if l.last != expect {
		t.Errorf("Expected '%s', got '%s'", expect, l.last)
	}

	// Third: code literals
	expect = "this is code"
	wrapped = "`" + expect + "`"
	r = strings.NewReader(wrapped)
	z = NewTokenizer(r, l)

	z.Next()
	if l.last != expect {
		t.Errorf("Expected '%s', got '%s'", expect, l.last)
	}

	// Alternate code literals
	expect = "this is code"
	wrapped = "«" + expect + "»"
	r = strings.NewReader(wrapped)
	z = NewTokenizer(r, l)

	z.Next()
	if l.last != expect {
		t.Errorf("Expected '%s', got '%s'", expect, l.last)
	}

}

func TestKeywords(t *testing.T) {

	expectMap := map[string]string {
		"IMPORT": "_IMPORT",
		"INCLUDES": "_INCLUDES",
		"ROUTE": "_ROUTE",
		"USING":"_USING",
		"DOES": "_DOES",
		"FROM": "_FROM",
		"        FROM": "_FROM",
		"IMPORTs": "IMPORTs", // This should be interpreted as a string.
		"FROMs": "FROMs", // This should be interpreted as a string.
		"DOE": "DOE", // This should be interpreted as a string.
		"ROUTER": "ROUTER", // This should be interpreted as a string.
	}

	for input, output := range expectMap {
		r := strings.NewReader(input)
		l := new(ListenerFixture)
		z := NewTokenizer(r, l)
		z.Next()

		if output != l.last {
			t.Errorf("Expected '%s', but got '%s'", output, l.last)
		}
	}

}

type ListenerFixture struct {
	last string
	err error
}

func (l *ListenerFixture) Error(err error){
	l.err = err
}
func (l *ListenerFixture) Literal(str string){
	l.last = str
}
func (l *ListenerFixture) Strval(str string){
	l.last = str
}
func (l *ListenerFixture) Import(){
	l.last = "_IMPORT"
}
func (l *ListenerFixture) Includes(){
	l.last = "_INCLUDES"
}
func (l *ListenerFixture) Route(){
	l.last = "_ROUTE"
}
func (l *ListenerFixture) Using(){
	l.last = "_USING"
}
func (l *ListenerFixture) Does(){
	l.last = "_DOES"
}
func (l *ListenerFixture) From(){
	l.last = "_FROM"
}
