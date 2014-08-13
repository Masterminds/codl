package parser

import (
	"testing"
	"strings"
)

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
