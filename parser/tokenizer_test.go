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
	if l.last.(string) != expect {
		t.Errorf("Expected '%s', got '%s'", expect, l.last.(string))
	}

	// Second: barewords
	expect = "bareword"
	r = strings.NewReader(expect)
	z = NewTokenizer(r, l)

	z.Next()
	if l.last.(string) != expect {
		t.Errorf("Expected '%s', got '%s'", expect, l.last.(string))
	}

	// Third: code literals
	// Second: barewords
	expect = "this is code"
	wrapped = "`" + expect + "`"
	r = strings.NewReader(wrapped)
	z = NewTokenizer(r, l)

	z.Next()
	if l.last.(string) != expect {
		t.Errorf("Expected '%s', got '%s'", expect, l.last.(string))
	}

}

type ListenerFixture struct {
	last interface{}
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
func (l *ListenerFixture) Import(){}
func (l *ListenerFixture) Include(){}
func (l *ListenerFixture) Route(){}
func (l *ListenerFixture) Using(){}
func (l *ListenerFixture) Does(){}
func (l *ListenerFixture) From(){}
